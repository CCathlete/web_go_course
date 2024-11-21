package views

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Template struct {
	HtmlTpl *template.Template
}

func ParseTemplate(templatePath string) (*Template, error) {
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error while parsing the template in %s: %w", templatePath, err)
	}

	return &Template{
		HtmlTpl: tpl,
	}, nil
}

func ParseFS(fs embed.FS, pattern ...string) (*Template, error) {
	tpl, err := template.ParseFS(fs, pattern...)
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

func (tpl *Template) Execute(w http.ResponseWriter, data any) {
	// Writing the html page into a buffer to make sure we don't have an error
	// before writing to the respinse writer.
	var actualRes bytes.Buffer
	if err := tpl.HtmlTpl.Execute(&actualRes, data); err != nil {
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
