package main

import (
	"bankist/db"
	"bankist/handler"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Database info
// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "Rady_2003"
// 	dbname   = "Bankist"
// )

func main() {
	connection := "postgresql://postgres:XUTGJtbFbuAHKTLHNTXvshYFdFXwPiuT@postgres.railway.internal:5432/railway"
    if connection == "" {
        log.Fatal("DATABASE_URL is not set")
    }
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

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
