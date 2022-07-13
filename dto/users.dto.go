package dto

import "database/sql"

type CreateUserRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Username       string `json:"username" binding:"required,alphanum"`
	Password       string `json:"password" binding:"min=6"`
	Balance        string `json:"balance"`
	Phone          string `json:"phone"`
	IdentityNumber string `json:"identity_number" binding:"required"`
	CreatedBy      int64  `json:"created_by"`
}

type CreateUserResponse struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Username       string         `json:"username"`
	Balance        sql.NullString `json:"balance"`
	Phone          string         `json:"phone"`
	IdentityNumber string         `json:"identity_number"`
}
