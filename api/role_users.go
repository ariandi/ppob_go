package api

import (
	"database/sql"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/ariandi/ppob_go/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

func newRoleUserResponse(roleUser db.RoleUser) dto.CreateRoleUserRes {
	return dto.CreateRoleUserRes{
		ID:     roleUser.ID,
		UserID: roleUser.UserID,
		RoleID: roleUser.RoleID,
	}
}

func (server *Server) createRoleUsers(c *gin.Context) {
	logrus.Println("start createRoleUsers")
	var req dto.CreateRoleUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userPayload, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Println("start createRoleUsers : user not found")
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
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
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//var roleUser []dto.RoleUser
	//roleUser = []
	resp := newRoleUserResponse(roleUsers)
	c.JSON(http.StatusOK, resp)
}

func (server *Server) getRoleUserByUserID(ctx *gin.Context) {
	logrus.Println("start getRoleUserByUserID.")
	var req dto.GetRoleUserByUserID
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetRoleUserByUserIDParams{
		UserID: req.UserID,
		Limit:  1,
		Offset: 0,
	}

	users, err := server.store.GetRoleUserByUserID(ctx, arg)
	if err != nil {
		logrus.Println("start getRoleUserByUserID.")
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	//if account.Owner != authPayload.Username {
	//	err := errors.New("account doesn't belong to the authenticated user")
	//	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	//	return
	//}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) listRoleUsers(c *gin.Context) {
	logrus.Println("start listUsers", c.Request.Body)

	var req dto.ListUserRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListUserParams{
		//Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var resp []dto.UserResponse
	for _, user := range users {
		dbUser := db.User{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			Username:       user.Username,
			Balance:        user.Balance,
			Phone:          user.Phone,
			IdentityNumber: user.IdentityNumber,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			CreatedBy:      user.CreatedBy,
			UpdatedBy:      user.UpdatedBy,
			//RoleID:         user.RoleID,
		}

		roleUserArg := db.GetRoleUserByUserIDParams{
			UserID: user.ID,
			Limit:  5,
			Offset: 0,
		}
		roleUsers, err := server.store.GetRoleUserByUserID(c, roleUserArg)
		if err != nil {
			logrus.Println("error GetRoleUserByUserID is : ", err.Error())
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		var roleUserResponse []dto.RoleUser
		for _, roleUser := range roleUsers {
			if roleUser.UserID == user.ID {
				roleUserDto := dto.RoleUser{
					ID:     roleUser.ID,
					RoleID: roleUser.RoleID,
					UserID: roleUser.UserID,
				}
				roleUserResponse = append(roleUserResponse, roleUserDto)
			}
		}

		if roleUserResponse == nil {
			roleUserResponse = []dto.RoleUser{}
		}

		u := newUserResponse(dbUser, roleUserResponse)
		resp = append(resp, u)
	}
	//resp := newUserResponse(users)
	c.JSON(http.StatusOK, resp)
}

func (server *Server) updateRoleUsers(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:                req.ID,
		SetName:           false,
		Name:              req.Name,
		SetPhone:          false,
		Phone:             req.Phone,
		SetIdentityNumber: false,
		IdentityNumber:    req.IdentityNumber,
		SetPassword:       false,
		UpdatedBy:         sql.NullInt64{Int64: 1, Valid: true},
	}

	if req.Name != "" {
		arg.SetName = true
	}

	if req.Phone != "" {
		arg.SetPhone = true
	}

	if req.IdentityNumber != "" {
		arg.SetIdentityNumber = true
	}

	if req.Password != "" {
		arg.SetPassword = true
		password, err := util.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	users, err := server.store.UpdateUser(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(users, []dto.RoleUser{})
	c.JSON(http.StatusOK, resp)
}

func (server *Server) softDeleteRoleUser(c *gin.Context) {
	logrus.Println("start softDeleteUser", c.Request.RequestURI)

	var req dto.UpdateInactiveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	logrus.Println("start get payload")
	authPayload := c.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userPayload, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		logrus.Println("cannot find username")
		c.JSON(http.StatusBadRequest, errorResponseString("you not allow to delete user"))
		return
	}
	arg := db.UpdateInactiveUserParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userPayload.ID, Valid: true},
	}

	_, err = server.store.UpdateInactiveUser(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    "",
	}
	c.JSON(http.StatusOK, resp)
}

func (server *Server) loginRoleUser(ctx *gin.Context) {
	logrus.Println("start login")
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Println("error validation is : ", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		logrus.Println("error get username is : ", err.Error())
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponseString("user not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password.String)
	if err != nil {
		logrus.Println("password not same : ", err.Error())
		ctx.JSON(http.StatusUnauthorized, errorResponseString("password is incorrect"))
		return
	}

	accessToken, accessPayload, err := server.TokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.TokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
	//	ID:           refreshPayload.ID,
	//	Username:     user.Username,
	//	RefreshToken: refreshToken,
	//	UserAgent:    ctx.Request.UserAgent(),
	//	ClientIp:     ctx.ClientIP(),
	//	IsBlocked:    false,
	//	ExpiresAt:    refreshPayload.ExpiredAt,
	//})
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//}

	rsp := dto.LoginUserResponse{
		//SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user, []dto.RoleUser{}),
	}
	ctx.JSON(http.StatusOK, rsp)
}
