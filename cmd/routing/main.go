package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"webGo/controllers"
	"webGo/views"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	myRouter := chi.NewRouter()
	// Set the routes in the router object.
	controllers.RouteAll(myRouter)
	fmt.Println("Starting the server on: 3000...")
	csrfMW := prepMiddleWare(false) //TODO: change this before deploying.
	http.ListenAndServe(":3000", csrfMW(myRouter))
}

// Creats a CSRF key and returns a middleware wrapper for
// an http router (http.Handler)
func prepMiddleWare(activate bool) func(http.Handler) http.Handler {
	csrfKey := make([]byte, 32)
	views.Must(rand.Read(csrfKey)) // Reads the content of a random num generator and writes it into the underlying array of the byteslice.
	csrfMW := csrf.Protect(
		csrfKey,
		csrf.Secure(activate),
	)
	return csrfMW
}
