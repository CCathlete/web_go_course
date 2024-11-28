package controllers

import (
	"fmt"
	"net/http"
	"webGo/models"
)

type Users struct {
	// Template struct we'll store all of the needed templaes in.
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

// Used as a handler function for GET request when
// getting the template of signing up a new user.
func (u Users) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Email string
		}
		data.Email = r.FormValue("email")
		// Putting the template object with the data inside the Users info.
		u.Templates.New.Execute(w, data)
	}
}

// Used as a handler function for POST request when
// creating a new user.
func (u Users) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Email: %s\n", r.FormValue("email"))
		fmt.Fprintf(w, "Password: %s\n", r.FormValue("password"))
	}
}
