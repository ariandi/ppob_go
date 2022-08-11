package dto

type CreatePartnerReq struct {
	Name        string `json:"name" binding:"required"`
	User        string `json:"user" binding:"required"`
	Secret      string `json:"secret" binding:"required"`
	AddInfo1    string `json:"add_info_1"`
	AddInfo2    string `json:"add_info_2"`
	ValidFrom   string `json:"valid_from"`
	ValidTo     string `json:"valid_to"`
	PaymentType string `json:"payment_type"`
	Status      string `json:"status"`
}

type PartnerRes struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	User        string `json:"user"`
	Secret      string `json:"secret"`
	AddInfo1    string `json:"add_info_1"`
	AddInfo2    string `json:"add_info_2"`
	ValidFrom   string `json:"valid_from"`
	ValidTo     string `json:"valid_to"`
	PaymentType string `json:"payment_type"`
	Status      string `json:"status"`
}

type GetPartnerReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdatePartnerRequest struct {
	ID          int64  `uri:"id" binding:"required,min=1"`
	Name        string `json:"name"`
	User        string `json:"user"`
	Secret      string `json:"secret"`
	AddInfo1    string `json:"add_info_1"`
	AddInfo2    string `json:"add_info_2"`
	ValidFrom   string `json:"valid_from"`
	ValidTo     string `json:"valid_to"`
	PaymentType string `json:"payment_type" binding:"required,paymentType"`
	Status      string `json:"status"`
}

type ListPartnerRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type UpdateInactivePartnerRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
