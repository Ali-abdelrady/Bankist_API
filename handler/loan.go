package handler

import (
	"bankist/db"
	"encoding/json"
	"net/http"
)

type LoanRequest struct {
    UserID int     `json:"user_id"`
    Amount float64 `json:"amount"`
}

func RequestLoan(w http.ResponseWriter, r *http.Request) {
    var req LoanRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    _, err := db.DB.Exec(`
        INSERT INTO operations (user_id, type, amount)
        VALUES ($1, 'loan', $2)
    `, req.UserID, req.Amount)

    if err != nil {
        http.Error(w, "Loan request failed", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Loan granted"})
}
