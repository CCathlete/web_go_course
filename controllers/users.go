package controllers

import (
	"net/http"
	"webGo/views"
)

type Users struct {
	// Template struct we'll store all of the needed templaes in.
	Templates struct {
		New *views.Template
	}
	Data any
}

func (u Users) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u.Templates.New.Execute(w, u.Data)
	}
}
