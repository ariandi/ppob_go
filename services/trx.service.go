package services

import (
	"database/sql"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type TransactionInterface interface {
	CreateTransactionService(ctx *gin.Context, in dto.CreateTransactionReq) (dto.TransactionRes, error)
	GetTransactionService(ctx *gin.Context, in dto.GetTransactionByTxIDReq) (dto.TransactionRes, error)
	ListTransactionService(ctx *gin.Context, in dto.ListTransactionRequest) ([]dto.TransactionRes, error)
	GetTransactionCountService(ctx *gin.Context, in dto.GetTransactionCountReq) (db.GetTransactionCountRow, error)
	UpdateTransactionService(ctx *gin.Context, in dto.UpdateTransactionRequest) (dto.TransactionRes, error)
	SoftDeleteTransactionService(ctx *gin.Context, in dto.UpdateInactiveTransactionRequest) error
	setCreateParams(arg db.CreateTransactionParams, in dto.CreateTransactionReq) db.CreateTransactionParams
	setUpdateParams(arg db.UpdateTransactionParams, in dto.UpdateTransactionRequest) db.UpdateTransactionParams
	TransactionRes(trx db.Transaction) dto.TransactionRes
	ExportTransaction(ctx *gin.Context, in dto.ListTransactionRequest) (*excelize.File, error)
	getTransactionList(ctx *gin.Context, in dto.ListTransactionRequest) ([]db.Transaction, error)
}

// TransactionService is
type TransactionService struct {
	store db.Store
}

var transactionService *TransactionService

// GetTransactionService is
func GetTransactionService(store db.Store) TransactionInterface {
	if transactionService == nil {
		transactionService = &TransactionService{
			store: store,
		}
	}
	return transactionService
}

func (o *TransactionService) CreateTransactionService(ctx *gin.Context, in dto.CreateTransactionReq) (dto.TransactionRes, error) {
	logrus.Println("[ProductService CreateTransactionService] start.")
	var out dto.TransactionRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.CreateTransactionParams{
		TxID:      in.TxID,
		BillID:    in.BillID,
		Status:    in.Status,
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg = o.setCreateParams(arg, in)

	transaction, err := o.store.CreateTransaction(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return out, err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.TransactionRes(transaction)
	return out, nil
}

func (o *TransactionService) GetTransactionService(ctx *gin.Context, in dto.GetTransactionByTxIDReq) (dto.TransactionRes, error) {
	logrus.Println("[TransactionService GetTransactionService] start.")
	var out dto.TransactionRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	trx, err := o.store.GetTransactionByTxID(ctx, in.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.TransactionRes(trx)
	return out, nil
}

func (o *TransactionService) GetTransactionCountService(ctx *gin.Context, in dto.GetTransactionCountReq) (db.GetTransactionCountRow, error) {
	logrus.Println("[TransactionService GetTransactionCountService] start.")
	var out db.GetTransactionCountRow

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	fromDt, err := time.Parse("2006-01-02", in.FromDate)
	if err != nil {
		return out, errors.New("from date format error")
	}

	toDt, err := time.Parse("2006-01-02", in.ToDate)
	if err != nil {
		return out, errors.New("to date format error")
	}
	args := db.GetTransactionCountParams{
		IsStatus: true,
		Status:   in.Status,
		Fromdt: sql.NullTime{
			Time:  fromDt,
			Valid: true,
		},
		Todt: sql.NullTime{
			Time:  toDt,
			Valid: true,
		},
	}
	out, err = o.store.GetTransactionCount(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	return out, nil
}

func (o *TransactionService) ListTransactionService(ctx *gin.Context, in dto.ListTransactionRequest) ([]dto.TransactionRes, error) {
	logrus.Println("[TransactionService ListTransactionService] start.")
	var out []dto.TransactionRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	trxList, err := o.getTransactionList(ctx, in)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	for _, trx := range trxList {
		u := o.TransactionRes(trx)
		out = append(out, u)
	}

	return out, nil
}

func (o *TransactionService) UpdateTransactionService(ctx *gin.Context, in dto.UpdateTransactionRequest) (dto.TransactionRes, error) {
	logrus.Println("[TransactionService UpdateTransactionService] start.")
	var out dto.TransactionRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	var arg = db.UpdateTransactionParams{
		ID:        in.TxID,
		Status:    in.Status,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg = o.setUpdateParams(arg, in)

	trx, err := o.store.UpdateTransaction(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return out, err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.TransactionRes(trx)

	return out, nil
}

func (o *TransactionService) SoftDeleteTransactionService(ctx *gin.Context, in dto.UpdateInactiveTransactionRequest) error {
	logrus.Println("[TransactionService SoftDeleteProductService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveTransactionParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = o.store.UpdateInactiveTransaction(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return err
	}

	return nil
}

func (o *TransactionService) setCreateParams(arg db.CreateTransactionParams, in dto.CreateTransactionReq) db.CreateTransactionParams {

	if in.CustName != "" {
		arg.CustName = sql.NullString{
			String: in.CustName,
			Valid:  true,
		}
	}
	if in.Amount != "" {
		arg.Amount = sql.NullString{
			String: in.Amount,
			Valid:  true,
		}
	}
	if in.Admin != "" {
		arg.Admin = sql.NullString{
			String: in.Admin,
			Valid:  true,
		}
	}
	if in.TotAmount != "" {
		arg.TotAmount = sql.NullString{
			String: in.TotAmount,
			Valid:  true,
		}
	}
	if in.FeePartner != "" {
		arg.FeePartner = sql.NullString{
			String: in.FeePartner,
			Valid:  true,
		}
	}
	if in.FeePpob != "" {
		arg.FeePpob = sql.NullString{
			String: in.FeePpob,
			Valid:  true,
		}
	}
	if in.CatID > 0 {
		arg.CatID = sql.NullInt64{
			Int64: in.CatID,
			Valid: true,
		}
	}
	if in.CatName != "" {
		arg.CatName = sql.NullString{
			String: in.CatName,
			Valid:  true,
		}
	}
	if in.ProdID > 0 {
		arg.ProdID = sql.NullInt64{
			Int64: in.ProdID,
			Valid: true,
		}
	}
	if in.ProdName != "" {
		arg.ProdName = sql.NullString{
			String: in.ProdName,
			Valid:  true,
		}
	}
	if in.PartnerID > 0 {
		arg.PartnerID = sql.NullInt64{
			Int64: in.PartnerID,
			Valid: true,
		}
	}
	if in.PartnerName != "" {
		arg.PartnerName = sql.NullString{
			String: in.PartnerName,
			Valid:  true,
		}
	}
	if in.ProviderID > 0 {
		arg.ProviderID = sql.NullInt64{
			Int64: in.ProviderID,
			Valid: true,
		}
	}
	if in.ProviderName != "" {
		arg.ProviderName = sql.NullString{
			String: in.ProviderName,
			Valid:  true,
		}
	}
	if in.ReqInqParams != "" {
		arg.ReqInqParams = sql.NullString{
			String: in.ReqInqParams,
			Valid:  true,
		}
	}

	return arg
}

func (o *TransactionService) setUpdateParams(arg db.UpdateTransactionParams, in dto.UpdateTransactionRequest) db.UpdateTransactionParams {

	if in.ReqPayParams != "" {
		arg.ReqPayParams = sql.NullString{
			String: in.ReqPayParams,
			Valid:  true,
		}
	}
	if in.ResPayParams != "" {
		arg.ResPayParams = sql.NullString{
			String: in.ResPayParams,
			Valid:  true,
		}
	}
	if in.ReqAdvParams != "" {
		arg.ReqAdvParams = sql.NullString{
			String: in.ReqAdvParams,
			Valid:  true,
		}
	}
	if in.ReqRevParams != "" {
		arg.ReqRevParams = sql.NullString{
			String: in.ReqRevParams,
			Valid:  true,
		}
	}
	if in.ReqCmtParams != "" {
		arg.ReqCmtParams = sql.NullString{
			String: in.ReqCmtParams,
			Valid:  true,
		}
	}

	return arg
}

func (o *TransactionService) TransactionRes(trx db.Transaction) dto.TransactionRes {
	return dto.TransactionRes{
		ID:           trx.ID,
		TxID:         trx.TxID,
		BillID:       trx.BillID,
		CustName:     trx.CustName.String,
		Amount:       trx.Amount.String,
		Admin:        trx.Admin.String,
		TotAmount:    trx.TotAmount.String,
		FeePartner:   trx.FeePartner.String,
		FeePpob:      trx.FeePpob.String,
		FirstBalance: trx.FirstBalance.String,
		LastBalance:  trx.LastBalance.String,
		Sn:           trx.Sn.String,
		CatID:        trx.CatID.Int64,
		CatName:      trx.CatName.String,
		ProdID:       trx.ProdID.Int64,
		ProdName:     trx.ProdName.String,
		PartnerID:    trx.PartnerID.Int64,
		PartnerName:  trx.PartnerName.String,
		ProviderID:   trx.ProviderID.Int64,
		ProviderName: trx.ProviderName.String,
		Status:       trx.Status,
		CreatedBy:    strconv.Itoa(int(trx.CreatedBy.Int64)),
		PaymentType:  trx.PaymentType.String,
		ReqInqParams: trx.ReqInqParams.String,
		ResInqParams: trx.ResInqParams.String,
		ReqPayParams: trx.ReqPayParams.String,
		ResPayParams: trx.ResPayParams.String,
		ReqCmtParams: trx.ReqCmtParams.String,
		ResCmtParams: trx.ResCmtParams.String,
		ReqAdvParams: trx.ReqAdvParams.String,
		ResAdvParams: trx.ResAdvParams.String,
		ReqRevParams: trx.ReqRevParams.String,
		ResRevParams: trx.ResRevParams.String,
	}
}

func (o *TransactionService) ExportTransaction(ctx *gin.Context, in dto.ListTransactionRequest) (*excelize.File, error) {
	logrus.Println("[TransactionService ExportTransaction] start.")

	//authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	//_, err := userService.validator(ctx, authPayload)
	//if err != nil {
	//	return errors.New("error in user validator")
	//}

	trxList, err := o.getTransactionList(ctx, in)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return nil, err
	}

	layoutFormat := "2006-01-02 15:04:05"
	sheetName := "Sheet1"

	f := excelize.NewFile()
	f.SetCellValue(sheetName, "A1", "tx_id")
	f.SetCellValue(sheetName, "B1", "locket")
	f.SetCellValue(sheetName, "C1", "bill_id")
	f.SetCellValue(sheetName, "D1", "category")
	f.SetCellValue(sheetName, "E1", "product")
	f.SetCellValue(sheetName, "F1", "partner")
	f.SetCellValue(sheetName, "G1", "status")
	f.SetCellValue(sheetName, "H1", "amount")
	f.SetCellValue(sheetName, "I1", "fee_partner")
	f.SetCellValue(sheetName, "J1", "fee_ppob")
	f.SetCellValue(sheetName, "K1", "first_balance")
	f.SetCellValue(sheetName, "L1", "last_balance")
	f.SetCellValue(sheetName, "M1", "sn")
	f.SetCellValue(sheetName, "N1", "date")

	for i, trx := range trxList {
		columnNumber := i + 2
		f.SetCellValue(sheetName, "A"+strconv.Itoa(columnNumber), trx.TxID)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(columnNumber), trx.CreatedBy)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(columnNumber), trx.BillID)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(columnNumber), trx.CatName)
		f.SetCellValue(sheetName, "E"+strconv.Itoa(columnNumber), trx.ProdName)
		f.SetCellValue(sheetName, "F"+strconv.Itoa(columnNumber), trx.PartnerName)
		f.SetCellValue(sheetName, "G"+strconv.Itoa(columnNumber), trx.Status)
		f.SetCellValue(sheetName, "H"+strconv.Itoa(columnNumber), trx.TotAmount)
		f.SetCellValue(sheetName, "I"+strconv.Itoa(columnNumber), trx.FeePartner)
		f.SetCellValue(sheetName, "J"+strconv.Itoa(columnNumber), trx.FeePpob)
		f.SetCellValue(sheetName, "K"+strconv.Itoa(columnNumber), trx.FirstBalance)
		f.SetCellValue(sheetName, "L"+strconv.Itoa(columnNumber), trx.LastBalance)
		f.SetCellValue(sheetName, "M"+strconv.Itoa(columnNumber), trx.Sn)
		f.SetCellValue(sheetName, "N"+strconv.Itoa(columnNumber), trx.CreatedAt.Time.Format(layoutFormat))
	}

	//now := time.Now()
	//
	//f.SetCellValue("Sheet1", "A4", now.Format(time.ANSIC))

	//if errSave := f.SaveAs("simple.xlsx"); err != nil {
	//	logrus.Println("[TransactionService ExportTransaction] error export excel file : ", errSave)
	//	log.Fatal(errSave)
	//}

	return f, nil
}

func (o *TransactionService) getTransactionList(ctx *gin.Context, in dto.ListTransactionRequest) ([]db.Transaction, error) {
	var out []db.Transaction
	arg := db.ListTransactionParams{
		Limit:    in.PageSize,
		Offset:   (in.PageID - 1) * in.PageSize,
		FromDate: in.FromDate,
		ToDate:   in.ToDate,
		Status:   in.Status,
		PaymentType: sql.NullString{
			String: in.PaymentType,
			Valid:  true,
		},
		IsType: true,
	}

	if in.Status != "" {
		arg.IsStatus = true
	}

	if in.CatID > 0 {
		arg.CatID = sql.NullInt64{
			Int64: in.CatID,
			Valid: true,
		}
		arg.IsCat = true
	}

	if in.PartnerID > 0 {
		arg.PartnerID = sql.NullInt64{
			Int64: in.PartnerID,
			Valid: true,
		}
		arg.IsPartner = true
	}

	if in.CreatedBy > 0 {
		arg.CreatedBy = sql.NullInt64{
			Int64: in.CreatedBy,
			Valid: true,
		}
		arg.IsCreated = true
	}

	if in.PaymentType == "" {
		arg.PaymentType = sql.NullString{
			String: "Payment",
			Valid:  true,
		}
	}

	trxList, err := o.store.ListTransaction(ctx, arg)
	if err != nil {
		return out, err
	}

	return trxList, nil
}
