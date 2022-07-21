package dto

import "database/sql"

type RoleUser struct {
	ID        int64         `json:"id"`
	RoleID    int64         `json:"role_id"`
	UserID    int64         `json:"user_id"`
	CreatedAt sql.NullTime  `json:"created_at"`
	UpdatedAt sql.NullTime  `json:"updated_at"`
	DeletedAt sql.NullTime  `json:"deleted_at"`
	CreatedBy sql.NullInt64 `json:"created_by"`
	UpdatedBy sql.NullInt64 `json:"updated_by"`
	DeletedBy sql.NullInt64 `json:"deleted_by"`
}

type CreateRoleUserReq struct {
	RoleID int64 `json:"role_id" binding:"required,min=1"`
	UserID int64 `json:"user_id" binding:"required,min=1"`
}

type CreateRoleUserRes struct {
	ID     int64 `json:"id"`
	RoleID int64 `json:"role_id" binding:"required,min=1"`
	UserID int64 `json:"user_id" binding:"required,min=1"`
}

type GetRoleUserByUserID struct {
	UserID int64 `uri:"user_id" binding:"required,min=1"`
}

type GetRoleUserByRoleID struct {
	RoleID int64 `uri:"role_id" binding:"required,min=1"`
}
