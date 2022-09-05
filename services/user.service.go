package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adjust/rmq/v4"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/ariandi/ppob_go/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type UserInterface interface {
	TestRedisMq(msg dto.LoginUserRequest) ([]string, error)
	CreateUserService(ctx *gin.Context, in dto.CreateUserRequest) (dto.UserResponse, error)
	CreateUserFirstService(ctx *gin.Context, in dto.CreateUserRequest) (dto.UserResponse, error)
	GetUserService(ctx *gin.Context, in dto.GetUserRequest) (dto.UserResponse, error)
	ListUserService(ctx *gin.Context, in dto.ListUserRequest) ([]dto.UserResponse, error)
	UpdateUserService(ctx *gin.Context, in dto.UpdateUserRequest) (dto.UserResponse, error)
	SoftDeleteUserService(ctx *gin.Context, in dto.UpdateInactiveUserRequest) error
	LoginUserService(ctx *gin.Context, in dto.LoginUserRequest) (dto.LoginUserResponse, error)
	newUserResponse(user db.User, roleUsers []dto.RoleUser) dto.UserResponse
	userRowToUserType(user db.GetUserByUsernameRow) db.User
	ValidateUserRole(user db.GetUserByUsernameRow) error
	getRoleByUser(ctx *gin.Context, in db.User) []dto.RoleUser
	validator(ctx *gin.Context, authPayload *token.Payload) (db.GetUserByUsernameRow, error)
}

// UserService is
type UserService struct {
	Store      db.Store
	Config     util.Config
	TokenMaker token.Maker
}

var userService *UserService
var redisConn rmq.Connection

const (
	AuthorizationPayloadKey = "authorization_payload"
)

// GetUserService is
func GetUserService(config util.Config, store db.Store, TokenMaker token.Maker) UserInterface {

	if userService == nil {

		userService = &UserService{
			Store:      store,
			Config:     config,
			TokenMaker: TokenMaker,
		}

		redisDb, errRedisDb := strconv.Atoi(config.RedisDB)
		if errRedisDb != nil {
			fmt.Println("=======================================")
			fmt.Println("Cannot get redis db config : ", errRedisDb)
			fmt.Println("=======================================")
		}

		var err error
		redisConn, err = rmq.OpenConnection("redisService", "tcp", config.RedisUrl, redisDb, nil)

		if err != nil {
			fmt.Println("=======================================")
			fmt.Println("Error connect Redis : ", err)
			fmt.Println("=======================================")
		}
	}

	return userService
}

func (o *UserService) TestRedisMq(msg dto.LoginUserRequest) ([]string, error) {

	ret := []string{}
	queueName := "test_"
	redisQueue, err := redisConn.OpenQueue(queueName)
	if err != nil {
		return ret, err
	}

	//messageItem := msg
	byt, err := json.Marshal(msg)
	if err != nil {
		return ret, err
	}

	err = redisQueue.Publish(string(byt))
	ret = append(ret, "Success")

	return ret, nil
}

