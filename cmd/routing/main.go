package main

import (
	"fmt"
	"net/http"
	"webGo/controllers"

	"github.com/go-chi/chi/v5"
)

func main() {
	myRouter := chi.NewRouter()
	// Set the routes in the router object.
	controllers.GetAll(myRouter)
	fmt.Println("Starting the server on: 3000...")
	http.ListenAndServe(":3000", myRouter)
}
