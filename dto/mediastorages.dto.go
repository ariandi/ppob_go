package dto

import (
	"time"
)

type CreateMediaReq struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	SecID   string `json:"sec_id"`
	TabID   string `json:"tab_id"`
}

type MediaRes struct {
	ID        int64     `json:"id"`
	SecID     string    `json:"sec_id"`
	TabID     string    `json:"tab_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
}

type GetMediaReq struct {
	ID    int64  `form:"id"`
	SecID string `form:"sec_id"`
	TabID string `form:"tab_id"`
}

type UpdateMediaRequest struct {
	ID   int64  `uri:"id" binding:"required,min=1"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ListMediaRequest struct {
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=50"`
	SecID    string `form:"sec_id"`
	TabID    string `form:"tab_id"`
}

type UpdateInactiveMediaRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
