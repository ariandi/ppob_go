package api

import (
	"fmt"
	"github.com/ariandi/ppob_go/dto"
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

	resp1, err := transactionService.CreateTransactionService(ctx, req)
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

	resp1, err := transactionService.GetTransactionService(ctx, req)
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

	resp1, err := transactionService.ListTransactionService(ctx, req)
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

func (server *Server) paymentExport(ctx *gin.Context) {
	logrus.Println("[Transactions paymentExport] start", ctx.Request.Body)

	var req dto.ListTransactionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	f, err := transactionService.ExportTransaction(ctx, req)
	if err != nil {
		return
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment;filename=userInputData.xlsx")
	ctx.Header("File-Name", "userInputData.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	errExport := f.Write(ctx.Writer)
	if errExport != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(errExport))
		return
	}
}

func (server *Server) updateTrx(ctx *gin.Context) {
	logrus.Println("[Transactions updateTrx] start.")
	var req dto.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := transactionService.UpdateTransactionService(ctx, req)
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
	err := transactionService.SoftDeleteTransactionService(ctx, req)
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

	resp1, err := transactionService.InqService(ctx, req)
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

func (server *Server) deposit(ctx *gin.Context) {
	logrus.Println("[Transactions deposit] start.")
	var req dto.DepositRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errs, _ := err.(validator.ValidationErrors)
		logrus.Info("[Transactions deposit] validation error ", errs)
		for _, v := range errs {
			field := v.Field()
			tag := v.Tag()

			errMsg := fmt.Sprintf("%v: %v", field, tag)
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponseString(errMsg))
			break
		}
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponseString("error in validation"))
		return
	}

	resp1, err := transactionService.DepositService(ctx, req)
	if err != nil {
		logrus.Info("[Transactions deposit] error deposit service ", err)
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

func (server *Server) depositApprove(ctx *gin.Context) {
	logrus.Println("[Transactions deposit] start.")
	var req dto.DepositApproveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errs, _ := err.(validator.ValidationErrors)
		logrus.Info("[Transactions deposit] validation error ", errs)
		for _, v := range errs {
			field := v.Field()
			tag := v.Tag()

			errMsg := fmt.Sprintf("%v: %v", field, tag)
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponseString(errMsg))
			break
		}
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponseString("error in validation"))
		return
	}

	resp1, err := transactionService.DepositApproveService(ctx, req)
	if err != nil {
		logrus.Info("[Transactions deposit] error deposit service ", err)
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
