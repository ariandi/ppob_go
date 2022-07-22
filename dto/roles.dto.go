package dto

import "database/sql"

type Role struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Level     int16         `json:"level"`
	CreatedAt sql.NullTime  `json:"created_at"`
	UpdatedAt sql.NullTime  `json:"updated_at"`
	DeletedAt sql.NullTime  `json:"deleted_at"`
	CreatedBy sql.NullInt64 `json:"created_by"`
	UpdatedBy sql.NullInt64 `json:"updated_by"`
	DeletedBy sql.NullInt64 `json:"deleted_by"`
}

type CreateRoleReq struct {
	Name  string `json:"name" binding:"required"`
	Level int16  `json:"level" binding:"required,min=1"`
}

type RoleRes struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Level int16  `json:"level"`
}

type GetRoleReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateRoleRequest struct {
	ID    int64  `uri:"id" binding:"required,min=1"`
	Name  string `json:"name" binding:"required"`
	Level int16  `json:"level" binding:"required,min=1"`
}

type ListRoleRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type UpdateInactiveRoleRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
