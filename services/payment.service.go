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
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type PaymentInterface interface {
	InqService(ctx *gin.Context, in dto.InqRequest) (dto.InqResponse, error)
	PayService(ctx *gin.Context, in dto.PayRequest) (dto.PayResponse, error)
	DepositService(ctx *gin.Context, in dto.DepositRequest) (dto.DepositResponse, error)
	DepositApproveService(ctx *gin.Context, in dto.DepositApproveRequest) (dto.DepositResponse, error)
	setTxID() string
	validateTrx(ctx *gin.Context, in dto.InqRequest) (dto.InqResponse, error)
	InqResult(in dto.InqSetResponse) dto.InqResponse
	InqResultSet(in dto.InqRequest, resultCd string, resultMsg string) dto.InqResponse
}

type PaymentService struct {
	store db.Store
}

var paymentService *PaymentService

func GetPaymentService(store db.Store) PaymentInterface {
	if paymentService == nil {
		paymentService = &PaymentService{
			store: store,
		}
	}
	return paymentService
}

func (o *PaymentService) InqService(ctx *gin.Context, in dto.InqRequest) (dto.InqResponse, error) {
	logrus.Println("[PaymentService InqService] start.")
	var ret dto.InqResponse

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return ret, errors.New("error in user validator")
	}

	logrus.Println("[PaymentService InqService] begin validate trx.")
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
		QueueName:   util.TrxInq,
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

func (o *PaymentService) PayService(ctx *gin.Context, in dto.PayRequest) (dto.PayResponse, error) {
	logrus.Println("[PaymentService PayService] start.")
	var ret dto.PayResponse

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return ret, errors.New("error in user validator")
	}

	logrus.Println("[PaymentService PayService] begin validate trx.")
	trx, err := o.store.GetTransactionByTxID(ctx, in.TxID)
	if err != nil {
		logrus.Info("[PaymentService PayService] ", util.TxIdNotFoundMsg)
		payRes := o.PayResultErrorSet(in, util.TxIdNotFoundCd, util.TxIdNotFoundMsg)
		return payRes, errors.New(util.TxIdNotFoundMsg)
	}

	logrus.Println("[PaymentService PayService] tr id is : ", trx.ID)

	respValidErr, err := o.validateTrxPay(ctx, in, trx)
	if err != nil {
		return respValidErr, err
	}

	inPayResponse := dto.PayResponse{
		TimeStamp:     in.TimeStamp,
		UserID:        in.UserID,
		RefID:         in.RefID,
		BillID:        in.BillID,
		AppName:       in.AppName,
		ProductCode:   in.ProductCode,
		MerchantToken: in.MerchantToken,
		ProductName:   trx.ProdName.String,
		Amount:        in.Amount,
		Admin:         in.Admin,
		TotalAmount:   in.TotalAmount,
		ResultCd:      util.PendingCd,
		ResultMsg:     util.PendingMsg,
		TxID:          in.TxID,
	}

	reqInqConsume := dto.InqRequestConsume{
		PayRequest:  in,
		PayResponse: inPayResponse,
		QueueName:   util.TrxPayment,
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

	return inPayResponse, nil
}

func (o *PaymentService) DepositService(ctx *gin.Context, in dto.DepositRequest) (dto.DepositResponse, error) {
	logrus.Println("[PaymentService DepositService] start.")
	var ret dto.DepositResponse

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userReq, err := userService.validator(ctx, authPayload)
	if err != nil {
		return ret, errors.New("error in user validator")
	}

	user := dto.UserResponse{
		ID:       userReq.ID,
		Name:     userReq.Name,
		Email:    userReq.Email,
		Username: userReq.Username,
		Balance:  userReq.Balance,
	}

	txID := o.setTxID()
	queueName := "deposit"
	redisQueue, err := redisConn.OpenQueue(queueName)
	if err != nil {
		return ret, err
	}

	ret = dto.DepositResponse{
		ResultCd:  util.SuccessCd,
		ResultMsg: util.SuccessMsg,
		TxID:      txID,
	}
	depositInqConsume := dto.DepositRequestConsume{
		DepositRequest:  in,
		DepositResponse: ret,
		UserRequest:     user,
		QueueName:       util.DEPOSIT_TYPE_REQUEST,
	}
	byt, err := json.Marshal(depositInqConsume)
	if err != nil {
		return ret, err
	}

	err = redisQueue.Publish(string(byt))

	return ret, nil
}

