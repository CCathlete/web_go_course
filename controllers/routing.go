package controllers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"webGo/templates"
	"webGo/views"

	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v2"
)

func getQuestionsTemplate(questionsFileName string) (interface{}, error) {
	// Creating a data bucket for the content.
	var qaYaml struct {
		Questions template.HTML `yaml:"Content"`
	}

	// Opening the file for reading only.
	file, err := templates.FS.Open(questionsFileName)
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

func GetAll(router *chi.Mux) {
	// router.Use(middleware.Logger)
	templateList := map[string]string{
		"":           "home.gohtml",
		"contact":    "contact.gohtml",
		"faq":        "faq.gohtml",
		"myproducts": "In_construction.gohtml",
	}

	for routeSuffix, templatePath := range templateList {

		// We need to assert at the end because Must returns an interface.
		template := views.Must(views.ParseFS(templates.FS, templatePath, "tailwind.gohtml")).(*views.Template)
		router.Get(fmt.Sprintf("/%s", routeSuffix),
			StaticHandler(template, routeSuffix))
	}
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found.", http.StatusNotFound)
	})
}

func GetAllNoEmbedding(router *chi.Mux) {
	// router.Use(middleware.Logger)
	// templateList := make(map[string]string) // Template list yaml routeSuffix: templatePath.
	yamlData := views.Must(views.ParseYamlNoEmbedding("/home/ccat/Repos/web_go_course/templatePaths.yaml"))
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

func getQuestionsTemplateNoEmbedding(questionsPath string) (interface{}, error) {
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
