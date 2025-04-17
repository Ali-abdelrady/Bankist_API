package handler

import (
	"bankist/db"
	"bankist/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetHistory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    rows, err := db.DB.Query(`
        SELECT operation_id, user_id, related_user_id, type, amount, created_at
        FROM operations
        WHERE user_id = $1
        ORDER BY created_at DESC
    `, userID)
    if err != nil {
        http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var ops []models.Operation
    for rows.Next() {
        var op models.Operation
        err := rows.Scan(&op.OperationID, &op.UserID, &op.RelatedUserID, &op.Type, &op.Amount, &op.CreatedAt)
        if err != nil {
            http.Error(w, "Error reading result", http.StatusInternalServerError)
            return
        }
        ops = append(ops, op)
    }
    w.Header().Set("Content-Type","application/json")
    json.NewEncoder(w).Encode(ops)
}
