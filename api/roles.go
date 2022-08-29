package api

import (
	"github.com/ariandi/ppob_go/dto"
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

	resp1, err := roleService.CreateRoleService(ctx, req)
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

	resp1, err := roleService.GetRoleService(ctx, req)
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

	resp1, err := roleService.ListRoleService(ctx, req)
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

	resp1, err := roleService.UpdateRoleService(ctx, req)
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
	err := roleService.SoftDeleteRoleService(ctx, req)
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
