package api

import (
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/services"
	"github.com/ariandi/ppob_go/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (server *Server) createUsers(c *gin.Context) {
	logrus.Println("[Users createUsers] start.")
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := c.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := userService.CreateUserService(req, authPayload, c, server.store)
	if err != nil {
		return
	}
	resp2 := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp1,
	}
	c.JSON(http.StatusOK, resp2)
}

func (server *Server) createUsersFirst(c *gin.Context) {
	logrus.Println("[Users createUsersFirst] start.")
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := userService.CreateUserFirstService(req, c, server.store)
	if err == nil {
		return
	}

	resp2 := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    resp1,
	}
	c.JSON(http.StatusOK, resp2)
}

func (server *Server) getUser(ctx *gin.Context) {
	logrus.Println("[Users getUser] start.")
	var req dto.GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := userService.GetUserService(req, authPayload, ctx, server.store)
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

func (server *Server) listUsers(ctx *gin.Context) {
	logrus.Println("start listUsers", ctx.Request.Body)

	var req dto.ListUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := userService.ListUserService(req, authPayload, ctx, server.store)
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

func (server *Server) updateUsers(ctx *gin.Context) {
	var req dto.UpdateUserRequest
	var reqID dto.UpdateUserIDRequest

	if err := ctx.ShouldBindUri(&reqID); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	req.ID = reqID.ID

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := userService.UpdateUserService(req, authPayload, ctx, server.store)
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

func (server *Server) softDeleteUser(ctx *gin.Context) {
	logrus.Println("start softDeleteUser", ctx.Request.RequestURI)

	var req dto.UpdateInactiveUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Println("error validation softDeleteUser", err)
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("start get payload")
	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	err := userService.SoftDeleteUserService(req, authPayload, ctx, server.store)
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

func (server *Server) testRedisMq(ctx *gin.Context) {
	var userService services.UserService
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	res, err := userService.TestRedisMq(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) loginUser(ctx *gin.Context) {
	logrus.Println("[User loginUser] start login")
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Println("[User loginUser] error validation is : ", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	rsp, err := userService.LoginUserService(req, server.TokenMaker, ctx, server.store, server.config)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}
