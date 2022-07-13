package dto

import "database/sql"

type CreateUserRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Username       string `json:"username" binding:"required,alphanum"`
	Password       string `json:"password" binding:"min=6"`
	Phone          string `json:"phone"`
	IdentityNumber string `json:"identity_number" binding:"required"`
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

type UpdateUserRequest struct {
	ID             int64  `uri:"id" binding:"required,min=1"`
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Username       string `json:"username" binding:"required,alphanum"`
	Password       string `json:"password" binding:"min=6"`
	Balance        string `json:"balance"`
	Phone          string `json:"phone"`
	IdentityNumber string `json:"identity_number" binding:"required"`
}

type GetUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}
