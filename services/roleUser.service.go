package services

import (
	"database/sql"
	"errors"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/dto"
	"github.com/ariandi/ppob_go/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

type RoleUserInterface interface {
	CreateRoleUserService(ctx *gin.Context, in dto.CreateRoleUserReq) (dto.RoleUserRes, error)
	GetRoleUserByUserIDService(ctx *gin.Context, in dto.GetRoleUserByUserID) ([]dto.RoleUserRes, error)
	ListRoleUserService(ctx *gin.Context, in dto.ListRoleUserRequest) ([]dto.RoleUserRes, error)
	UpdateRoleUserService(ctx *gin.Context, in dto.UpdateRoleUserRequest) (dto.RoleUserRes, error)
	SoftDeleteRoleUserService(ctx *gin.Context, in dto.UpdateInactiveROleUserRequest) error
	roleUserResponse(roleUser db.RoleUser) dto.RoleUserRes
}

type RoleUserService struct {
	store db.Store
}

// GetRoleUserService is
func GetRoleUserService(store db.Store) RoleUserInterface {
	return &RoleUserService{
		store: store,
	}
}

func (o *RoleUserService) CreateRoleUserService(ctx *gin.Context, in dto.CreateRoleUserReq) (dto.RoleUserRes, error) {
	logrus.Println("[RoleUserService CreateRoleUserService] start.")
	var out dto.RoleUserRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.CreateRoleUserParams{
		RoleID:    in.RoleID,
		UserID:    in.UserID,
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	cat, err := o.store.CreateRoleUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return out, err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.roleUserResponse(cat)
	return out, nil
}

func (o *RoleUserService) GetRoleUserByUserIDService(ctx *gin.Context, in dto.GetRoleUserByUserID) ([]dto.RoleUserRes, error) {
	logrus.Println("[RoleUserService GetRoleUserService] start.")
	var out []dto.RoleUserRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.GetRoleUserByUserIDParams{
		UserID: in.UserID,
		Limit:  5,
		Offset: 0,
	}
	getRolesByUserID, err := o.store.GetRoleUserByUserID(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	for _, roleUser := range getRolesByUserID {
		resArg := o.roleUserResponse(roleUser)
		out = append(out, resArg)
	}

	return out, nil
}

func (o *RoleUserService) ListRoleUserService(ctx *gin.Context, in dto.ListRoleUserRequest) ([]dto.RoleUserRes, error) {
	logrus.Println("[RoleUserService GetRoleUserService] start.")
	var out []dto.RoleUserRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.ListRoleUserParams{
		Limit:  in.PageSize,
		Offset: (in.PageID - 1) * in.PageSize,
	}

	roleUsers, err := o.store.ListRoleUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	for _, roleUser := range roleUsers {
		u := o.roleUserResponse(roleUser)
		out = append(out, u)
	}

	return out, nil
}

func (o *RoleUserService) UpdateRoleUserService(ctx *gin.Context, in dto.UpdateRoleUserRequest) (dto.RoleUserRes, error) {
	logrus.Println("[RoleUserService UpdateRoleUserService] start.")
	var out dto.RoleUserRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.UpdateRoleUserParams{
		ID:        in.ID,
		UserID:    in.UserID,
		RoleID:    in.RoleID,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	cat, err := o.store.UpdateRoleUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, dto.ErrorResponse(err))
				return out, err
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.roleUserResponse(cat)

	return out, nil
}

func (o *RoleUserService) SoftDeleteRoleUserService(ctx *gin.Context, in dto.UpdateInactiveROleUserRequest) error {
	logrus.Println("[RoleUserService SoftDeleteRoleUserService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveRoleUserParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = o.store.UpdateInactiveRoleUser(ctx, arg)
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

func (o *RoleUserService) roleUserResponse(roleUser db.RoleUser) dto.RoleUserRes {
	return dto.RoleUserRes{
		ID:     roleUser.ID,
		UserID: roleUser.UserID,
		RoleID: roleUser.RoleID,
	}
}
