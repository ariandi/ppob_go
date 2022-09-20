package dto

type InqRequest struct {
	TimeStamp     string `json:"time_stamp" binding:"required"`
	UserID        string `json:"user_id" binding:"required"`
	RefID         string `json:"ref_id" binding:"required"`
	BillID        string `json:"bill_id" binding:"required"`
	AppName       string `json:"app_name" binding:"required"`
	ProductCode   string `json:"product_code" binding:"required"`
	MerchantToken string `json:"merchant_token" binding:"required"`
	Amount        int64  `json:"amount"`
}

type PayRequest struct {
	TimeStamp     string `json:"time_stamp" binding:"required"`
	UserID        string `json:"user_id" binding:"required"`
	RefID         string `json:"ref_id" binding:"required"`
	BillID        string `json:"bill_id" binding:"required"`
	AppName       string `json:"app_name" binding:"required"`
	ProductCode   string `json:"product_code" binding:"required"`
	MerchantToken string `json:"merchant_token" binding:"required"`
	Amount        int64  `json:"amount"`
	Admin         int64  `json:"admin"`
	TotalAmount   int64  `json:"total_amount"`
	TxID          string `json:"tx_id" binding:"required"`
}

type InqResponse struct {
	TimeStamp     string `json:"time_stamp"`
	UserID        string `json:"user_id"`
	RefID         string `json:"ref_id"`
	BillID        string `json:"bill_id"`
	AppName       string `json:"app_name"`
	ProductCode   string `json:"product_code"`
	MerchantToken string `json:"merchant_token"`
	ProductName   string `json:"product_name"`
	Amount        int64  `json:"amount"`
	Admin         int64  `json:"admin"`
	TotalAmount   int64  `json:"total_amount"`
	ResultCd      string `json:"result_cd"`
	ResultMsg     string `json:"result_msg"`
	TxID          string `json:"tx_id"`
}

type PayResponse struct {
	TimeStamp     string `json:"time_stamp"`
	UserID        string `json:"user_id"`
	RefID         string `json:"ref_id"`
	BillID        string `json:"bill_id"`
	AppName       string `json:"app_name"`
	ProductCode   string `json:"product_code"`
	MerchantToken string `json:"merchant_token"`
	ProductName   string `json:"product_name"`
	Amount        int64  `json:"amount"`
	Admin         int64  `json:"admin"`
	TotalAmount   int64  `json:"total_amount"`
	ResultCd      string `json:"result_cd"`
	ResultMsg     string `json:"result_msg"`
	TxID          string `json:"tx_id"`
}

type DepositRequest struct {
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Name    string `json:"name" binding:"required"`
	AppName string `json:"app_name" binding:"required"`
	Amount  int64  `json:"amount" binding:"required"`
}

type DepositApproveRequest struct {
	TxID string `json:"tx_id" binding:"required"`
}

type DepositResponse struct {
	ResultCd  string `json:"result_cd"`
	ResultMsg string `json:"result_msg"`
	TxID      string `json:"tx_id"`
}

type InqRequestConsume struct {
	InqRequest  InqRequest
	InqResponse InqResponse
	PayRequest  PayRequest
	PayResponse PayResponse
	QueueName   string
}

type DepositRequestConsume struct {
	DepositRequest        DepositRequest
	DepositApproveRequest DepositApproveRequest
	DepositResponse       DepositResponse
	UserRequest           UserResponse
	QueueName             string
}

type InqSetResponse struct {
	InqData     InqRequest
	ProductName string `json:"product_name"`
	Amount      int64  `json:"amount"`
	Admin       int64  `json:"admin"`
	TotalAmount int64  `json:"total_amount"`
	ResultCd    string `json:"result_cd"`
	ResultMsg   string `json:"result_msg"`
	TxID        string `json:"tx_id"`
}
