package models

import "time"

// aqui vai servir para transferir os dados para o postgres
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

// aqui vai servir para o usuario transferir os dados para o sistema como DTO
type CreateUserRequest struct {
	Username string
	Password string
}
