package api

import (
	"github.com/ariandi/ppob_go/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (server *Server) createMedia(ctx *gin.Context) {
	logrus.Println("[Categories createMedia] start.")
	var req dto.CreateMediaReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := mediaService.CreateMediaService(ctx, req)
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

func (server *Server) getMedia(ctx *gin.Context) {
	logrus.Println("[Categories getMedia] start.")
	var req dto.GetMediaReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := mediaService.GetMediaService(ctx, req)
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

func (server *Server) listMedia(ctx *gin.Context) {
	logrus.Println("[Categories listMedia] start", ctx.Request.Body)

	var req dto.ListMediaRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := mediaService.ListMediaService(ctx, req)
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

func (server *Server) updateMedia(ctx *gin.Context) {
	logrus.Println("[Categories updateMedia] start.")
	var req dto.UpdateMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := mediaService.UpdateMediaService(ctx, req)
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

func (server *Server) softDeleteMedia(ctx *gin.Context) {
	logrus.Println("[Categories softDeleteMedia] start.")
	var req dto.UpdateInactiveMediaRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		logrus.Println("[Categories softDeleteMedia] error validation.")
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("[Categories softDeleteMedia] start get payload")
	err := mediaService.SoftDeleteMediaService(ctx, req)
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
