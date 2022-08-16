package dto

type CreateProductReq struct {
	Name       string `json:"name" binding:"required"`
	CatID      int64  `json:"cat_id" binding:"required"`
	Amount     string `json:"amount"`
	ProviderID int64  `json:"provider_id" binding:"required"`
	Status     string `json:"status"`
	Parent     int64  `json:"parent"`
}

type ProductRes struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CatID      int64  `json:"cat_id"`
	Amount     string `json:"amount"`
	ProviderID int64  `json:"provider_id"`
	Status     string `json:"status"`
	Parent     int64  `json:"parent"`
}

type GetProductReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListProductRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type UpdateProductRequest struct {
	ID         int64  `uri:"id" binding:"required,min=1"`
	Name       string `json:"name" binding:"required"`
	CatID      int64  `json:"cat_id" binding:"required"`
	Amount     string `json:"amount"`
	ProviderID int64  `json:"provider_id" binding:"required"`
	Status     string `json:"status"`
	Parent     int64  `json:"parent"`
}

type UpdateInactiveProductRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
