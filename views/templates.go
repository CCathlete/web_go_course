package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Template struct {
	HtmlTpl *template.Template
}

func ParseTempate(w http.ResponseWriter, templatePath string) *Template {
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Error when parsing the tamplate: %v", err)
		http.Error(w, fmt.Sprintf("Error when parsing template: %v", err), http.StatusInternalServerError)
		return nil
	}

	return &Template{
		HtmlTpl: tpl,
	}
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
