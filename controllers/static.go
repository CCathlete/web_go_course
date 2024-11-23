package controllers

import (
	"net/http"
	"webGo/views"
)

func StaticHandler(tpl *views.Template, routeSuffix string) http.HandlerFunc {
	// var data *template.HTML
	var data interface{}
	switch routeSuffix {
	case "faq":
		// var err error
		// data, err = getQuestionsTemplate("QA.yaml")
		// if err != nil {
		// 	return func(w http.ResponseWriter, r *http.Request) {
		// 		http.Error(w, fmt.Sprintf("Error while parsing the questions yaml: %v", err), http.StatusInternalServerError)
		// 	}
		// }
		return FAQ(tpl)
	case "signup":
		usersC := Users{
			Templates: struct {
				New Template
			}{
				tpl,
			},
		}
		return usersC.New()
	default:
		return func(w http.ResponseWriter, r *http.Request) {
			tpl.Execute(w, data)
		}
	}
}

func FAQ(tpl *views.Template) http.HandlerFunc {
	questions := []struct{ Question, Answer string }{
		{
			Question: "Alternative Q1",
			Answer:   "Alternative A1",
		},
		{
			Question: "Alternative Q2",
			Answer:   "Alternative A2",
		},
		{
			Question: "Alternative Q3",
			Answer:   "Alternative A3",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
