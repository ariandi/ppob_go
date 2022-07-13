package api

import (
	"database/sql"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

func newUserResponse(user db.User) dto.CreateUserResponse {
	return dto.CreateUserResponse{
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
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		Username:       req.Username,
		CreatedBy:      sql.NullInt64{Int64: 1, Valid: true},
		Phone:          req.Phone,
		Balance:        sql.NullString{String: "0.00", Valid: true},
		IdentityNumber: req.IdentityNumber,
	}

	if req.Password != "" {
		password, err := util.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	users, err := server.store.CreateUser(c, arg)
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

	resp := newUserResponse(users)
	c.JSON(http.StatusOK, resp)
}

func (server *Server) getUser(ctx *gin.Context) {
	var req dto.GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
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

func (server *Server) listUsers(c *gin.Context) {
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

	c.JSON(http.StatusOK, users)
}

func (server *Server) updateUsers(c *gin.Context) {
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

	resp := newUserResponse(users)
	c.JSON(http.StatusOK, resp)
}
