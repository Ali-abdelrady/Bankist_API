package handler

import (
	"bankist/db"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var balance float64
    query := "SELECT get_balance($1)"
    err = db.DB.QueryRow(query, userID).Scan(&balance)
    if err != nil {
        http.Error(w, "Error getting balance", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}
