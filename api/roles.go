package api

import (
	"database/sql"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RoleResponse(role db.Role) dto.RoleRes {
	return dto.RoleRes{
		ID:    role.ID,
		Name:  role.Name,
		Level: role.Level,
	}
}

func (server *Server) createRole(c *gin.Context) {
	logrus.Println("[Roles createRole] start.")
	var req dto.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := c.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userPayload, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Println("[Roles createRole] : user not found")
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	err = userService.ValidateUserRole(userPayload)
	if err != nil {
		logrus.Println("createRoleUsers, ValidateUserRole : ", err)
		c.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return
	}

	arg := db.CreateRoleParams{
		Name:      req.Name,
		Level:     req.Level,
		CreatedBy: sql.NullInt64{Int64: userPayload.ID, Valid: true},
	}

	role, err := server.store.CreateRole(c, arg)
	if err != nil {
		logrus.Println("[Roles createRole] : error create role users")
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				c.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	resp1 := RoleResponse(role)
	resp2 := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp1,
	}
	c.JSON(http.StatusOK, resp2)
}

func (server *Server) getRole(ctx *gin.Context) {
	logrus.Println("[Roles getRole] start.")
	var req dto.GetRoleReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userPayload, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Println("[Roles getRole] user not found")
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	err = userService.ValidateUserRole(userPayload)
	if err != nil {
		logrus.Println("[Roles getRole] createRoleUsers, ValidateUserRole : ", err)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return
	}

	role, err := server.store.GetRole(ctx, req.ID)
	if err != nil {
		logrus.Println("[Roles getRole] start getRoleUserByUserID.")
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	resp1 := RoleResponse(role)
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
	userPayload, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Println("[Roles listRole] start createRoleUsers : user not found")
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	err = userService.ValidateUserRole(userPayload)
	if err != nil {
		logrus.Println("[Roles listRole] createRoleUsers, ValidateUserRole : ", err)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return
	}

	arg := db.ListRoleParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roles, err := server.store.ListRole(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	var resp1 []dto.RoleRes
	for _, role := range roles {
		u := RoleResponse(role)
		resp1 = append(resp1, u)
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
	userPayload, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Println("[Roles updateRole] start createRoleUsers : user not found")
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	err = userService.ValidateUserRole(userPayload)
	if err != nil {
		logrus.Println("[Roles updateRole] createRoleUsers, ValidateUserRole : ", err)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return
	}

	arg := db.UpdateRoleParams{
		ID:        req.ID,
		Name:      req.Name,
		Level:     req.Level,
		UpdatedBy: sql.NullInt64{Int64: userPayload.ID, Valid: true},
	}

	role, err := server.store.UpdateRole(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	resp1 := RoleResponse(role)
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("[Roles softDeleteRole] start get payload")
	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userPayload, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Println("start createRoleUsers : user not found")
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	err = userService.ValidateUserRole(userPayload)
	if err != nil {
		logrus.Println("[Roles softDeleteRole] createRoleUsers, ValidateUserRole : ", err)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return
	}

	arg := db.UpdateInactiveRoleParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userPayload.ID, Valid: true},
	}

	_, err = server.store.UpdateInactiveRole(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	resp := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    "",
	}
	ctx.JSON(http.StatusOK, resp)
}
