package api

import (
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (server *Server) createRole(ctx *gin.Context) {
	logrus.Println("[Roles createRole] start.")
	var req dto.CreateRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := roleService.CreateRoleService(req, authPayload, ctx, server.store)
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

func (server *Server) getRole(ctx *gin.Context) {
	logrus.Println("[Roles getRole] start.")
	var req dto.GetRoleReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := roleService.GetRoleService(req, authPayload, ctx, server.store)
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

func (server *Server) listRole(ctx *gin.Context) {
	logrus.Println("[Roles listRole] start listRoleUsers", ctx.Request.Body)

	var req dto.ListRoleRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := roleService.ListRoleService(req, authPayload, ctx, server.store)
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

func (server *Server) updateRole(ctx *gin.Context) {
	logrus.Println("[Roles updateRole] start.")
	var req dto.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := roleService.UpdateRoleService(req, authPayload, ctx, server.store)
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

func (server *Server) softDeleteRole(ctx *gin.Context) {
	logrus.Println("[Roles softDeleteRole] start.")
	var req dto.UpdateInactiveRoleRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		logrus.Println("[Roles softDeleteRole] error validation.")
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("[Roles softDeleteRole] start get payload")
	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	err := roleService.SoftDeleteRoleService(req, authPayload, ctx, server.store)
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
