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

// UserService is
type UserService struct {
}

var userService *UserService
var redisConn rmq.Connection

// GetUserService is
func GetUserService(config util.Config) *UserService {

	if userService == nil {

		userService = new(UserService)

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
	queueName := "test_123456"
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

func (o *UserService) CreateUserService(req dto.CreateUserRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.UserResponse, error) {
	logrus.Println("[UserService CreateUserService] start.")
	var result dto.UserResponse
	user, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		Username:       req.Username,
		CreatedBy:      sql.NullInt64{Int64: user.ID, Valid: true},
		Phone:          req.Phone,
		Balance:        sql.NullString{String: "0.00", Valid: true},
		IdentityNumber: req.IdentityNumber,
	}

	if req.Password != "" {
		password, errHash := util.HashPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
			return result, errHash
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	users, err := store.CreateUserTx(ctx, arg, authPayload, req.RoleID)
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

	return users, nil
}

func (o *UserService) CreateUserFirstService(req dto.CreateUserRequest, ctx *gin.Context, store db.Store) (dto.UserResponse, error) {
	logrus.Println("[UserService CreateUserFirstService] start.")
	var result dto.UserResponse
	if req.Username != "dbduabelas" {
		logrus.Println("[UserService CreateUserFirstService] error username not allow.")
		ctx.JSON(http.StatusNotFound, dto.ErrorResponseString("error username not allow"))
		return result, errors.New("error username not allow")
	}

	userPayload, err := store.GetUserByUsername(ctx, "dbduabelas")
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

	if req.Password != "" {
		password, errPswd := util.HashPassword(req.Password)
		if errPswd != nil {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(errPswd))
			return result, errPswd
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	users, err := store.CreateUser(ctx, arg)
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

	result = o.newUserResponse(users, []dto.RoleUser{})
	return result, nil
}

func (o *UserService) GetUserService(req dto.GetUserRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.UserResponse, error) {
	logrus.Println("[UserService GetUserService] start.")
	var result dto.UserResponse
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	user, err := store.GetUser(ctx, req.ID)
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
	roleUsers := getRoleByUser(argUser, ctx, store)
	userArg := db.User{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Username:       user.Username,
		Balance:        user.Balance,
		Phone:          user.Phone,
		IdentityNumber: user.IdentityNumber,
	}
	resp1 := o.newUserResponse(userArg, roleUsers)

	return resp1, nil
}

func (o *UserService) ListUserService(req dto.ListUserRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) ([]dto.UserResponse, error) {
	logrus.Println("[UserService ListUserService] start.")
	var result []dto.UserResponse
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.ListUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := store.ListUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	var resp1 []dto.UserResponse
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
		roleUserResponse := getRoleByUser(argUser, ctx, store)

		u := o.newUserResponse(dbUser, roleUserResponse)
		resp1 = append(resp1, u)
	}
	return resp1, nil
}

func (o *UserService) UpdateUserService(req dto.UpdateUserRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.UserResponse, error) {
	logrus.Println("[UserService UpdateUserService] start.")
	var result dto.UserResponse
	user, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.UpdateUserParams{
		ID:             req.ID,
		Email:          req.Email,
		Name:           req.Name,
		Phone:          req.Phone,
		IdentityNumber: req.IdentityNumber,
		UpdatedBy:      sql.NullInt64{Int64: user.ID, Valid: true},
	}

	if req.Password != "" {
		password, errPasswd := util.HashPassword(req.Password)
		if errPasswd != nil {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(errPasswd))
			return result, errPasswd
		}
		arg.Password = sql.NullString{String: password, Valid: true}
	}

	users, err := store.UpdateUserTx(ctx, arg, authPayload, req.RoleID)
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

	return users, nil
}

func (o *UserService) SoftDeleteUserService(req dto.UpdateInactiveUserRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) error {
	logrus.Println("[UserService SoftDeleteUserService] start.")
	user, err := validator(store, ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveUserParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: user.ID, Valid: true},
	}

	_, err = store.UpdateInactiveUser(ctx, arg)
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

func (o *UserService) LoginUserService(req dto.LoginUserRequest, tokenMaker token.Maker, ctx *gin.Context, store db.Store, config util.Config) (dto.LoginUserResponse, error) {
	logrus.Println("[UserService LoginUserService] start.")
	var result dto.LoginUserResponse
	user, err := store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		logrus.Println("[UserService loginUser] error get username is : ", err.Error())
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponseString("user not found"))
			return result, err
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	err = util.CheckPassword(req.Password, user.Password.String)
	if err != nil {
		logrus.Println("[UserService loginUser] password not same : ", err.Error())
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponseString("password is incorrect"))
		return result, err
	}

	accessToken, accessPayload, err := tokenMaker.CreateToken(
		user.Username,
		config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	refreshToken, refreshPayload, err := tokenMaker.CreateToken(
		user.Username,
		config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	argUser := db.User{
		ID: user.ID,
	}
	roleUsers := getRoleByUser(argUser, ctx, store)
	userRes := userRowToUserType(user)

	result = dto.LoginUserResponse{
		//SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  o.newUserResponse(userRes, roleUsers),
	}
	return result, nil
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

func (o *UserService) ValidateUserRole(user db.GetUserByUsernameRow) error {
	if user.RoleID.Int64 != 1 {
		return fmt.Errorf("you not allow to access this service")
	}

	return nil
}

func getRoleByUser(user db.User, ctx *gin.Context, store db.Store) []dto.RoleUser {
	roleUserArg := db.GetRoleUserByUserIDParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 0,
	}
	roleUsers, _ := store.GetRoleUserByUserID(ctx, roleUserArg)
	var roleUserResponse []dto.RoleUser
	for _, roleUser := range roleUsers {
		if roleUser.UserID == user.ID {
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

func validator(store db.Store, ctx *gin.Context, authPayload *token.Payload) (db.GetUserByUsernameRow, error) {
	logrus.Println("[UserService validator] start.")
	var res db.GetUserByUsernameRow
	userPayload, err := store.GetUserByUsername(ctx, authPayload.Username)
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
