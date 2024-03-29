package api

import (
	"github.com/ariandi/ppob_go/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (server *Server) createSelling(ctx *gin.Context) {
	logrus.Println("[Selling createSelling] start.")
	var req dto.CreateSellingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := sellingService.CreateSellingService(ctx, req)
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

func (server *Server) getSelling(ctx *gin.Context) {
	logrus.Println("[Selling getSelling] start.")
	var req dto.GetSellingReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := sellingService.GetSellingService(ctx, req)
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

func (server *Server) listSelling(ctx *gin.Context) {
	logrus.Println("[Selling listSelling] start", ctx.Request.Body)

	var req dto.ListSellingRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := sellingService.ListSellingService(ctx, req)
	if err != nil {
		return
	}

	logrus.Info(resp1)
	resp2 := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp1,
	}
	ctx.JSON(http.StatusOK, resp2)
}

func (server *Server) updateSelling(ctx *gin.Context) {
	logrus.Println("[Selling updateSelling] start.")
	var req dto.UpdateSellingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := sellingService.UpdateSellingService(ctx, req)
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

func (server *Server) softDeleteSelling(ctx *gin.Context) {
	logrus.Println("[Selling softDeleteSelling] start.")
	var req dto.UpdateInactiveSellingRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		logrus.Println("[Selling softDeleteSelling] error validation.")
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("[Selling softDeleteSelling] start get payload")
	err := sellingService.SoftDeleteSellingService(ctx, req)
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
