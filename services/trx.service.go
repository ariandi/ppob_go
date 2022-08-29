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
	store db.Store
}

var transactionService *TransactionService

// GetTransactionService is
func GetTransactionService(store db.Store) *TransactionService {
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

func (o *TransactionService) ListTransactionService(ctx *gin.Context, in dto.ListTransactionRequest) ([]dto.TransactionRes, error) {
	logrus.Println("[TransactionService ListTransactionService] start.")
	var out []dto.TransactionRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.ListTransactionParams{
		Limit:  in.PageSize,
		Offset: (in.PageID - 1) * in.PageSize,
	}

	trxList, err := o.store.ListTransaction(ctx, arg)
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

func (o *TransactionService) InqService(ctx *gin.Context, in dto.InqRequest) (dto.InqResponse, error) {
	logrus.Println("[TransactionService InqService] start.")
	var ret dto.InqResponse

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return ret, errors.New("error in user validator")
	}

	logrus.Println("[TransactionService InqService] begin validate trx.")
	respValidErr, err := o.validateTrx(ctx, in)
	if err != nil {
		return respValidErr, err
	}

	txID := o.setTxID()
	prodID, err := strconv.Atoi(in.ProductCode)
	prod, _ := o.store.GetProduct(ctx, int64(prodID))
	category, _ := o.store.GetCategory(ctx, prod.CatID)
	partnerArg := db.GetPartnerByParamsParams{
		IsUser:     true,
		UserParams: in.AppName,
	}
	partner, _ := o.store.GetPartnerByParams(ctx, partnerArg)
	sellingArg := db.ListSellingByParamsParams{
		Limit:     1,
		Offset:    0,
		IsPartner: true,
		PartnerID: sql.NullInt64{
			Int64: partner.ID,
			Valid: true,
		},
		IsCategory: true,
		CategoryID: sql.NullInt64{
			Int64: prod.CatID,
			Valid: true,
		},
	}
	var selling db.Selling
	sellings, _ := o.store.ListSellingByParams(ctx, sellingArg)
	for _, sell := range sellings {
		selling = sell
	}

	prodAmount, _ := strconv.ParseFloat(prod.Amount, 64)
	amountF, _ := strconv.ParseFloat(selling.Amount.String, 64)
	upSellingF, _ := strconv.ParseFloat(category.UpSelling.String, 64)
	amount := int(amountF)
	upSelling := int(upSellingF)
	totAmount := int(prodAmount) + amount + upSelling
	inInqSetResponse := dto.InqSetResponse{
		InqData:     in,
		ProductName: prod.Name,
		Amount:      int64(prodAmount + upSellingF),
		Admin:       int64(amountF),
		TotalAmount: int64(totAmount),
		ResultCd:    util.SuccessCd,
		ResultMsg:   util.SuccessMsg,
		TxID:        txID,
	}
	ret = o.InqResult(inInqSetResponse)

	reqInqConsume := dto.InqRequestConsume{
		InqRequest:  in,
		InqResponse: ret,
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

	return ret, nil
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

func (o *TransactionService) setTxID() string {

	unique := util.RandomInt(1, 9999)
	uniqueString := strconv.Itoa(int(unique))
	t := time.Now()
	txID := "MADU_" + t.Format("200601021504") + uniqueString

	return txID
}

func (o *TransactionService) validateTrx(ctx *gin.Context, in dto.InqRequest) (dto.InqResponse, error) {

	inqRes := dto.InqResponse{}

	partnerArg := db.GetPartnerByParamsParams{
		IsUser:     true,
		UserParams: in.AppName,
	}
	partner, err := o.store.GetPartnerByParams(ctx, partnerArg)
	if err != nil {
		logrus.Info("select partner error : ", err)
		inqRes = o.InqResultSet(in, util.AppNameNotFoundCd, util.AppNameNotFoundMsg)
		return inqRes, errors.New("app name not exist")
	}

	_, err = o.store.GetUserByUsername(ctx, in.UserID)
	if err != nil {
		logrus.Info("select user error : ", err)
		inqRes = o.InqResultSet(in, util.UserNotFoundCd, util.UserNotFoundMsg)
		return inqRes, err
	}

	prodCode, _ := strconv.ParseInt(in.ProductCode, 10, 64)
	prod, err := o.store.GetProduct(ctx, prodCode)
	if err != nil {
		logrus.Info("select prod error", err)
		inqRes = o.InqResultSet(in, util.ProductNotFoundCd, util.ProductNotFoundMsg)
		return inqRes, err
	}

	if prod.ProviderCode == "" || prod.ProviderCode == "-" {
		logrus.Info("select prod code is empty", err)
		err = errors.New(util.ProductNotFoundMsg)
		inqRes = o.InqResultSet(in, util.ProductNotFoundCd, util.ProductNotFoundMsg)
		return inqRes, err
	}

	_, err = o.store.GetCategory(ctx, prod.CatID)
	if err != nil {
		logrus.Info("select GetCategory error : ", err)
		inqRes = o.InqResultSet(in, util.CategoryNotFoundCd, util.CategoryNotFoundMsg)
		return inqRes, err
	}

	_, err = o.store.GetProvider(ctx, prod.ProviderID)
	if err != nil {
		logrus.Info("select GetProvider error : ", err)
		inqRes = o.InqResultSet(in, util.ProviderNotFoundCd, util.ProviderNotFoundMsg)
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
	sellingPrice, err := o.store.ListSellingByParams(ctx, sellingArg)
	if err != nil || len(sellingPrice) == 0 {
		logrus.Info("select ListSellingByParams error : ", err)
		if err == nil {
			err = errors.New(util.SellingPriceNotFoundMsg)
		}
		inqRes = o.InqResultSet(in, util.SellingPriceNotFoundCd, util.SellingPriceNotFoundMsg)
		return inqRes, err
	}

	if len(in.TimeStamp) != 14 {
		logrus.Info("select time stamp error", err)
		err = errors.New(util.TimeStampLengthInvalidMsg)
		inqRes = o.InqResultSet(in, util.TimeStampLengthInvalidCd, util.TimeStampLengthInvalidMsg)
		return inqRes, err
	}

	layoutFormat := "20060102150405"
	value := in.TimeStamp
	timeStampStr, err := time.Parse(layoutFormat, value)
	if err != nil {
		logrus.Info("[TrxService setTxID] err timestamp format : ", err)
		err = errors.New(util.TimeStampFormatInvalidMsg)
		inqRes = o.InqResultSet(in, util.TimeStampFormatInvalidCd, util.TimeStampFormatInvalidMsg)
		return inqRes, err
	}

	sha256Req := in.BillID + in.ProductCode + in.UserID + in.RefID + partner.Secret + in.TimeStamp
	hash := sha256.Sum256([]byte(sha256Req))
	sha256Res := hash[:]
	logrus.Info("[TrxService setTxID] local token params is : ", sha256Req)
	logrus.Info("[TrxService setTxID] local token is : ", hex.EncodeToString(sha256Res))
	logrus.Info("[TrxService setTxID] merchant token is : ", in.MerchantToken)

	if hex.EncodeToString(sha256Res) != in.MerchantToken {
		logrus.Info("token not same", err)
		err = errors.New(util.MerchantTokenErrorMsg)
		inqRes = o.InqResultSet(in, util.MerchantTokenErrorCd, util.MerchantTokenErrorMsg)
		return inqRes, err
	}

	logrus.Info("[TrxService setTxID] ref id date is : ", timeStampStr.Format("2006-01-02 15:04:05"))
	trxRefArg := db.GetTransactionByRefIDParams{
		RefID: in.RefID,
		PartnerID: sql.NullInt64{
			Int64: partner.ID,
			Valid: true,
		},
	}
	refID, err := o.store.GetTransactionByRefID(ctx, trxRefArg)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("tx id not found", err)
		} else {
			logrus.Info("select trx table error", err)
			inqRes = o.InqResultSet(in, util.GeneralErrorCd, util.GeneralErrorMsg)
			return inqRes, err
		}
	}

	if refID.ID > 0 {
		err = errors.New(util.RefIDAlreadyUsedMsg)
		inqRes = o.InqResultSet(in, util.RefIDAlreadyUsedCd, util.RefIDAlreadyUsedMsg)
		return inqRes, err
	}

	logrus.Info("[TrxService setTxID] begin check pending trx.")
	_, err = o.store.GetTransactionPending(ctx, in.BillID)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("tx id not found", err)
		} else {
			logrus.Info("select trx table error", err)
			err = errors.New(util.GeneralErrorMsg)
			inqRes = o.InqResultSet(in, util.GeneralErrorCd, util.GeneralErrorMsg)
			return inqRes, err
		}
	}

	if refID.BillID != "" {
		err = errors.New(util.StillPendingTransactionMsg)
		inqRes = o.InqResultSet(in, util.StillPendingTransactionCd, util.StillPendingTransactionMsg)
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
		CreatedBy:    strconv.Itoa(int(trx.CreatedBy.Int64)),
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
		ProductName:   in.ProductName,
		Amount:        in.Amount,
		Admin:         in.Admin,
		TotalAmount:   in.TotalAmount,
		ResultCd:      in.ResultCd,
		ResultMsg:     in.ResultMsg,
		TxID:          in.TxID,
	}
}

func (o *TransactionService) InqResultSet(in dto.InqRequest, resultCd string, resultMsg string) dto.InqResponse {
	inInqSetResponse := dto.InqSetResponse{
		InqData:   in,
		ResultCd:  resultCd,
		ResultMsg: resultMsg,
	}
	return o.InqResult(inInqSetResponse)
}
