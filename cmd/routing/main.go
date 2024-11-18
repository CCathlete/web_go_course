package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"webGo/views"

	"github.com/go-chi/chi/v5"
	"github.com/go-yaml/yaml"
)

// TODO: Need to fix the part of adding data parsed from the QA yaml.
// func executeTemplate(w http.ResponseWriter, innerData any) {
// 	// Setting up the response's header before further processing.
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	// Executing the template and writing the http page resulting from it.
// 	viewTpl.Execute(w, innerData)
// }

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	templatePath := "/home/ccat/Repos/web_go_course/templates/home.gohtml"
// 	executeTemplate(w, templatePath, nil)
// }

// func contactHandler(w http.ResponseWriter, r *http.Request) {
// 	templatePath := "/home/ccat/Repos/web_go_course/templates/contact.gohtml"
// 	executeTemplate(w, templatePath, nil)
// }

func getQuestionsTemplate(questionsPath string) (interface{}, error) {
	// Creating a data bucket for the content.
	var qaYaml struct {
		Questions template.HTML `yaml:"Content"`
	}

	// Opening the file for reading only.
	file, err := os.Open(questionsPath)
	if err != nil {
		return nil, fmt.Errorf("error while opening internal QA file: %w", err)
	}
	defer file.Close()

	// Spilling the data into the bucket in small chunks.
	if err := yaml.NewDecoder(file).Decode(&qaYaml); err != nil && err != io.EOF {
		return nil, fmt.Errorf("error while opening internal QA file: %w", err)
	}

	// Converting the questions into a string, replacing parts and converting back to template.HTML.
	formattedContent := template.HTML(strings.ReplaceAll(string(qaYaml.Questions), "\n", "<br>"))
	qaYaml.Questions = formattedContent // The questions in html template form and newlines.

	return &qaYaml, nil
}

func staticHandler(tpl *views.Template, routeSuffix string) http.HandlerFunc {
	// var data *template.HTML
	var data interface{}
	var err error
	if routeSuffix == "faq" {
		data, err = getQuestionsTemplate("QA.yaml")
		if err != nil {
			return func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, fmt.Sprintf("Error while parsing the questions yaml: %v", err), http.StatusInternalServerError)
			}
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, data)
	}
}

func getAll(router *chi.Mux) {
	// router.Use(middleware.Logger)
	// templateList := make(map[string]string) // Template list yaml routeSuffix: templatePath.
	yamlData, err := views.ParseYaml("templatePaths.yaml")
	if err != nil {
		panic(err)
	}
	templateList := yamlData.(map[interface{}]interface{})

	for routeSuffix, templatePath := range templateList {
		routeSuffix := routeSuffix.(string)
		templatePath := templatePath.(string)
		template, err := views.ParseTemplate(templatePath)
		if err != nil {
			panic(err)
		}
		router.Get(fmt.Sprintf("/%s", routeSuffix), staticHandler(template, routeSuffix))
	}
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
