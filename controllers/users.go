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
	UserService    *models.UserService
	SessionService *models.SessionService
}

// Returns an initialised users controller.
func NewUserController() Users {
	return Users{
		UserService:    &models.UserService{},
		SessionService: &models.SessionService{},
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
		u.Templates.New.Execute(w, r, data)
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
			return
		}
		session, err := u.SessionService.Create(userInfo.ID)
		if err != nil {
			log.Println(err)
			http.Error(w,
				"Error when creating a new session for the new user.",
				http.StatusInternalServerError)
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		// Creating a cookie with the session token and setting it in the response.
		http.SetCookie(w, NewCookie(CookieSession, session.Token))
		http.Redirect(w, r, "/users/me", http.StatusFound)
	}
}

// Used as a handler function for GET request when
// getting the template of signing in an existing user.
func (u Users) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Email string
		}
		data.Email = r.FormValue("email")
		// Putting the template object with the data inside the Users info.
		u.Templates.SignIn.Execute(w, r, data)
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

		session, err := u.SessionService.Create(userInfo.ID)
		if err != nil {
			log.Println(err)
			http.Error(w,
				"Error when creating a new session for the user.",
				http.StatusInternalServerError)
			return
		}

		// Creating a cookie with the session token and setting it in the response.
		http.SetCookie(w, NewCookie(CookieSession, session.Token))
		http.Redirect(w, r, "/users/me", http.StatusFound)
	}
}

// Takes up a web requests and prints put the current user information.
func (u Users) CurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(CookieSession)
		if err != nil {
			log.Println(err)
			http.Error(w,
				"No active session, please sign in.",
				http.StatusInternalServerError)
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		// We're passing the user service in because we want to query
		// the DB where the users are in.
		user, err := u.SessionService.User(u.UserService, tokenCookie.Value)
		if err != nil {
			log.Println(err)
			http.Error(w,
				"No active session, please sign in.",
				http.StatusInternalServerError)
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		fmt.Fprintf(w, "Current user: %v\n", user)
		// fmt.Fprintf(w, "Current user: %s\n", user.Email)
	}
}