func (o *PaymentService) DepositApproveService(ctx *gin.Context, in dto.DepositApproveRequest) (dto.DepositResponse, error) {
	logrus.Println("[PaymentService DepositApproveService] start.")
	var ret dto.DepositResponse

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userReq, err := userService.validator(ctx, authPayload)
	if err != nil {
		return ret, errors.New("error in user validator")
	}

	if authPayload.Username != "dbduabelas" {
		return ret, errors.New("you not allow to approve deposit")
	}

	_, err = o.store.GetTransactionByTxID(ctx, in.TxID)
	if err != nil {
		logrus.Info("[PaymentService DepositApproveService] select tx id not found : ", err)
		ret = dto.DepositResponse{
			ResultCd:  util.TransactionNotFoundCd,
			ResultMsg: util.TransactionNotFoundMsg,
			TxID:      in.TxID,
		}
		return ret, err
	}

	user := dto.UserResponse{
		ID:       userReq.ID,
		Name:     userReq.Name,
		Email:    userReq.Email,
		Username: userReq.Username,
		Balance:  userReq.Balance,
	}

	queueName := "deposit"
	redisQueue, err := redisConn.OpenQueue(queueName)
	if err != nil {
		return ret, err
	}

	ret = dto.DepositResponse{
		ResultCd:  util.SuccessCd,
		ResultMsg: util.SuccessMsg,
		TxID:      in.TxID,
	}
	depositInqConsume := dto.DepositRequestConsume{
		DepositApproveRequest: in,
		DepositResponse:       ret,
		UserRequest:           user,
		QueueName:             util.DEPOSIT_TYPE_APPROVE,
	}
	byt, err := json.Marshal(depositInqConsume)
	if err != nil {
		return ret, err
	}

	err = redisQueue.Publish(string(byt))

	return ret, nil
}

func (o *PaymentService) setTxID() string {

	unique := util.RandomInt(1, 9999)
	uniqueString := strconv.Itoa(int(unique))
	t := time.Now()
	txID := "MADU_" + t.Format("200601021504") + uniqueString

	return txID
}

