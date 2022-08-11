package dto

type CreateCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

type CategoryRes struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetCategoryReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateCategoryRequest struct {
	ID   int64  `uri:"id" binding:"required,min=1"`
	Name string `json:"name" binding:"required"`
}

type ListCategoryRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type UpdateInactiveCategoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
