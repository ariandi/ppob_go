package api

import (
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (server *Server) createProduct(ctx *gin.Context) {
	logrus.Println("[Products createProduct] start.")
	var req dto.CreateProductReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := productService.CreateProductService(req, authPayload, ctx, server.store)
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

func (server *Server) getProduct(ctx *gin.Context) {
	logrus.Println("[Products getProduct] start.")
	var req dto.GetProductReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := productService.GetProductService(req, authPayload, ctx, server.store)
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

func (server *Server) listProduct(ctx *gin.Context) {
	logrus.Println("[Products listProduct] start", ctx.Request.Body)

	var req dto.ListProductRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := productService.ListProductService(req, authPayload, ctx, server.store)
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

func (server *Server) updateProduct(ctx *gin.Context) {
	logrus.Println("[Products updateProduct] start.")
	var req dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := productService.UpdateProductService(req, authPayload, ctx, server.store)
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

func (server *Server) softDeleteProduct(ctx *gin.Context) {
	logrus.Println("[Products softDeleteProduct] start.")
	var req dto.UpdateInactiveProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		logrus.Println("[Products softDeleteProduct] error validation.")
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("[Products softDeleteProduct] start get payload")
	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	err := productService.SoftDeleteProductService(req, authPayload, ctx, server.store)
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
