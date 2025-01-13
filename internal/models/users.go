package models

type User struct {
	ID        string  `json:"user_id"    db:"user_id"`
	Login     string  `json:"login"      db:"login"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
	DeletedAt *string `json:"deleted_at" db:"deleted_at"`
}
