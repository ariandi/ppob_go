package dto

type CreateTransactionReq struct {
	TxID         string `json:"tx_id" binding:"required"`
	BillID       string `json:"bill_id"`
	CustName     string `json:"cust_name"`
	Amount       string `json:"amount"`
	Admin        string `json:"admin"`
	TotAmount    string `json:"tot_amount"`
	FeePartner   string `json:"fee_partner"`
	FeePpob      string `json:"fee_ppob"`
	ValidFrom    string `json:"valid_from"`
	ValidTo      string `json:"valid_to"`
	CatID        int64  `json:"cat_id"`
	CatName      string `json:"cat_name"`
	ProdID       int64  `json:"prod_id"`
	ProdName     string `json:"prod_name"`
	PartnerID    int64  `json:"partner_id"`
	PartnerName  string `json:"partner_name"`
	ProviderID   int64  `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	Status       string `json:"status"`
	PaymentType  string `json:"payment_type"`
	ReqInqParams string `json:"req_inq_params"`
	ResInqParams string `json:"res_inq_params"`
	ReqPayParams string `json:"req_pay_params"`
	ResPayParams string `json:"res_pay_params"`
	ReqCmtParams string `json:"req_cmt_params"`
	ResCmtParams string `json:"res_cmt_params"`
	ReqAdvParams string `json:"req_adv_params"`
	ResAdvParams string `json:"res_adv_params"`
	ReqRevParams string `json:"req_rev_params"`
	ResRevParams string `json:"res_rev_params"`
}

type TransactionRes struct {
	ID           int64  `json:"id"`
	TxID         string `json:"tx_id"`
	BillID       string `json:"bill_id"`
	CustName     string `json:"cust_name"`
	Amount       string `json:"amount"`
	Admin        string `json:"admin"`
	TotAmount    string `json:"tot_amount"`
	FeePartner   string `json:"fee_partner"`
	FeePpob      string `json:"fee_ppob"`
	FirstBalance string `json:"first_balance"`
	LastBalance  string `json:"last_balance"`
	Sn           string `json:"sn"`
	ValidFrom    string `json:"valid_from"`
	ValidTo      string `json:"valid_to"`
	CatID        int64  `json:"cat_id"`
	CatName      string `json:"cat_name"`
	ProdID       int64  `json:"prod_id"`
	ProdName     string `json:"prod_name"`
	PartnerID    int64  `json:"partner_id"`
	PartnerName  string `json:"partner_name"`
	ProviderID   int64  `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	Status       string `json:"status"`
	PaymentType  string `json:"payment_type"`
	CreatedBy    string `json:"created_by"`
	ReqInqParams string `json:"req_inq_params"`
	ResInqParams string `json:"res_inq_params"`
	ReqPayParams string `json:"req_pay_params"`
	ResPayParams string `json:"res_pay_params"`
	ReqCmtParams string `json:"req_cmt_params"`
	ResCmtParams string `json:"res_cmt_params"`
	ReqAdvParams string `json:"req_adv_params"`
	ResAdvParams string `json:"res_adv_params"`
	ReqRevParams string `json:"req_rev_params"`
	ResRevParams string `json:"res_rev_params"`
}

type GetTransactionReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type GetTransactionByTxIDReq struct {
	ID string `uri:"tx_id" binding:"required,min=1"`
}

type GetTransactionCountReq struct {
	Status   string `form:"status" binding:"required"`
	FromDate string `form:"from_date" binding:"required"`
	ToDate   string `form:"to_date" binding:"required"`
}

type ListTransactionRequest struct {
	PageID      int32  `form:"page_id" binding:"required,min=1"`
	PageSize    int32  `form:"page_size" binding:"required,min=5,max=500"`
	FromDate    string `form:"from_date" binding:"required"`
	ToDate      string `form:"to_date" binding:"required"`
	CatID       int64  `form:"cat_id"`
	PartnerID   int64  `form:"partner_id"`
	CreatedBy   int64  `form:"created_by"`
	Status      string `form:"status"`
	PaymentType string `form:"payment_type"`
}

type UpdateTransactionRequest struct {
	TxID         int64  `uri:"tx_id" binding:"required,min=1"`
	Status       string `json:"status"`
	ReqPayParams string `json:"req_pay_params"`
	ResPayParams string `json:"res_pay_params"`
	ReqCmtParams string `json:"req_cmt_params"`
	ResCmtParams string `json:"res_cmt_params"`
	ReqAdvParams string `json:"req_adv_params"`
	ResAdvParams string `json:"res_adv_params"`
	ReqRevParams string `json:"req_rev_params"`
	ResRevParams string `json:"res_rev_params"`
}

type UpdateInactiveTransactionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
