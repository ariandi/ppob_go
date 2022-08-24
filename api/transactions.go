package api

import (
	"fmt"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (server *Server) createTrx(ctx *gin.Context) {
	logrus.Println("[Transactions createTrx] start.")
	var req dto.CreateTransactionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := transactionService.CreateTransactionService(req, authPayload, ctx, server.store)
	if err != nil {
		return
	}
	resp2 := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp1,
	}
	ctx.JSON(http.StatusOK, resp2)
}

func (server *Server) getTrx(ctx *gin.Context) {
	logrus.Println("[Transactions getTrx] start.")
	var req dto.GetTransactionByTxIDReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := transactionService.GetTransactionService(req, authPayload, ctx, server.store)
	if err != nil {
		return
	}

	resp2 := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp1,
	}
	ctx.JSON(http.StatusOK, resp2)
}

func (server *Server) listTrx(ctx *gin.Context) {
	logrus.Println("[Transactions listTrx] start", ctx.Request.Body)

	var req dto.ListTransactionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := transactionService.ListTransactionService(req, authPayload, ctx, server.store)
	if err != nil {
		return
	}

	resp2 := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp1,
	}
	ctx.JSON(http.StatusOK, resp2)
}

func (server *Server) updateTrx(ctx *gin.Context) {
	logrus.Println("[Transactions updateTrx] start.")
	var req dto.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := transactionService.UpdateTransactionService(req, authPayload, ctx, server.store)
	if err != nil {
		return
	}

	resp2 := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp1,
	}
	ctx.JSON(http.StatusOK, resp2)
}

func (server *Server) softDeleteTrx(ctx *gin.Context) {
	logrus.Println("[Transactions softDeleteTrx] start.")
	var req dto.UpdateInactiveTransactionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		logrus.Println("[Transactions softDeleteTrx] error validation.")
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("[Transactions softDeleteTrx] start get payload")
	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	err := transactionService.SoftDeleteTransactionService(req, authPayload, ctx, server.store)
	if err != nil {
		return
	}

	resp := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    "",
	}
	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) inquiry(ctx *gin.Context) {
	logrus.Println("[Transactions inquiry] start.")
	var req dto.InqRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errs, _ := err.(validator.ValidationErrors)
		logrus.Info("ok", errs)
		for _, v := range errs {
			field := v.Field()
			tag := v.Tag()

			errMsg := fmt.Sprintf("%v: %v", field, tag)
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponseString(errMsg))
			break
		}
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := transactionService.InqService(req, authPayload, ctx, server.store)
	if err != nil {
		ctx.JSON(http.StatusOK, resp1)
		return
	}

	resp2 := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp1,
	}
	ctx.JSON(http.StatusOK, resp2)
}
