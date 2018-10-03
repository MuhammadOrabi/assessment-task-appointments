// Entrypoint for API
package main

import (
 	"log"
 	"net/http"
 	"os"
	"github.com/gorilla/handlers"
	"app/appointments"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := appointments.NewRouter() // create routes

	// These two lines are important in order to allow access from the front-end side to the methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"}) 
 	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	// Launch server with CORS validations
 	log.Fatal(http.ListenAndServe("0.0.0.0:" + port, handlers.CORS(allowedOrigins, allowedMethods)(router)))
}