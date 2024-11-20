package controllers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"webGo/views"

	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v2"
)

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

func StaticHandler(tpl *views.Template, routeSuffix string) http.HandlerFunc {
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

func GetAll(router *chi.Mux) {
	// router.Use(middleware.Logger)
	// templateList := make(map[string]string) // Template list yaml routeSuffix: templatePath.
	yamlData := views.Must(views.ParseYaml("templatePaths.yaml"))
	templateList := yamlData.(map[interface{}]interface{})

	for routeSuffix, templatePath := range templateList {
		routeSuffix := routeSuffix.(string)
		templatePath := templatePath.(string)

		// We need to assert at the end because Must returns an interface.
		template := views.Must(views.ParseTemplate(templatePath)).(*views.Template)
		router.Get(fmt.Sprintf("/%s", routeSuffix),
			StaticHandler(template, routeSuffix))
	}
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found.", http.StatusNotFound)
	})
}
