package main

import (
	"bankist/db"
	"bankist/handler"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)


func main() {
	connection := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	db.InitDB(connection)

	r := mux.NewRouter()

	r.HandleFunc("/api/users/{id}/balance", handler.GetBalance).Methods("GET")
	r.HandleFunc("/api/users/{id}/history", handler.GetHistory).Methods("GET")
	r.HandleFunc("/api/transfer", handler.TransferMoney).Methods("POST")
	r.HandleFunc("/api/loan", handler.RequestLoan).Methods("POST")
	r.HandleFunc("/api/users", handler.SignUp).Methods("POST")
	r.HandleFunc("/api/login", handler.Login).Methods("POST")

	// CORS configuration
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"}) // Allow all origins (for dev)
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	log.Println("Server is running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
