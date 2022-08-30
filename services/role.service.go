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

type RoleInterface interface {
	CreateRoleService(ctx *gin.Context, in dto.CreateRoleReq) (dto.RoleRes, error)
	GetRoleService(ctx *gin.Context, in dto.GetRoleReq) (dto.RoleRes, error)
	ListRoleService(ctx *gin.Context, in dto.ListRoleRequest) ([]dto.RoleRes, error)
	UpdateRoleService(ctx *gin.Context, in dto.UpdateRoleRequest) (dto.RoleRes, error)
	SoftDeleteRoleService(ctx *gin.Context, in dto.UpdateInactiveRoleRequest) error
	RoleResponse(role db.Role) dto.RoleRes
}

// RoleService is
type RoleService struct {
	store db.Store
}

var roleService *RoleService

// GetRoleService is
func GetRoleService(store db.Store) RoleInterface {
	if roleService == nil {
		roleService = &RoleService{
			store: store,
		}
	}
	return roleService
}

func (o *RoleService) CreateRoleService(ctx *gin.Context, in dto.CreateRoleReq) (dto.RoleRes, error) {
	logrus.Println("[RoleService CreateRoleService] start.")
	var out dto.RoleRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.CreateRoleParams{
		Name:      in.Name,
		Level:     in.Level,
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	role, err := o.store.CreateRole(ctx, arg)
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

	out = o.RoleResponse(role)
	return out, nil
}

func (o *RoleService) GetRoleService(ctx *gin.Context, in dto.GetRoleReq) (dto.RoleRes, error) {
	logrus.Println("[RoleService GetRoleService] start.")
	var out dto.RoleRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	role, err := o.store.GetRole(ctx, in.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.RoleResponse(role)
	return out, nil
}

func (o *RoleService) ListRoleService(ctx *gin.Context, in dto.ListRoleRequest) ([]dto.RoleRes, error) {
	logrus.Println("[RoleService ListRoleService] start.")
	var out []dto.RoleRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.ListRoleParams{
		Limit:  in.PageSize,
		Offset: (in.PageID - 1) * in.PageSize,
	}

	roles, err := o.store.ListRole(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	for _, role := range roles {
		u := o.RoleResponse(role)
		out = append(out, u)
	}

	return out, nil
}

func (o *RoleService) UpdateRoleService(ctx *gin.Context, in dto.UpdateRoleRequest) (dto.RoleRes, error) {
	logrus.Println("[RoleService UpdateRoleService] start.")
	var out dto.RoleRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.UpdateRoleParams{
		ID:        in.ID,
		Name:      in.Name,
		Level:     in.Level,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	role, err := o.store.UpdateRole(ctx, arg)
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

	out = o.RoleResponse(role)

	return out, nil
}

func (o *RoleService) SoftDeleteRoleService(ctx *gin.Context, in dto.UpdateInactiveRoleRequest) error {
	logrus.Println("[RoleService SoftDeleteRoleService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveRoleParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = o.store.UpdateInactiveRole(ctx, arg)
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

func (o *RoleService) RoleResponse(role db.Role) dto.RoleRes {
	return dto.RoleRes{
		ID:    role.ID,
		Name:  role.Name,
		Level: role.Level,
	}
}