func (o *UserService) CreateUserService(ctx *gin.Context, in dto.CreateUserRequest) (dto.UserResponse, error) {
	logrus.Println("[UserService CreateUserService] start.")
	var result dto.UserResponse

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	user, err := o.validator(ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.CreateUserParams{
		Name:           in.Name,
		Email:          in.Email,
		Username:       in.Username,
		CreatedBy:      sql.NullInt64{Int64: user.ID, Valid: true},
		Phone:          in.Phone,
		Balance:        sql.NullString{String: "0.00", Valid: true},
		IdentityNumber: in.IdentityNumber,
	}

	if in.Password != "" {
		password, errHash := util.HashPassword(in.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
			return result, errHash
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	out, err := o.Store.CreateUserTx(ctx, arg, authPayload, in.RoleID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return result, err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	return out, nil
}

func (o *UserService) CreateUserFirstService(ctx *gin.Context, in dto.CreateUserRequest) (dto.UserResponse, error) {
	logrus.Println("[UserService CreateUserFirstService] start.")
	var result dto.UserResponse

	if in.Username != "dbduabelas" {
		logrus.Println("[UserService CreateUserFirstService] error username not allow.")
		ctx.JSON(http.StatusNotFound, dto.ErrorResponseString("error username not allow"))
		return result, errors.New("error username not allow")
	}

	userPayload, err := o.Store.GetUserByUsername(ctx, "dbduabelas")
	if err == nil {
		logrus.Println("[UserService CreateUserFirstService] user already exist.")
		ctx.JSON(http.StatusNotFound, dto.ErrorResponseString("user already exist"))
		return result, errors.New("user already exist")
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

	if in.Password != "" {
		password, errPswd := util.HashPassword(in.Password)
		if errPswd != nil {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(errPswd))
			return result, errPswd
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	users, err := o.Store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return result, err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	out := o.newUserResponse(users, []dto.RoleUser{})
	return out, nil
}

func (o *UserService) GetUserService(ctx *gin.Context, in dto.GetUserRequest) (dto.UserResponse, error) {
	logrus.Println("[UserService GetUserService] start.")
	var result dto.UserResponse

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := o.validator(ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	user, err := o.Store.GetUser(ctx, in.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return result, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	argUser := db.User{
		ID: user.ID,
	}
	roleUsers := o.getRoleByUser(ctx, argUser)
	userArg := db.User{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Username:       user.Username,
		Balance:        user.Balance,
		Phone:          user.Phone,
		IdentityNumber: user.IdentityNumber,
	}
	out := o.newUserResponse(userArg, roleUsers)

	return out, nil
}

func (o *UserService) ListUserService(ctx *gin.Context, in dto.ListUserRequest) ([]dto.UserResponse, error) {
	logrus.Println("[UserService ListUserService] start.")
	var result []dto.UserResponse

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := o.validator(ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.ListUserParams{
		Limit:  in.PageSize,
		Offset: (in.PageID - 1) * in.PageSize,
		RoleID: in.RoleID,
	}

	if in.RoleID > 0 {
		arg.IsRole = true
	}

	users, err := o.Store.ListUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	var out []dto.UserResponse
	for _, user := range users {
		dbUser := db.User{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			Username:       user.Username,
			Balance:        user.Balance,
			Phone:          user.Phone,
			BankCode:       user.BankCode,
			IdentityNumber: user.IdentityNumber,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			CreatedBy:      user.CreatedBy,
			UpdatedBy:      user.UpdatedBy,
		}

		argUser := db.User{
			ID: user.ID,
		}
		roleUserResponse := o.getRoleByUser(ctx, argUser)

		u := o.newUserResponse(dbUser, roleUserResponse)
		out = append(out, u)
	}
	return out, nil
}

func (o *UserService) UpdateUserService(ctx *gin.Context, in dto.UpdateUserRequest) (dto.UserResponse, error) {
	logrus.Println("[UserService UpdateUserService] start.")
	var result dto.UserResponse

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	user, err := o.validator(ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.UpdateUserParams{
		ID:             in.ID,
		Email:          in.Email,
		Name:           in.Name,
		Phone:          in.Phone,
		IdentityNumber: in.IdentityNumber,
		UpdatedBy:      sql.NullInt64{Int64: user.ID, Valid: true},
	}

	if in.Password != "" {
		password, errPasswd := util.HashPassword(in.Password)
		if errPasswd != nil {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(errPasswd))
			return result, errPasswd
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	out, err := o.Store.UpdateUserTx(ctx, arg, authPayload, in.RoleID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return result, err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	return out, nil
}

func (o *UserService) SoftDeleteUserService(ctx *gin.Context, in dto.UpdateInactiveUserRequest) error {
	logrus.Println("[UserService SoftDeleteUserService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	user, err := o.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveUserParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: user.ID, Valid: true},
	}

	_, err = o.Store.UpdateInactiveUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return err
	}

	return nil
}

func (o *UserService) LoginUserService(ctx *gin.Context, in dto.LoginUserRequest) (dto.LoginUserResponse, error) {
	logrus.Println("[UserService LoginUserService] start.")
	var result dto.LoginUserResponse

	user, err := o.Store.GetUserByUsername(ctx, in.Username)
	if err != nil {
		logrus.Println("[UserService loginUser] error get username is : ", err.Error())
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponseString("user not found"))
			return result, err
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	err = util.CheckPassword(in.Password, user.Password.String)
	if err != nil {
		logrus.Println("[UserService loginUser] password not same : ", err.Error())
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponseString("password is incorrect"))
		return result, err
	}

	accessToken, accessPayload, err := o.TokenMaker.CreateToken(
		user.Username,
		o.Config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	refreshToken, refreshPayload, err := o.TokenMaker.CreateToken(
		user.Username,
		o.Config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	argUser := db.User{
		ID: user.ID,
	}
	roleUsers := o.getRoleByUser(ctx, argUser)
	userRes := o.userRowToUserType(user)

	out := dto.LoginUserResponse{
		//SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  o.newUserResponse(userRes, roleUsers),
	}
	return out, nil
}

func (o *UserService) newUserResponse(user db.User, roleUsers []dto.RoleUser) dto.UserResponse {
	return dto.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Username:       user.Username,
		Balance:        user.Balance,
		Phone:          user.Phone,
		IdentityNumber: user.IdentityNumber,
		BankCode:       user.BankCode.Int64,
		Role:           roleUsers,
	}
}

func (o *UserService) userRowToUserType(user db.GetUserByUsernameRow) db.User {
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

func (o *UserService) ValidateUserRole(user db.GetUserByUsernameRow) error {
	if user.RoleID.Int64 != 1 {
		return fmt.Errorf("you not allow to access this service")
	}

	return nil
}

func (o *UserService) getRoleByUser(ctx *gin.Context, in db.User) []dto.RoleUser {
	roleUserArg := db.GetRoleUserByUserIDParams{
		UserID: in.ID,
		Limit:  5,
		Offset: 0,
	}
	roleUsers, _ := o.Store.GetRoleUserByUserID(ctx, roleUserArg)
	var roleUserResponse []dto.RoleUser
	for _, roleUser := range roleUsers {
		if roleUser.UserID == in.ID {
			roleUserDto := dto.RoleUser{
				ID:        roleUser.ID,
				RoleID:    roleUser.RoleID,
				UserID:    roleUser.UserID,
				CreatedAt: roleUser.CreatedAt,
				CreatedBy: roleUser.CreatedBy,
				UpdatedBy: roleUser.UpdatedBy,
				UpdatedAt: roleUser.UpdatedAt,
			}
			roleUserResponse = append(roleUserResponse, roleUserDto)
		}
	}

	if roleUserResponse == nil {
		roleUserResponse = []dto.RoleUser{}
	}

	return roleUserResponse
}

func (o *UserService) validator(ctx *gin.Context, authPayload *token.Payload) (db.GetUserByUsernameRow, error) {
	logrus.Println("[UserService validator] start.")
	var res db.GetUserByUsernameRow
	userPayload, err := o.Store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return res, err
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return res, err
	}

	err = userService.ValidateUserRole(userPayload)
	if err != nil {
		logrus.Println("[UserService CreateUserService] createRoleUsers, ValidateUserRole : ", err)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
		return res, err
	}

	return userPayload, nil
}
