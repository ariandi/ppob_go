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

type CategoryInterface interface {
	CreateCategoryService(ctx *gin.Context, in dto.CreateCategoryReq) (dto.CategoryRes, error)
	GetCategoryService(ctx *gin.Context, in dto.GetCategoryReq) (dto.CategoryRes, error)
	ListCategoryService(ctx *gin.Context, in dto.ListCategoryRequest) ([]dto.CategoryRes, error)
	UpdateCategoryService(ctx *gin.Context, in dto.UpdateCategoryRequest) (dto.CategoryRes, error)
	SoftDeleteCategoryService(ctx *gin.Context, in dto.UpdateInactiveCategoryRequest) error
	CatResponse(cat db.Category) dto.CategoryRes
}

// CategoryService is
type CategoryService struct {
	store db.Store
}

var categoryService *CategoryService

// GetCategoryService is
func GetCategoryService(store db.Store) CategoryInterface {
	if categoryService == nil {
		categoryService = &CategoryService{
			store: store,
		}
	}
	return categoryService
}

func (o *CategoryService) CreateCategoryService(ctx *gin.Context, in dto.CreateCategoryReq) (dto.CategoryRes, error) {
	logrus.Println("[CategoryService CreateCategoryService] start.")
	var out dto.CategoryRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.CreateCategoryParams{
		Name: in.Name,
		UpSelling: sql.NullString{
			String: in.UpSelling,
			Valid:  true,
		},
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	cat, err := o.store.CreateCategory(ctx, arg)
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

	out = o.CatResponse(cat)
	return out, nil
}

func (o *CategoryService) GetCategoryService(ctx *gin.Context, in dto.GetCategoryReq) (dto.CategoryRes, error) {
	logrus.Println("[CategoryService GetCategoryService] start.")
	var out dto.CategoryRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	cat, err := o.store.GetCategory(ctx, in.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.CatResponse(cat)
	return out, nil
}

func (o *CategoryService) ListCategoryService(ctx *gin.Context, in dto.ListCategoryRequest) ([]dto.CategoryRes, error) {
	logrus.Println("[CategoryService GetCategoryService] start.")
	var out []dto.CategoryRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.ListCategoryParams{
		Limit:  in.PageSize,
		Offset: (in.PageID - 1) * in.PageSize,
	}

	categories, err := o.store.ListCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	for _, cat := range categories {
		u := o.CatResponse(cat)
		out = append(out, u)
	}

	return out, nil
}

func (o *CategoryService) UpdateCategoryService(ctx *gin.Context, in dto.UpdateCategoryRequest) (dto.CategoryRes, error) {
	logrus.Println("[CategoryService UpdateCategoryService] start.")
	var out dto.CategoryRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.UpdateCategoryParams{
		ID:           in.ID,
		SetName:      true,
		Name:         in.Name,
		SetUpSelling: true,
		UpSelling: sql.NullString{
			String: in.UpSelling,
			Valid:  true,
		},
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	cat, err := o.store.UpdateCategory(ctx, arg)
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

	out = o.CatResponse(cat)

	return out, nil
}

func (o *CategoryService) SoftDeleteCategoryService(ctx *gin.Context, in dto.UpdateInactiveCategoryRequest) error {
	logrus.Println("[CategoryService SoftDeleteCategoryService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveCategoryParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = o.store.UpdateInactiveCategory(ctx, arg)
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

func (o *CategoryService) CatResponse(cat db.Category) dto.CategoryRes {
	return dto.CategoryRes{
		ID:        cat.ID,
		Name:      cat.Name,
		UpSelling: cat.UpSelling.String,
	}
}
