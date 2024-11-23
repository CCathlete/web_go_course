package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	// Template struct we'll store all of the needed templaes in.
	Templates struct {
		New Template
	}
	Data any
}

func (u Users) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u.Templates.New.Execute(w, u.Data)
	}
}

func (u Users) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Email: %s\n", r.FormValue("email"))
		fmt.Fprintf(w, "Password: %s\n", r.FormValue("password"))
	}
}
