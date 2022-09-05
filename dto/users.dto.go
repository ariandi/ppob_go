package dto

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type CreateUserRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Username       string `json:"username" binding:"required,alphanum"`
	Password       string `json:"password" binding:"min=6"`
	Phone          string `json:"phone"`
	IdentityNumber string `json:"identity_number" binding:"required"`
	RoleID         int64  `json:"role_id" binding:"required"`
}

type UserResponse struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Username       string         `json:"username"`
	Balance        sql.NullString `json:"balance"`
	Phone          string         `json:"phone"`
	IdentityNumber string         `json:"identity_number"`
	BankCode       int64          `json:"bank_code"`
	Role           []RoleUser     `json:"role"`
}

type UpdateUserRequest struct {
	ID             int64
	Name           string `json:"name"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Balance        string `json:"balance"`
	Phone          string `json:"phone"`
	IdentityNumber string `json:"identity_number"`
	RoleID         int64  `json:"role_id" binding:"required"`
}

type UpdateUserIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateInactiveUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type GetUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=200"`
	RoleID   int64 `form:"role_id"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}
