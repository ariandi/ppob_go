package services

import (
	"database/sql"
	"errors"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

// TransactionService is
type TransactionService struct {
}

var transactionService *TransactionService

// GetTransactionService is
func GetTransactionService() *TransactionService {
	if transactionService == nil {
		transactionService = new(TransactionService)
	}
	return transactionService
}

func (o *TransactionService) CreateTransactionService(req dto.CreateTransactionReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.TransactionRes, error) {
	logrus.Println("[ProductService CreateTransactionService] start.")
	var result dto.TransactionRes

	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.CreateTransactionParams{
		TxID:      req.TxID,
		BillID:    req.BillID,
		Status:    req.Status,
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg = o.setCreateParams(arg, req)

	transaction, err := store.CreateTransaction(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return result, err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	result = o.TransactionRes(transaction)
	return result, nil
}

func (o *TransactionService) GetTransactionService(req dto.GetTransactionByTxIDReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.TransactionRes, error) {
	logrus.Println("[TransactionService GetTransactionService] start.")
	var result dto.TransactionRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	trx, err := store.GetTransactionByTxID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return result, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	result = o.TransactionRes(trx)
	return result, nil
}

func (o *TransactionService) ListTransactionService(req dto.ListTransactionRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) ([]dto.TransactionRes, error) {
	logrus.Println("[TransactionService ListTransactionService] start.")
	var result []dto.TransactionRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.ListTransactionParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	trxs, err := store.ListTransaction(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	for _, trx := range trxs {
		u := o.TransactionRes(trx)
		result = append(result, u)
	}

	return result, nil
}

func (o *TransactionService) UpdateTransactionService(req dto.UpdateTransactionRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.TransactionRes, error) {
	logrus.Println("[TransactionService UpdateTransactionService] start.")
	var result dto.TransactionRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	var arg = db.UpdateTransactionParams{
		ID:        req.TxID,
		Status:    req.Status,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg = o.setUpdateParams(arg, req)

	trx, err := store.UpdateTransaction(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return result, err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	result = o.TransactionRes(trx)

	return result, nil
}

func (o *TransactionService) SoftDeleteTransactionService(req dto.UpdateInactiveTransactionRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) error {
	logrus.Println("[TransactionService SoftDeleteProductService] start.")
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveTransactionParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = store.UpdateInactiveTransaction(ctx, arg)
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

func (o TransactionService) setCreateParams(arg db.CreateTransactionParams, req dto.CreateTransactionReq) db.CreateTransactionParams {

	if req.CustName != "" {
		arg.CustName = sql.NullString{
			String: req.CustName,
			Valid:  true,
		}
	}
	if req.Amount != "" {
		arg.Amount = sql.NullString{
			String: req.Amount,
			Valid:  true,
		}
	}
	if req.Admin != "" {
		arg.Admin = sql.NullString{
			String: req.Admin,
			Valid:  true,
		}
	}
	if req.TotAmount != "" {
		arg.TotAmount = sql.NullString{
			String: req.TotAmount,
			Valid:  true,
		}
	}
	if req.FeePartner != "" {
		arg.FeePartner = sql.NullString{
			String: req.FeePartner,
			Valid:  true,
		}
	}
	if req.FeePpob != "" {
		arg.FeePpob = sql.NullString{
			String: req.FeePpob,
			Valid:  true,
		}
	}
	if req.CatID > 0 {
		arg.CatID = sql.NullInt64{
			Int64: req.CatID,
			Valid: true,
		}
	}
	if req.CatName != "" {
		arg.CatName = sql.NullString{
			String: req.CatName,
			Valid:  true,
		}
	}
	if req.ProdID > 0 {
		arg.ProdID = sql.NullInt64{
			Int64: req.ProdID,
			Valid: true,
		}
	}
	if req.ProdName != "" {
		arg.ProdName = sql.NullString{
			String: req.ProdName,
			Valid:  true,
		}
	}
	if req.PartnerID > 0 {
		arg.PartnerID = sql.NullInt64{
			Int64: req.PartnerID,
			Valid: true,
		}
	}
	if req.PartnerName != "" {
		arg.PartnerName = sql.NullString{
			String: req.PartnerName,
			Valid:  true,
		}
	}
	if req.ProviderID > 0 {
		arg.ProviderID = sql.NullInt64{
			Int64: req.ProviderID,
			Valid: true,
		}
	}
	if req.ProviderName != "" {
		arg.ProviderName = sql.NullString{
			String: req.ProviderName,
			Valid:  true,
		}
	}
	if req.ReqInqParams != "" {
		arg.ReqInqParams = sql.NullString{
			String: req.ReqInqParams,
			Valid:  true,
		}
	}

	return arg
}

func (o TransactionService) setUpdateParams(arg db.UpdateTransactionParams, req dto.UpdateTransactionRequest) db.UpdateTransactionParams {

	if req.ReqPayParams != "" {
		arg.ReqPayParams = sql.NullString{
			String: req.ReqPayParams,
			Valid:  true,
		}
	}
	if req.ResPayParams != "" {
		arg.ResPayParams = sql.NullString{
			String: req.ResPayParams,
			Valid:  true,
		}
	}
	if req.ReqAdvParams != "" {
		arg.ReqAdvParams = sql.NullString{
			String: req.ReqAdvParams,
			Valid:  true,
		}
	}
	if req.ReqRevParams != "" {
		arg.ReqRevParams = sql.NullString{
			String: req.ReqRevParams,
			Valid:  true,
		}
	}
	if req.ReqCmtParams != "" {
		arg.ReqCmtParams = sql.NullString{
			String: req.ReqCmtParams,
			Valid:  true,
		}
	}

	return arg
}

func (o TransactionService) TransactionRes(trx db.Transaction) dto.TransactionRes {
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
		CatID:        trx.CatID.Int64,
		CatName:      trx.CatName.String,
		ProdID:       trx.ProdID.Int64,
		ProdName:     trx.ProdName.String,
		PartnerID:    trx.PartnerID.Int64,
		PartnerName:  trx.PartnerName.String,
		ProviderID:   trx.ProviderID.Int64,
		ProviderName: trx.ProviderName.String,
		Status:       trx.Status,
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
