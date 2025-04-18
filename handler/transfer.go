package handler

import (
	"bankist/db"
	"encoding/json"
	"net/http"
)

type TransferRequest struct {
	FromUserID int     `json:"from_user_id"`
	ToUsername string  `json:"to_username"` // Changed from ToUserID to ToUsername
	Amount     float64 `json:"amount"`
}

func TransferMoney(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 🔍 Get ToUserID from username
	var toUserID int
	err := db.DB.QueryRow(`SELECT user_id FROM users WHERE username = $1`, req.ToUsername).Scan(&toUserID)
	if err != nil {
		http.Error(w, "Recipient user not found", http.StatusNotFound)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// 🧾 Withdraw from sender
	_, err = tx.Exec(`
        INSERT INTO operations (user_id, related_user_id, type, amount)
        VALUES ($1, $2, 'transfer_out', $3)
    `, req.FromUserID, toUserID, req.Amount)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Transfer failed (sender)", http.StatusInternalServerError)
		return
	}

	// 💰 Deposit to receiver
	_, err = tx.Exec(`
        INSERT INTO operations (user_id, related_user_id, type, amount)
        VALUES ($1, $2, 'transfer_in', $3)
    `, toUserID, req.FromUserID, req.Amount)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Transfer failed (receiver)", http.StatusInternalServerError)
		return
	}

	tx.Commit()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Transfer completed"})
}