package dto

type CreateSellingReq struct {
	PartnerID  int64  `json:"partner_id" binding:"required"`
	CategoryID int64  `json:"category_id" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}

type SellingRes struct {
	ID         int64  `json:"id"`
	PartnerID  int64  `json:"partner_id" binding:"required"`
	CategoryID int64  `json:"category_id" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}

type GetSellingReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListSellingRequest struct {
	PageID     int32 `form:"page_id" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=5,max=500"`
	PartnerID  int64 `form:"partner_id"  binding:"required,min=1"`
	CategoryID int64 `form:"category_id"  binding:"required,min=1"`
}

type UpdateSellingRequest struct {
	ID         int64  `uri:"id" binding:"required,min=1"`
	PartnerID  int64  `json:"partner_id" binding:"required"`
	CategoryID int64  `json:"category_id" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}

type UpdateInactiveSellingRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
