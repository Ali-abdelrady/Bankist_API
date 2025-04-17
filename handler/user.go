package handler

import (
	"bankist/db"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
    FullName string `json:"full_name"`
}
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Message string `json:"message"`
    UserID  int    `json:"user_id"`
}
func SignUp(w http.ResponseWriter, r *http.Request) {
    var req SignupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    _, err = db.DB.Exec(`
        INSERT INTO users (username, password_hash, full_name)
        VALUES ($1, $2, $3)
    `, req.Username, hashed, req.FullName)

    if err != nil {
        log.Println("Error inserting user:", err)
        http.Error(w, "Could not create user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}


func Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    var storedHash string
    var userID int
    err := db.DB.QueryRow(`
        SELECT user_id, password_hash FROM users WHERE username = $1 AND is_active = TRUE
    `, req.Username).Scan(&userID, &storedHash)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Compare password with bcrypt hash
    if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(LoginResponse{
        Message: "Login successful",
        UserID:  userID,
    })
}