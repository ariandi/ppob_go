package api

import (
	"fmt"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (server *Server) createUsers(ctx *gin.Context) {
	logrus.Println("[Users createUsers] start.")
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := userService.CreateUserService(ctx, req)
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

func (server *Server) createUsersFirst(ctx *gin.Context) {
	logrus.Println("[Users createUsersFirst] start.")
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	resp1, err := userService.CreateUserFirstService(ctx, req)
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

func (server *Server) getUser(ctx *gin.Context) {
	logrus.Println("[Users getUser] start.")
	var req dto.GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	//authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	resp1, err := userService.GetUserService(ctx, req)
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

	resp1, err := userService.ListUserService(ctx, req)
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
	resp1, err := userService.UpdateUserService(ctx, req)
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
	err := userService.SoftDeleteUserService(ctx, req)
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
		errs, _ := err.(validator.ValidationErrors)
		logrus.Info("ok", errs)
		for _, v := range errs {
			field := v.Field()
			tag := v.Tag()

			errMsg := fmt.Sprintf("%v: %v", field, tag)
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponseString(errMsg))
			break
		}

		if len(errs) == 0 {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		}
		return
	}

	rsp, err := userService.LoginUserService(ctx, req)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}
