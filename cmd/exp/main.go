package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
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

	// TODO: find out why it prints only in the console or only as a web page.
	// // In case we didn't get an error we can now stream the data into the resonse writer.
	// if _, err := io.Copy(w, &actualRes); err != nil {
	// 	log.Printf("Error when writing the response: %v", err)
	// 	http.Error(w, fmt.Sprintf("Error when writing the response: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	// We stream the data into os.stdout.
	if _, err := io.Copy(os.Stdout, &actualRes); err != nil {
		log.Printf("Error when writing the response: %v", err)
		http.Error(w, fmt.Sprintf("Error when writing the response: %v", err), http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templatePath := "hello.gohtml"

	user := struct {
		Name, Bio string
		Age       int
		Birthday  time.Time
		Address   struct {
			HouseNum     int
			Street, City string
		}
	}{
		Name:     "Ken",
		Bio:      "<script>alert(Haha you've been h4x0r3d!);</script>",
		Birthday: time.Date(1988, time.January, 21, 0, 0, 0, 0, time.UTC),
		Address: struct {
			HouseNum     int
			Street, City string
		}{
			HouseNum: 100,
			Street:   "Derech Haking",
			City:     "Prešov"},
	}

	executeTemplate(w, templatePath, user)
}

func main() {
	router := chi.NewRouter()
	router.Get("/", homeHandler)
	http.ListenAndServe(":3000", router)
}
