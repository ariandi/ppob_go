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

// CategoryService is
type CategoryService struct {
}

var categoryService *CategoryService

// GetCategoryService is
func GetCategoryService() *CategoryService {
	if categoryService == nil {
		categoryService = new(CategoryService)
	}
	return categoryService
}

func (o *CategoryService) CreateCategoryService(req dto.CreateCategoryReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.CategoryRes, error) {
	logrus.Println("[CategoryService CreateCategoryService] start.")
	var result dto.CategoryRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.CreateCategoryParams{
		Name:      req.Name,
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	cat, err := store.CreateCategory(ctx, arg)
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

	result = o.CatResponse(cat)
	return result, nil
}

func (o *CategoryService) GetCategoryService(req dto.GetCategoryReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.CategoryRes, error) {
	logrus.Println("[CategoryService GetCategoryService] start.")
	var result dto.CategoryRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	cat, err := store.GetCategory(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return result, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	result = o.CatResponse(cat)
	return result, nil
}

func (o *CategoryService) ListCategoryService(req dto.ListCategoryRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) ([]dto.CategoryRes, error) {
	logrus.Println("[CategoryService GetCategoryService] start.")
	var result []dto.CategoryRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.ListCategoryParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	categories, err := store.ListCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	for _, cat := range categories {
		u := o.CatResponse(cat)
		result = append(result, u)
	}

	return result, nil
}

func (o *CategoryService) UpdateCategoryService(req dto.UpdateCategoryRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.CategoryRes, error) {
	logrus.Println("[CategoryService UpdateCategoryService] start.")
	var result dto.CategoryRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.UpdateCategoryParams{
		ID:        req.ID,
		SetName:   true,
		Name:      req.Name,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	cat, err := store.UpdateCategory(ctx, arg)
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

	result = o.CatResponse(cat)

	return result, nil
}

func (o *CategoryService) SoftDeleteCategoryService(req dto.UpdateInactiveCategoryRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) error {
	logrus.Println("[CategoryService SoftDeleteCategoryService] start.")
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveCategoryParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = store.UpdateInactiveCategory(ctx, arg)
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
		ID:   cat.ID,
		Name: cat.Name,
	}
}
