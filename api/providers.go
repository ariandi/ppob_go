package api

import (
	"github.com/ariandi/ppob_go/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (server *Server) createProvider(ctx *gin.Context) {
	logrus.Println("[Providers createProvider] start.")
	var req dto.CreateProviderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := providerService.CreateProviderService(ctx, req)
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

func (server *Server) getProvider(ctx *gin.Context) {
	logrus.Println("[Providers getProvider] start.")
	var req dto.GetProviderReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := providerService.GetProviderService(ctx, req)
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

func (server *Server) listProvider(ctx *gin.Context) {
	logrus.Println("[Providers listProvider] start", ctx.Request.Body)

	var req dto.ListProviderRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := providerService.ListProviderService(ctx, req)
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

func (server *Server) updateProvider(ctx *gin.Context) {
	logrus.Println("[Providers updateProvider] start.")
	var req dto.UpdateProviderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := providerService.UpdateProviderService(ctx, req)
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

func (server *Server) softDeleteProvider(ctx *gin.Context) {
	logrus.Println("[Providers softDeleteProvider] start.")
	var req dto.UpdateInactiveProviderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		logrus.Println("[Providers softDeleteProvider] error validation.")
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("[Providers softDeleteProvider] start get payload")
	err := providerService.SoftDeleteProviderService(ctx, req)
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
