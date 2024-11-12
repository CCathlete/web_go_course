package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-yaml/yaml"
)

func executeTemplate(w http.ResponseWriter, templatePath string, templateBody any) {
	// Setting up the response's header before further processing.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Error when parsing the tamplate: %v", err)
		http.Error(w, fmt.Sprintf("Error when parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	// Writing the html page into a buffer to make sure we don't have an error
	// before writing to the respinse writer.
	var actualRes bytes.Buffer
	if err := tpl.Execute(&actualRes, templateBody); err != nil {
		http.Error(w, fmt.Sprintf("Error when executing html: %v", err), http.StatusInternalServerError)
		return
	}

	// In case we didn't get an error we can now stream the data into the resonse writer.
	if _, err := io.Copy(w, &actualRes); err != nil {
		log.Printf("Error when writing the response: %v", err)
		http.Error(w, fmt.Sprintf("Error when writing the response: %v", err), http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templatePath := "/home/ccat/Repos/web_go_course/templates/home.gohtml"
	executeTemplate(w, templatePath, nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	templatePath := "/home/ccat/Repos/web_go_course/templates/contact.gohtml"
	executeTemplate(w, templatePath, nil)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	// Creating a data bucket for the content.
	var qaYaml struct {
		Questions template.HTML `yaml:"Content"`
	}

	// Opening the file for reading only.
	file, err := os.Open("QA.yaml")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while opening internal QA file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Spilling the data into the bucket in small chunks.
	if err := yaml.NewDecoder(file).Decode(&qaYaml); err != nil && err != io.EOF {
		http.Error(w, fmt.Sprintf("Error while opening internal QA file: %v", err), http.StatusInternalServerError)
		return
	}

	// Converting the questions into a string , replacing parts and converting back to template.HTML.
	formattedContent := template.HTML(strings.ReplaceAll(string(qaYaml.Questions), "\n", "<br>"))
	qaYaml.Questions = formattedContent
	templatePath := "/home/ccat/Repos/web_go_course/templates/faq.gohtml"
	executeTemplate(w, templatePath, qaYaml)
}

func getAll(router *chi.Mux) {
	// router.Use(middleware.Logger)
	router.Get("/", homeHandler)
	router.Get("/contact", contactHandler)
	router.Get("/faq", faqHandler)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found.", http.StatusNotFound)
	})
}

func main() {
	myRouter := chi.NewRouter()
	// Set the routes in the router object.
	getAll(myRouter)
	fmt.Println("Starting the server on: 3000...")
	http.ListenAndServe(":3000", myRouter)
}
