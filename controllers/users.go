package controllers

import (
	"fmt"
	"log"
	"net/http"
	"webGo/models"
)

type Users struct {
	// Template struct we'll store all of the needed templaes in.
	Templates struct {
		New    Template
		SignIn Template
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
		email := r.FormValue("email")
		password := r.FormValue("password")
		// The open db connection is passed inside UserService
		userInfo, err := u.UserService.Create(email, password)
		if err != nil {
			log.Println(err)
			http.Error(w,
				"Error when creating a new user entry.",
				http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "User created: %+v", userInfo)
	}
}

// Used as a handler function for GET request when
// getting the template of signing in a new user.
func (u Users) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Email string
		}
		data.Email = r.FormValue("email")
		// Putting the template object with the data inside the Users info.
		u.Templates.SignIn.Execute(w, data)
	}
}
