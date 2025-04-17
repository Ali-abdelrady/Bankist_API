package models

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

type Operation struct {
	OperationID   int     `json:"operation_id"`
	UserID        int     `json:"user_id"`
	RelatedUserID *int    `json:"related_user_id,omitempty"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	CreatedAt     string  `json:"created_at"`
}
