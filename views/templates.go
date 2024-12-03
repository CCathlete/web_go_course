package views

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

type Template struct {
	HtmlTpl *template.Template
}

func ParseFS(fs embed.FS, patterns ...string) (*Template, error) {
	// Initialises a template object with the name of the base template
	// to prapare it for parsFS.
	tpl := template.New(patterns[0])
	// Inner template functions must be defined before parsing.
	// I assume this is because the names of functions in the templates are already mentioned so when the parser sees a function it needs to know where to link it.
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
		},
	)
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return nil, fmt.Errorf("error when parsing template: %w", err)
	}

	return &Template{
		HtmlTpl: tpl,
	}, nil
}

func Must(tpl any, err error) any {
	if err != nil {
		panic(err)
	}

	return tpl
}

func (tpl *Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	// The inner template field is a pointer so we need to use a copy
	// before putting in the csrf token during execution to handle multiple users.
	innerTemplate, err := tpl.HtmlTpl.Clone()
	if err != nil {
		log.Printf("Error cloning template (Execute): %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	innerTemplate = innerTemplate.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
		},
	)
	// Writing the html page into a buffer to make sure we don't have an error
	// before writing to the respinse writer.
	var actualRes bytes.Buffer
	if err := innerTemplate.Execute(&actualRes, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, fmt.Sprintf("Error when rendering html: %v", err), http.StatusInternalServerError)
		return
	}

	// In case we didn't get an error we can now stream the data into the resonse writer.
	if _, err := io.Copy(w, &actualRes); err != nil {
		log.Printf("Error when writing the response: %v", err)
		http.Error(w, fmt.Sprintf("Error when writing the response: %v", err), http.StatusInternalServerError)
		return
	}
}

// func ParseTemplate(templatePath string) (*Template, error) {
// 	tpl, err := template.ParseFiles(templatePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while parsing the template in %s: %w", templatePath, err)
// 	}

// 	return &Template{
// 		HtmlTpl: tpl,
// 	}, nil
// }
