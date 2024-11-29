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

// Returns an initialised users controller.
func NewUserController() Users {
	return Users{
		UserService: &models.UserService{},
	}
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
// getting the template of signing in an exsting user.
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

// Used as a handler function for POST request for
// authentication and processing the data of an existing user.
func (u Users) ProcessSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Email, Password string
		}
		data.Email = r.FormValue("email")
		data.Password = r.FormValue("password")
		fmt.Println("")
		log.Printf("Original credentials: \n%+v\n\n", data)
		userInfo, err := u.UserService.Authenticate(data.Email, data.Password)
		if err != nil {
			log.Println("ProcessSignIn: ", err)
			http.Error(w, "Authentication error.", http.StatusUnauthorized)
			return
		}

		// FOR FUTURE USE:
		// Putting the template object with the data inside the Users info.
		// u.Templates.SignIn.Execute(w, data)

		fmt.Fprintf(w, "Authentication successful: \n%+v", userInfo)
	}
}

// Takes up a web requests and prints put the current user information.
func (u Users) CurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, err := r.Cookie("email")
		if err != nil {
			fmt.Fprint(w, "The email cookie couldn't be read.")
			return
		}

		fmt.Fprintf(w, "Email cookie: %s\n", email.Value)
		fmt.Fprintf(w, "Headers: %+v\n", r.Header)
	}
}
