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

func newRoleUserResponse(roleUser db.RoleUser) dto.RoleUserRes {
	return dto.RoleUserRes{
		ID:     roleUser.ID,
		UserID: roleUser.UserID,
		RoleID: roleUser.RoleID,
	}
}

func (server *Server) createRoleUsers(c *gin.Context) {
	logrus.Println("start createRoleUsers")
	var req dto.CreateRoleUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := c.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userPayload, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Println("start createRoleUsers : user not found")
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

	arg := db.CreateRoleUserParams{
		UserID:    req.UserID,
		RoleID:    req.RoleID,
		CreatedBy: sql.NullInt64{Int64: userPayload.ID, Valid: true},
	}

	roleUsers, err := server.store.CreateRoleUser(c, arg)
	if err != nil {
		logrus.Println("start createRoleUsers : error create role users")
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

	//var roleUser []dto.RoleUser
	//roleUser = []
	resp := newRoleUserResponse(roleUsers)
	c.JSON(http.StatusOK, resp)
}

//func (server *Server) getRoleUserByID(ctx *gin.Context) {
//	logrus.Println("start getRoleUserByID.")
//	var req dto.GetRoleUserByUserID
//	if err := ctx.ShouldBindUri(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
//		return
//	}
//
//	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
//	userPayload, err := server.store.GetUserByUsername(ctx, authPayload.Username)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			logrus.Println("start createRoleUsers : user not found")
//			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
//		return
//	}
//
//	roleUsers, err := server.store.GetRoleUserByID(ctx, req.UserID)
//	if err != nil {
//		logrus.Println("error get GetRoleUserByID.")
//		if err == sql.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
//			return
//		}
//
//		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
//		return
//	}
//
//	var resp []dto.RoleUserRes
//	for _, roleUser := range roleUsers {
//		resArg := newRoleUserResponse(roleUser)
//		resp = append(resp, resArg)
//	}
//	ctx.JSON(http.StatusOK, roleUsers)
//}

func (server *Server) getRoleUserByUserID(ctx *gin.Context) {
	logrus.Println("start getRoleUserByUserID.")
	var req dto.GetRoleUserByUserID
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

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
		logrus.Println("createRoleUsers, ValidateUserRole : ", err)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return
	}

	arg := db.GetRoleUserByUserIDParams{
		UserID: req.UserID,
		Limit:  5,
		Offset: 0,
	}

	roleUsers, err := server.store.GetRoleUserByUserID(ctx, arg)
	if err != nil {
		logrus.Println("start getRoleUserByUserID.")
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	var resp []dto.RoleUserRes
	for _, roleUser := range roleUsers {
		resArg := newRoleUserResponse(roleUser)
		resp = append(resp, resArg)
	}
	ctx.JSON(http.StatusOK, roleUsers)
}

func (server *Server) listRoleUsers(ctx *gin.Context) {
	logrus.Println("start listRoleUsers", ctx.Request.Body)

	var req dto.ListRoleUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

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
		logrus.Println("createRoleUsers, ValidateUserRole : ", err)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return
	}

	arg := db.ListRoleUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roleUsers, err := server.store.ListRoleUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	var resp []dto.RoleUserRes
	for _, roleUser := range roleUsers {
		dbRoleUser := db.RoleUser{
			ID:     roleUser.ID,
			RoleID: roleUser.RoleID,
			UserID: roleUser.UserID,
		}

		u := newRoleUserResponse(dbRoleUser)
		resp = append(resp, u)
	}
	//resp := newUserResponse(users)
	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) updateRoleUsers(ctx *gin.Context) {
	var req dto.UpdateRoleUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

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
		logrus.Println("createRoleUsers, ValidateUserRole : ", err)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return
	}

	arg := db.UpdateRoleUserParams{
		ID:        req.ID,
		UserID:    req.UserID,
		RoleID:    req.RoleID,
		UpdatedBy: sql.NullInt64{Int64: userPayload.ID, Valid: true},
	}

	roleUser, err := server.store.UpdateRoleUser(ctx, arg)
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

	resp := newRoleUserResponse(roleUser)
	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) softDeleteRoleUser(ctx *gin.Context) {
	logrus.Println("[role_user softDeleteRoleUser] start. softDeleteRoleUser", ctx.Request.RequestURI)
	var req dto.UpdateInactiveROleUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("[role_user softDeleteRoleUser] start get payload")
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
		logrus.Println("createRoleUsers, ValidateUserRole : ", err)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return
	}

	arg := db.UpdateInactiveRoleUserParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userPayload.ID, Valid: true},
	}

	_, err = server.store.UpdateInactiveRoleUser(ctx, arg)
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