func (o *PaymentService) validateTrx(ctx *gin.Context, in dto.InqRequest) (dto.InqResponse, error) {
	logrus.Info("[PaymentService validateTrx] start.")
	inqRes := dto.InqResponse{}

	partnerArg := db.GetPartnerByParamsParams{
		IsUser:     true,
		UserParams: in.AppName,
	}
	partner, err := o.store.GetPartnerByParams(ctx, partnerArg)
	if err != nil {
		logrus.Info("[PaymentService validateTrx] select partner error : ", err)
		inqRes = o.InqResultSet(in, util.AppNameNotFoundCd, util.AppNameNotFoundMsg)
		return inqRes, errors.New("app name not exist")
	}

	_, err = o.store.GetUserByUsername(ctx, in.UserID)
	if err != nil {
		logrus.Info("[PaymentService validateTrx] select user error : ", err)
		inqRes = o.InqResultSet(in, util.UserNotFoundCd, util.UserNotFoundMsg)
		return inqRes, err
	}

	prodCode, _ := strconv.ParseInt(in.ProductCode, 10, 64)
	prod, err := o.store.GetProduct(ctx, prodCode)
	if err != nil {
		logrus.Info("[PaymentService validateTrx] select prod error", err)
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
		logrus.Info("[PaymentService validateTrx] select GetCategory error : ", err)
		inqRes = o.InqResultSet(in, util.CategoryNotFoundCd, util.CategoryNotFoundMsg)
		return inqRes, err
	}

	_, err = o.store.GetProvider(ctx, prod.ProviderID)
	if err != nil {
		logrus.Info("[PaymentService validateTrx] select GetProvider error : ", err)
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
		logrus.Info("[PaymentService validateTrx] select ListSellingByParams error : ", err)
		if err == nil {
			err = errors.New(util.SellingPriceNotFoundMsg)
		}
		inqRes = o.InqResultSet(in, util.SellingPriceNotFoundCd, util.SellingPriceNotFoundMsg)
		return inqRes, err
	}

	if len(in.TimeStamp) != 14 {
		logrus.Info("[PaymentService validateTrx] select time stamp error", err)
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
	logrus.Info("[PaymentService validateTrx] local token params is : ", sha256Req)
	logrus.Info("[PaymentService validateTrx] local token is : ", hex.EncodeToString(sha256Res))
	logrus.Info("[PaymentService validateTrx] merchant token is : ", in.MerchantToken)

	if hex.EncodeToString(sha256Res) != in.MerchantToken {
		logrus.Info("[PaymentService validateTrx] token not same", err)
		err = errors.New(util.MerchantTokenErrorMsg)
		inqRes = o.InqResultSet(in, util.MerchantTokenErrorCd, util.MerchantTokenErrorMsg)
		return inqRes, err
	}

	logrus.Info("[PaymentService validateTrx] ref id date is : ", timeStampStr.Format("2006-01-02 15:04:05"))
	trxRefArg := db.GetTransactionByRefIDParams{
		IsReff:    true,
		RefID:     in.RefID,
		IsPartner: true,
		PartnerID: sql.NullInt64{
			Int64: partner.ID,
			Valid: true,
		},
	}
	refID, err := o.store.GetTransactionByRefID(ctx, trxRefArg)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("[PaymentService validateTrx] tx id not found : ", err)
		} else {
			logrus.Info("[PaymentService validateTrx] select trx table error", err)
			inqRes = o.InqResultSet(in, util.GeneralErrorCd, util.GeneralErrorMsg)
			return inqRes, err
		}
	}

	logrus.Info("[PaymentService validateTrx] refID.Status : ", refID.Status)
	if refID.Status != "" {
		logrus.Info("[PaymentService validateTrx] " + util.RefIDAlreadyUsedMsg)
		err = errors.New(util.RefIDAlreadyUsedMsg)
		inqRes = o.InqResultSet(in, util.RefIDAlreadyUsedCd, util.RefIDAlreadyUsedMsg)
		return inqRes, err
	}

	logrus.Info("[PaymentService validateTrx] begin check pending trx.")
	checkBillID, err := o.store.GetTransactionPending(ctx, in.BillID)
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

	if checkBillID.BillID != "" {
		err = errors.New(util.StillPendingTransactionMsg)
		inqRes = o.InqResultSet(in, util.StillPendingTransactionCd, util.StillPendingTransactionMsg)
		return inqRes, err
	}

	return inqRes, nil
}

