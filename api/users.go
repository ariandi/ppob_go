package api

import (
	"database/sql"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/services"
	"github.com/ariandi/ppob_go/token"
	"github.com/ariandi/ppob_go/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

func newUserResponse(user db.User, roleUsers []dto.RoleUser) dto.UserResponse {
	return dto.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Username:       user.Username,
		Balance:        user.Balance,
		Phone:          user.Phone,
		IdentityNumber: user.IdentityNumber,
		Role:           roleUsers,
	}
}

func userRowToUserType(user db.GetUserByUsernameRow) db.User {
	return db.User{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Username:       user.Username,
		Balance:        user.Balance,
		Phone:          user.Phone,
		IdentityNumber: user.IdentityNumber,
	}
}

func (server *Server) createUsers(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	authPayload := c.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userPayload, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		Username:       req.Username,
		CreatedBy:      sql.NullInt64{Int64: userPayload.ID, Valid: true},
		Phone:          req.Phone,
		Balance:        sql.NullString{String: "0.00", Valid: true},
		IdentityNumber: req.IdentityNumber,
	}

	if req.Password != "" {
		password, err := util.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
			return
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	users, err := server.store.CreateUser(c, arg)
	if err != nil {
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
	resp := newUserResponse(users, []dto.RoleUser{})
	c.JSON(http.StatusOK, resp)
}

func (server *Server) createUsersFirst(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	//authPayload := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	userPayload, err := server.store.GetUserByUsername(c, "dbduabelas")
	if err == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user already exist"})
		return
	}
	arg := db.CreateUserParams{
		Name:           "Ariandi Nugraha",
		Email:          "dbduabelas@gmail.com",
		Username:       "dbduabelas",
		CreatedBy:      sql.NullInt64{Int64: userPayload.ID, Valid: true},
		Phone:          "081219836581",
		Balance:        sql.NullString{String: "0.00", Valid: true},
		IdentityNumber: "3201011411870003",
	}

	if req.Password != "" {
		password, err := util.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
			return
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	users, err := server.store.CreateUser(c, arg)
	if err != nil {
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

	resp := newUserResponse(users, []dto.RoleUser{})
	c.JSON(http.StatusOK, resp)
}

func (server *Server) getUser(ctx *gin.Context) {
	var req dto.GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	users, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	//if account.Owner != authPayload.Username {
	//	err := errors.New("account doesn't belong to the authenticated user")
	//	ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse(err))
	//	return
	//}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) listUsers(c *gin.Context) {
	logrus.Println("start listUsers", c.Request.Body)

	var req dto.ListUserRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
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
		}

		argUser := db.User{
			ID: user.ID,
		}
		roleUserResponse := getRoleByUser(argUser, c, server)

		u := newUserResponse(dbUser, roleUserResponse)
		resp = append(resp, u)
	}
	//resp := newUserResponse(users)
	c.JSON(http.StatusOK, resp)
}

func (server *Server) updateUsers(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
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
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
			return
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	users, err := server.store.UpdateUser(c, arg)
	if err != nil {
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

	resp := newUserResponse(users, []dto.RoleUser{})
	c.JSON(http.StatusOK, resp)
}

func (server *Server) softDeleteUser(c *gin.Context) {
	logrus.Println("start softDeleteUser", c.Request.RequestURI)

	var req dto.UpdateInactiveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	logrus.Println("start get payload")
	authPayload := c.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userPayload, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		logrus.Println("cannot find username")
		c.JSON(http.StatusBadRequest, dto.ErrorResponseString("you not allow to delete user"))
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
				c.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	resp := dto.ResponseDefault{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    "",
	}
	c.JSON(http.StatusOK, resp)
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
	logrus.Println("start login")
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Println("error validation is : ", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		logrus.Println("error get username is : ", err.Error())
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponseString("user not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password.String)
	if err != nil {
		logrus.Println("password not same : ", err.Error())
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponseString("password is incorrect"))
		return
	}

	accessToken, accessPayload, err := server.TokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.TokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
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
	//	ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
	//	return
	//}

	argUser := db.User{
		ID: user.ID,
	}
	roleUsers := getRoleByUser(argUser, ctx, server)
	userRes := userRowToUserType(user)

	rsp := dto.LoginUserResponse{
		//SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(userRes, roleUsers),
	}
	ctx.JSON(http.StatusOK, rsp)
}

func getRoleByUser(user db.User, ctx *gin.Context, server *Server) []dto.RoleUser {
	roleUserArg := db.GetRoleUserByUserIDParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 0,
	}
	roleUsers, _ := server.store.GetRoleUserByUserID(ctx, roleUserArg)
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

	return roleUserResponse
}
