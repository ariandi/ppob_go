package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/ariandi/ppob_go/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
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

func (o *TransactionService) InqService(req dto.InqRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.InqResponse, error) {
	logrus.Println("[TransactionService InqService] start.")
	var ret dto.InqResponse
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return ret, errors.New("error in user validator")
	}

	logrus.Println("[TransactionService InqService] begin validate trx.")
	respValidErr, err := o.validateTrx(req, ctx, store)
	if err != nil {
		return respValidErr, err
	}

	txID := o.setTxID()
	reqInqConsume := dto.InqRequestConsume{
		InqRequest: req,
		TxID:       txID,
	}

	queueName := "transactions"
	redisQueue, err := redisConn.OpenQueue(queueName)
	if err != nil {
		return ret, err
	}

	byt, err := json.Marshal(reqInqConsume)
	if err != nil {
		return ret, err
	}

	err = redisQueue.Publish(string(byt))
	in := dto.InqSetResponse{
		InqData:   req,
		TxID:      txID,
		ResultCd:  util.SuccessCd,
		ResultMsg: util.SuccessMsg,
	}
	ret = o.InqResult(in)

	return ret, nil
}

func (o *TransactionService) setCreateParams(arg db.CreateTransactionParams, req dto.CreateTransactionReq) db.CreateTransactionParams {

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

func (o *TransactionService) setUpdateParams(arg db.UpdateTransactionParams, req dto.UpdateTransactionRequest) db.UpdateTransactionParams {

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

func (o *TransactionService) setTxID() string {

	unique := util.RandomInt(1, 9999)
	uniqueString := strconv.Itoa(int(unique))
	t := time.Now()
	txID := "MADU_" + t.Format("200601021504") + uniqueString

	return txID
}

func (o *TransactionService) validateTrx(req dto.InqRequest, ctx *gin.Context, store db.Store) (dto.InqResponse, error) {

	inqRes := dto.InqResponse{}

	partnerArg := db.GetPartnerByParamsParams{
		IsUser:     true,
		UserParams: req.AppName,
	}
	partner, err := store.GetPartnerByParams(ctx, partnerArg)
	if err != nil {
		logrus.Info("select partner error : ", err)
		inqRes = o.InqResultSet(req, util.AppNameNotFoundCd, util.AppNameNotFoundMsg)
		return inqRes, errors.New("app name not exist")
	}

	_, err = store.GetUserByUsername(ctx, req.UserID)
	if err != nil {
		logrus.Info("select user error : ", err)
		inqRes = o.InqResultSet(req, util.UserNotFoundCd, util.UserNotFoundMsg)
		return inqRes, err
	}

	prodCode, _ := strconv.ParseInt(req.ProductCode, 10, 64)
	prod, err := store.GetProduct(ctx, prodCode)
	if err != nil {
		logrus.Info("select prod error", err)
		inqRes = o.InqResultSet(req, util.ProductNotFoundCd, util.ProductNotFoundMsg)
		return inqRes, err
	}

	if prod.ProviderCode == "" || prod.ProviderCode == "-" {
		logrus.Info("select prod code is empty", err)
		err = errors.New(util.ProductNotFoundMsg)
		inqRes = o.InqResultSet(req, util.ProductNotFoundCd, util.ProductNotFoundMsg)
		return inqRes, err
	}

	_, err = store.GetCategory(ctx, prod.CatID)
	if err != nil {
		logrus.Info("select GetCategory error : ", err)
		inqRes = o.InqResultSet(req, util.CategoryNotFoundCd, util.CategoryNotFoundMsg)
		return inqRes, err
	}

	_, err = store.GetProvider(ctx, prod.ProviderID)
	if err != nil {
		logrus.Info("select GetProvider error : ", err)
		inqRes = o.InqResultSet(req, util.ProviderNotFoundCd, util.ProviderNotFoundMsg)
		return inqRes, err
	}

	sellingArg := db.ListSellingByParamsParams{
		Limit:     10,
		Offset:    0,
		IsPartner: true,
		PartnerID: sql.NullInt64{
			Int64: prod.ProviderID,
			Valid: true,
		},
		IsCategory: true,
		CategoryID: sql.NullInt64{
			Int64: prod.CatID,
			Valid: true,
		},
	}
	sellingPrice, err := store.ListSellingByParams(ctx, sellingArg)
	if err != nil || len(sellingPrice) == 0 {
		logrus.Info("select ListSellingByParams error : ", err)
		if err == nil {
			err = errors.New(util.SellingPriceNotFoundMsg)
		}
		inqRes = o.InqResultSet(req, util.SellingPriceNotFoundCd, util.SellingPriceNotFoundMsg)
		return inqRes, err
	}

	if len(req.TimeStamp) != 14 {
		logrus.Info("select time stamp error", err)
		err = errors.New(util.TimeStampLengthInvalidMsg)
		inqRes = o.InqResultSet(req, util.TimeStampLengthInvalidCd, util.TimeStampLengthInvalidMsg)
		return inqRes, err
	}

	layoutFormat := "20060102150405"
	value := req.TimeStamp
	timeStampStr, err := time.Parse(layoutFormat, value)
	if err != nil {
		logrus.Info("[TrxService setTxID] err timestamp format : ", err)
		err = errors.New(util.TimeStampFormatInvalidMsg)
		inqRes = o.InqResultSet(req, util.TimeStampFormatInvalidCd, util.TimeStampFormatInvalidMsg)
		return inqRes, err
	}

	sha256Req := req.BillID + req.ProductCode + req.UserID + req.RefID + partner.Secret + req.TimeStamp
	hash := sha256.Sum256([]byte(sha256Req))
	sha256Res := hash[:]
	logrus.Info("[TrxService setTxID] local token params is : ", sha256Req)
	logrus.Info("[TrxService setTxID] local token is : ", hex.EncodeToString(sha256Res))
	logrus.Info("[TrxService setTxID] merchant token is : ", req.MerchantToken)

	if hex.EncodeToString(sha256Res) != req.MerchantToken {
		logrus.Info("token not same", err)
		err = errors.New(util.MerchantTokenErrorMsg)
		inqRes = o.InqResultSet(req, util.MerchantTokenErrorCd, util.MerchantTokenErrorMsg)
		return inqRes, err
	}

	logrus.Info("[TrxService setTxID] ref id date is : ", timeStampStr.Format("2006-01-02 15:04:05"))
	trxRefArg := db.GetTransactionByRefIDParams{
		RefID: req.RefID,
		PartnerID: sql.NullInt64{
			Int64: partner.ID,
			Valid: true,
		},
	}
	refID, err := store.GetTransactionByRefID(ctx, trxRefArg)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("tx id not found", err)
		} else {
			logrus.Info("select trx table error", err)
			inqRes = o.InqResultSet(req, util.GeneralErrorCd, util.GeneralErrorMsg)
			return inqRes, err
		}
	}

	if refID.ID > 0 {
		err = errors.New(util.RefIDAlreadyUsedMsg)
		inqRes = o.InqResultSet(req, util.RefIDAlreadyUsedCd, util.RefIDAlreadyUsedMsg)
		return inqRes, err
	}

	logrus.Info("[TrxService setTxID] begin check pending trx.")
	_, err = store.GetTransactionPending(ctx, req.BillID)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("tx id not found", err)
		} else {
			logrus.Info("select trx table error", err)
			err = errors.New(util.GeneralErrorMsg)
			inqRes = o.InqResultSet(req, util.GeneralErrorCd, util.GeneralErrorMsg)
			return inqRes, err
		}
	}

	if refID.BillID != "" {
		err = errors.New(util.StillPendingTransactionMsg)
		inqRes = o.InqResultSet(req, util.StillPendingTransactionCd, util.StillPendingTransactionMsg)
		return inqRes, err
	}

	return inqRes, nil
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

func (o *TransactionService) InqResult(in dto.InqSetResponse) dto.InqResponse {
	return dto.InqResponse{
		TimeStamp:     in.InqData.TimeStamp,
		UserID:        in.InqData.UserID,
		RefID:         in.InqData.RefID,
		BillID:        in.InqData.BillID,
		AppName:       in.InqData.AppName,
		ProductCode:   in.InqData.ProductCode,
		MerchantToken: in.InqData.MerchantToken,
		Amount:        in.InqData.Amount,
		TxID:          in.TxID,
		ResultCd:      in.ResultCd,
		ResultMsg:     in.ResultMsg,
	}
}

func (o *TransactionService) InqResultSet(req dto.InqRequest, resultCd string, resultMsg string) dto.InqResponse {
	in := dto.InqSetResponse{
		InqData:   req,
		ResultCd:  resultCd,
		ResultMsg: resultMsg,
	}
	return o.InqResult(in)
}
