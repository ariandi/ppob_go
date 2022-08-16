package dto

type CreateProviderReq struct {
	Name      string `json:"name" binding:"required"`
	User      string `json:"user" binding:"required"`
	Secret    string `json:"secret" binding:"required"`
	AddInfo1  string `json:"add_info1"`
	AddInfo2  string `json:"add_info2"`
	ValidFrom string `json:"valid_from"`
	ValidTo   string `json:"valid_to"`
	BaseUrl   string `json:"base_url"`
	Method    string `json:"method"`
	Inq       string `json:"inq"`
	Pay       string `json:"pay"`
	Adv       string `json:"adv"`
	Cmt       string `json:"cmt"`
	Rev       string `json:"rev"`
	Status    string `json:"status"`
}

type ProviderRes struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Secret    string `json:"secret"`
	AddInfo1  string `json:"add_info1"`
	AddInfo2  string `json:"add_info2"`
	ValidFrom string `json:"valid_from"`
	ValidTo   string `json:"valid_to"`
	BaseUrl   string `json:"base_url"`
	Method    string `json:"method"`
	Inq       string `json:"inq"`
	Pay       string `json:"pay"`
	Adv       string `json:"adv"`
	Cmt       string `json:"cmt"`
	Rev       string `json:"rev"`
	Status    string `json:"status"`
}

type GetProviderReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListProviderRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type UpdateProviderRequest struct {
	ID        int64  `uri:"id" binding:"required,min=1"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Secret    string `json:"secret"`
	AddInfo1  string `json:"add_info1"`
	AddInfo2  string `json:"add_info2"`
	ValidFrom string `json:"valid_from"`
	ValidTo   string `json:"valid_to"`
	BaseUrl   string `json:"base_url"`
	Method    string `json:"method"`
	Inq       string `json:"inq"`
	Pay       string `json:"pay"`
	Adv       string `json:"adv"`
	Cmt       string `json:"cmt"`
	Rev       string `json:"rev"`
	Status    string `json:"status"`
}

type UpdateInactiveProviderRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
