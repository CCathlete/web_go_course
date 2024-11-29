package controllers

import (
	"fmt"
	"net/http"
	"webGo/models"
)

func StaticGetHandler(tpl Template, data any, routeSuffix string) http.HandlerFunc {
	// Preparing the template for the response to GET requests.
	usersC := NewUserController()

	switch routeSuffix {

	case "faq":
		return FAQ(tpl)

	case "signup":
		// Assigning the template into the field for new user creation.
		usersC.Templates.New = tpl
		// Getting the template of the sign up page.
		return usersC.New()

	case "signin":
		// Assigning the template into the field for sign in.
		usersC.Templates.SignIn = tpl
		// Getting the template of the sign in page.
		return usersC.SignIn()

	case "users/me":
		// Returning info about the current user (from cookies).
		return usersC.CurrentUser()

	default:
		return func(w http.ResponseWriter, r *http.Request) {
			tpl.Execute(w, data)
		}
	}
}

func FAQ(tpl Template) http.HandlerFunc {
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

func StaticPostHandler(tpl Template, routeSuffix string) http.HandlerFunc {
	usersC := NewUserController()

	switch routeSuffix {

	case "users":
		// Sending back the form (insite the signup template)
		// with the information for a creation of a new user.
		db, err := models.ConnectToDB()
		if err != nil {
			return func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, fmt.Sprintf("Error connetcing to DB: %v", err), http.StatusInternalServerError)
			}
		}

		// Preparing the template and db connection for POST requests.
		usersC.Templates.New = tpl
		usersC.UserService.DB = db
		// Using the data we got from the POST request.
		return usersC.Create()

	default:
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("POST request doing something.")
		}
	}
}

/*
	case "faq":
		// var err error
		// data, err = getQuestionsTemplate("QA.yaml")
		// if err != nil {
		// 	return func(w http.ResponseWriter, r *http.Request) {
		// 		http.Error(w, fmt.Sprintf("Error while parsing the questions yaml: %v", err), http.StatusInternalServerError)
		// 	}
		// }
		return FAQ(tpl)
*/
