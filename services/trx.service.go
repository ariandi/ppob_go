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

	result = TransactionRes(transaction)
	return result, nil
}

//func (o *ProductService) GetProductService(req dto.GetProductReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.ProductRes, error) {
//	logrus.Println("[ProductService GetProductService] start.")
//	var result dto.ProductRes
//	_, err := validator(store, ctx, authPayload)
//	if err != nil {
//		return result, errors.New("error in user validator")
//	}
//
//	prod, err := store.GetProduct(ctx, req.ID)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
//			return result, err
//		}
//
//		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
//		return result, err
//	}
//
//	result = ProductRes(prod)
//	return result, nil
//}
//
//func (o *ProductService) ListProductService(req dto.ListProductRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) ([]dto.ProductRes, error) {
//	logrus.Println("[ProductService ListProductService] start.")
//	var result []dto.ProductRes
//	_, err := validator(store, ctx, authPayload)
//	if err != nil {
//		return result, errors.New("error in user validator")
//	}
//
//	arg := db.ListProductParams{
//		Limit:  req.PageSize,
//		Offset: (req.PageID - 1) * req.PageSize,
//	}
//
//	products, err := store.ListProduct(ctx, arg)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
//		return result, err
//	}
//
//	for _, prod := range products {
//		u := ProductRes(prod)
//		result = append(result, u)
//	}
//
//	return result, nil
//}
//
//func (o *ProductService) UpdateProductService(req dto.UpdateProductRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.ProductRes, error) {
//	logrus.Println("[ProductService UpdateProductService] start.")
//	var result dto.ProductRes
//	userValid, err := validator(store, ctx, authPayload)
//	if err != nil {
//		return result, errors.New("error in user validator")
//	}
//
//	var arg = db.UpdateProductParams{
//		ID:         req.ID,
//		Name:       req.Name,
//		CatID:      req.CatID,
//		Amount:     req.Amount,
//		ProviderID: req.ProviderID,
//		Status:     req.Status,
//		Parent:     req.Parent,
//		UpdatedBy:  sql.NullInt64{Int64: userValid.ID, Valid: true},
//	}
//
//	arg = o.setUpdateProd(arg, req)
//
//	prod, err := store.UpdateProduct(ctx, arg)
//	if err != nil {
//		if pqErr, ok := err.(*pq.Error); ok {
//			switch pqErr.Code.Name() {
//			case "foreign_key_violation", "unique_violation":
//				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
//				return result, err
//			}
//		}
//		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
//		return result, err
//	}
//
//	result = ProductRes(prod)
//
//	return result, nil
//}
//
//func (o *ProductService) SoftDeleteProductService(req dto.UpdateInactiveProductRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) error {
//	logrus.Println("[ProductService SoftDeleteProviderService] start.")
//	userValid, err := validator(store, ctx, authPayload)
//	if err != nil {
//		return errors.New("error in user validator")
//	}
//
//	arg := db.UpdateInactiveProductParams{
//		ID:        req.ID,
//		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
//	}
//
//	_, err = store.UpdateInactiveProduct(ctx, arg)
//	if err != nil {
//		if pqErr, ok := err.(*pq.Error); ok {
//			switch pqErr.Code.Name() {
//			case "foreign_key_violation", "unique_violation":
//				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
//				return err
//			}
//		}
//		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
//		return err
//	}
//
//	return nil
//}

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

func TransactionRes(trx db.Transaction) dto.TransactionRes {
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