func (o *PaymentService) validateTrxPay(ctx *gin.Context, in dto.PayRequest, trx db.Transaction) (dto.PayResponse, error) {
	logrus.Info("[PaymentService validateTrxPay] start.")
	payRes := dto.PayResponse{}

	if in.BillID != trx.BillID {
		logrus.Info("[PaymentService validateTrxPay] ", util.BillIdDifferentMsg)
		err := errors.New(util.BillIdDifferentMsg)
		payRes = o.PayResultErrorSet(in, util.BillIdDifferentCd, util.BillIdDifferentMsg)
		return payRes, err
	}

	if trx.Status != util.INQ {
		logrus.Info("[PaymentService validateTrxPay] ", util.TxNotInqMsg)
		err := errors.New(util.TxNotInqMsg)
		payRes = o.PayResultErrorSet(in, util.TxNotInqCd, util.TxNotInqMsg)
		return payRes, err
	}

	partnerArg := db.GetPartnerByParamsParams{
		IsUser:     true,
		UserParams: in.AppName,
	}
	partner, err := o.store.GetPartnerByParams(ctx, partnerArg)
	if err != nil {
		logrus.Info("[PaymentService validateTrxPay] select partner error : ", err)
		payRes = o.PayResultErrorSet(in, util.AppNameNotFoundCd, util.AppNameNotFoundMsg)
		return payRes, errors.New("app name not exist")
	}

	if trx.PartnerID.Int64 != partner.ID {
		logrus.Info("[PaymentService validateTrxPay] partner not same error : ", err)
		payRes = o.PayResultErrorSet(in, util.AppNameNotFoundCd, util.AppNameNotFoundMsg)
		return payRes, errors.New("app name not same")
	}

	user, err := o.store.GetUserByUsername(ctx, in.UserID)
	if err != nil {
		logrus.Info("[PaymentService validateTrxPay] select user error : ", err)
		payRes = o.PayResultErrorSet(in, util.UserNotFoundCd, util.UserNotFoundMsg)
		return payRes, err
	}

	if trx.CreatedBy.Int64 != user.ID {
		logrus.Info("[PaymentService validateTrxPay] user not same error : ", err)
		payRes = o.PayResultErrorSet(in, util.AppNameNotFoundCd, util.AppNameNotFoundMsg)
		return payRes, errors.New("user id not same")
	}

	prodCode, _ := strconv.ParseInt(in.ProductCode, 10, 64)
	prod, err := o.store.GetProduct(ctx, prodCode)
	if err != nil {
		logrus.Info("[PaymentService validateTrxPay] select prod error", err)
		payRes = o.PayResultErrorSet(in, util.ProductNotFoundCd, util.ProductNotFoundMsg)
		return payRes, err
	}

	if prod.ProviderCode == "" || prod.ProviderCode == "-" {
		logrus.Info("[PaymentService validateTrxPay] select prod code is empty", err)
		err = errors.New(util.ProductNotFoundMsg)
		payRes = o.PayResultErrorSet(in, util.ProductNotFoundCd, util.ProductNotFoundMsg)
		return payRes, err
	}

	if trx.ProdID.Int64 != prod.ID {
		logrus.Info("[PaymentService validateTrxPay] product id not same error : ", err)
		payRes = o.PayResultErrorSet(in, util.AppNameNotFoundCd, util.AppNameNotFoundMsg)
		return payRes, errors.New("product id not same")
	}

	cat, err := o.store.GetCategory(ctx, prod.CatID)
	if err != nil {
		logrus.Info("[PaymentService validateTrxPay] select GetCategory error : ", err)
		payRes = o.PayResultErrorSet(in, util.CategoryNotFoundCd, util.CategoryNotFoundMsg)
		return payRes, err
	}

	if trx.CatID.Int64 != cat.ID {
		logrus.Info("[PaymentService validateTrxPay] category id not same error : ", err)
		payRes = o.PayResultErrorSet(in, util.AppNameNotFoundCd, util.AppNameNotFoundMsg)
		return payRes, errors.New("category id not same")
	}

	provider, err := o.store.GetProvider(ctx, prod.ProviderID)
	if err != nil {
		logrus.Info("[PaymentService validateTrxPay] select GetProvider error : ", err)
		payRes = o.PayResultErrorSet(in, util.ProviderNotFoundCd, util.ProviderNotFoundMsg)
		return payRes, err
	}

	if trx.ProviderID.Int64 != provider.ID {
		logrus.Info("[PaymentService validateTrxPay] provider id not same error : ", err)
		payRes = o.PayResultErrorSet(in, util.AppNameNotFoundCd, util.AppNameNotFoundMsg)
		return payRes, errors.New("provider id not same")
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
		logrus.Info("[PaymentService validateTrxPay] select ListSellingByParams error : ", err)
		if err == nil {
			err = errors.New(util.SellingPriceNotFoundMsg)
		}
		payRes = o.PayResultErrorSet(in, util.SellingPriceNotFoundCd, util.SellingPriceNotFoundMsg)
		return payRes, err
	}

	if len(in.TimeStamp) != 14 {
		logrus.Info("[PaymentService validateTrxPay] select time stamp error", err)
		err = errors.New(util.TimeStampLengthInvalidMsg)
		payRes = o.PayResultErrorSet(in, util.TimeStampLengthInvalidCd, util.TimeStampLengthInvalidMsg)
		return payRes, err
	}

	layoutFormat := "20060102150405"
	value := in.TimeStamp
	timeStampStr, err := time.Parse(layoutFormat, value)
	if err != nil {
		logrus.Info("[PaymentService validateTrxPay] err timestamp format : ", err)
		err = errors.New(util.TimeStampFormatInvalidMsg)
		payRes = o.PayResultErrorSet(in, util.TimeStampFormatInvalidCd, util.TimeStampFormatInvalidMsg)
		return payRes, err
	}

	sha256Req := in.BillID + in.ProductCode + in.UserID + in.RefID + partner.Secret + in.TimeStamp
	hash := sha256.Sum256([]byte(sha256Req))
	sha256Res := hash[:]
	logrus.Info("[PaymentService validateTrxPay] local token params is : ", sha256Req)
	logrus.Info("[PaymentService validateTrxPay] local token is : ", hex.EncodeToString(sha256Res))
	logrus.Info("[PaymentService validateTrxPay] merchant token is : ", in.MerchantToken)

	if hex.EncodeToString(sha256Res) != in.MerchantToken {
		logrus.Info("[PaymentService validateTrxPay] token not same", err)
		err = errors.New(util.MerchantTokenErrorMsg)
		payRes = o.PayResultErrorSet(in, util.MerchantTokenErrorCd, util.MerchantTokenErrorMsg)
		return payRes, err
	}

	logrus.Info("[PaymentService validateTrxPay] ref id date is : ", timeStampStr.Format("2006-01-02 15:04:05"))
	trxRefArg := db.GetTransactionByRefIDParams{
		IsReff:    true,
		RefID:     in.RefID,
		IsPartner: true,
		PartnerID: sql.NullInt64{
			Int64: partner.ID,
			Valid: true,
		},
	}
	refID, err := o.store.GetTransactionByRefID(ctx, trxRefArg)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("[PaymentService validateTrxPay] ", util.RefIdNotFoundMsg)
			payRes = o.PayResultErrorSet(in, util.RefIdNotFoundCd, util.RefIdNotFoundMsg)
			return payRes, err
		} else {
			logrus.Info("[PaymentService validateTrxPay] select trx table error", err)
			payRes = o.PayResultErrorSet(in, util.GeneralErrorCd, util.GeneralErrorMsg)
			return payRes, err
		}
	}

	if refID.TxID != in.TxID {
		logrus.Info("[PaymentService validateTrxPay] ", util.RefIdDifferentMsg)
		err = errors.New(util.RefIdDifferentMsg)
		payRes = o.PayResultErrorSet(in, util.RefIdDifferentCd, util.RefIdDifferentMsg)
		return payRes, err
	}

	logrus.Info("[PaymentService validateTrxPay] begin check pending trx.")
	checkBillID, err := o.store.GetTransactionPending(ctx, in.BillID)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("[PaymentService validateTrxPay] tx id not found ", err)
		} else {
			logrus.Info("[PaymentService validateTrxPay] select trx table error ", err)
			err = errors.New(util.GeneralErrorMsg)
			payRes = o.PayResultErrorSet(in, util.GeneralErrorCd, util.GeneralErrorMsg)
			return payRes, err
		}
	}

	if checkBillID.BillID != "" {
		err = errors.New(util.StillPendingTransactionMsg)
		payRes = o.PayResultErrorSet(in, util.StillPendingTransactionCd, util.StillPendingTransactionMsg)
		return payRes, err
	}

	userBalance, err := strconv.ParseFloat(user.Balance.String, 64)
	trxBalance, err := strconv.ParseFloat(trx.DeductedBalance.String, 64)
	if userBalance < trxBalance {
		logrus.Info("[PaymentService validateTrxPay] begin check pending trx.")
		err = errors.New(util.BalanceNotEnoughMsg)
		payRes = o.PayResultErrorSet(in, util.BalanceNotEnoughCd, util.BalanceNotEnoughMsg)
		return payRes, err
	}

	trxAmount, err := strconv.ParseFloat(trx.Amount.String, 64)
	if int64(trxAmount) != in.Amount {
		logrus.Info("[PaymentService validateTrxPay] trx.Amount.String : ", int64(trxAmount))
		logrus.Info("[PaymentService validateTrxPay] in.Amount : ", in.Amount)
		logrus.Info("[PaymentService validateTrxPay] amount is different.")
		err = errors.New(util.AmountDifferentMsg)
		payRes = o.PayResultErrorSet(in, util.AmountDifferentCd, util.AmountDifferentMsg)
		return payRes, err
	}

	trxAdmin, err := strconv.ParseFloat(trx.Admin.String, 64)
	if int64(trxAdmin) != in.Admin {
		logrus.Info("[PaymentService validateTrxPay] admin is different.")
		err = errors.New(util.AmountDifferentMsg)
		payRes = o.PayResultErrorSet(in, util.AmountDifferentCd, util.AmountDifferentMsg)
		return payRes, err
	}

	trxTotAmt, err := strconv.ParseFloat(trx.TotAmount.String, 64)
	if int64(trxTotAmt) != in.TotalAmount {
		logrus.Info("[PaymentService validateTrxPay] total amount is different.")
		err = errors.New(util.AmountDifferentMsg)
		payRes = o.PayResultErrorSet(in, util.AmountDifferentCd, util.AmountDifferentMsg)
		return payRes, err
	}

	if trx.Status == "0" {
		logrus.Info("[PaymentService validateTrxPay] total amount is different.")
		err = errors.New(util.AlreadySuccessMsg)
		payRes = o.PayResultErrorSet(in, util.AlreadySuccessCd, util.AlreadySuccessMsg)
		return payRes, err
	}

	return payRes, nil
}

func (o *PaymentService) InqResult(in dto.InqSetResponse) dto.InqResponse {
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

func (o *PaymentService) InqResultSet(in dto.InqRequest, resultCd string, resultMsg string) dto.InqResponse {
	inInqSetResponse := dto.InqSetResponse{
		InqData:   in,
		ResultCd:  resultCd,
		ResultMsg: resultMsg,
	}
	return o.InqResult(inInqSetResponse)
}

func (o *PaymentService) PayResult(in dto.InqSetResponse) dto.PayResponse {
	return dto.PayResponse{
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

func (o *PaymentService) PayResultErrorSet(in dto.PayRequest, resultCd string, resultMsg string) dto.PayResponse {
	return dto.PayResponse{
		TimeStamp:     in.TimeStamp,
		UserID:        in.UserID,
		RefID:         in.RefID,
		BillID:        in.BillID,
		AppName:       in.AppName,
		ProductCode:   in.ProductCode,
		MerchantToken: in.MerchantToken,
		ProductName:   "",
		Amount:        0,
		Admin:         0,
		TotalAmount:   0,
		ResultCd:      resultCd,
		ResultMsg:     resultMsg,
		TxID:          in.TxID,
	}
}
