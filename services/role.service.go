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

// RoleService is
type RoleService struct {
}

var roleService *RoleService

// GetRoleService is
func GetRoleService() *RoleService {
	if roleService == nil {
		roleService = new(RoleService)
	}
	return roleService
}

func (o *RoleService) CreateRoleService(req dto.CreateRoleReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.RoleRes, error) {
	logrus.Println("[RoleService CreateRoleService] start.")
	var result dto.RoleRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.CreateRoleParams{
		Name:      req.Name,
		Level:     req.Level,
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	role, err := store.CreateRole(ctx, arg)
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

	result = RoleResponse(role)
	return result, nil
}

func (o *RoleService) GetRoleService(req dto.GetRoleReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.RoleRes, error) {
	logrus.Println("[RoleService GetRoleService] start.")
	var result dto.RoleRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	role, err := store.GetRole(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return result, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	result = RoleResponse(role)
	return result, nil
}

func (o *RoleService) ListRoleService(req dto.ListRoleRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) ([]dto.RoleRes, error) {
	logrus.Println("[RoleService ListRoleService] start.")
	var result []dto.RoleRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.ListRoleParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roles, err := store.ListRole(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	for _, role := range roles {
		u := RoleResponse(role)
		result = append(result, u)
	}

	return result, nil
}

func (o *RoleService) UpdateRoleService(req dto.UpdateRoleRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.RoleRes, error) {
	logrus.Println("[RoleService UpdateRoleService] start.")
	var result dto.RoleRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.UpdateRoleParams{
		ID:        req.ID,
		Name:      req.Name,
		Level:     req.Level,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	role, err := store.UpdateRole(ctx, arg)
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

	result = RoleResponse(role)

	return result, nil
}

func (o *RoleService) SoftDeleteRoleService(req dto.UpdateInactiveRoleRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) error {
	logrus.Println("[RoleService SoftDeleteRoleService] start.")
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveRoleParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = store.UpdateInactiveRole(ctx, arg)
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

func RoleResponse(role db.Role) dto.RoleRes {
	return dto.RoleRes{
		ID:    role.ID,
		Name:  role.Name,
		Level: role.Level,
	}
}
