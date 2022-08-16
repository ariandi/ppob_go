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

// ProductService is
type ProductService struct {
}

var productService *ProductService

// GetProductService is
func GetProductService() *ProductService {
	if productService == nil {
		productService = new(ProductService)
	}
	return productService
}

func (o *ProductService) CreateProductService(req dto.CreateProductReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.ProductRes, error) {
	logrus.Println("[ProductService CreateProductService] start.")
	var result dto.ProductRes

	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.CreateProductParams{
		CatID:      0,
		Name:       req.Name,
		Amount:     req.Amount,
		ProviderID: req.ProviderID,
		Status:     req.Status,
		Parent:     req.Parent,
		CreatedBy:  sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	if req.Status == "" {
		arg.Status = "active"
	}

	prod, err := store.CreateProduct(ctx, arg)
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

	result = o.ProductRes(prod)
	return result, nil
}

func (o *ProductService) GetProductService(req dto.GetProductReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.ProductRes, error) {
	logrus.Println("[ProductService GetProductService] start.")
	var result dto.ProductRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	prod, err := store.GetProduct(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return result, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	result = o.ProductRes(prod)
	return result, nil
}

func (o *ProductService) ListProductService(req dto.ListProductRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) ([]dto.ProductRes, error) {
	logrus.Println("[ProductService ListProductService] start.")
	var result []dto.ProductRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.ListProductParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	products, err := store.ListProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	for _, prod := range products {
		u := o.ProductRes(prod)
		result = append(result, u)
	}

	return result, nil
}

func (o *ProductService) UpdateProductService(req dto.UpdateProductRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.ProductRes, error) {
	logrus.Println("[ProductService UpdateProductService] start.")
	var result dto.ProductRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	var arg = db.UpdateProductParams{
		ID:         req.ID,
		Name:       req.Name,
		CatID:      req.CatID,
		Amount:     req.Amount,
		ProviderID: req.ProviderID,
		Status:     req.Status,
		Parent:     req.Parent,
		UpdatedBy:  sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg = o.setUpdateProd(arg, req)

	prod, err := store.UpdateProduct(ctx, arg)
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

	result = o.ProductRes(prod)

	return result, nil
}

func (o *ProductService) SoftDeleteProductService(req dto.UpdateInactiveProductRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) error {
	logrus.Println("[ProductService SoftDeleteProviderService] start.")
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveProductParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = store.UpdateInactiveProduct(ctx, arg)
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

func (o ProductService) setUpdateProd(arg db.UpdateProductParams, req dto.UpdateProductRequest) db.UpdateProductParams {

	if req.Name != "" {
		arg.SetName = true
	}
	if req.CatID > 0 {
		arg.SetCat = true
	}
	if req.Amount != "" {
		arg.SetAmount = true
	}
	if req.ProviderID > 0 {
		arg.SetProvider = true
	}
	if req.Status != "" {
		arg.SetStatus = true
	}
	if req.Parent > 0 {
		arg.SetParent = true
	}

	return arg
}

func (o ProductService) ProductRes(prod db.Product) dto.ProductRes {
	return dto.ProductRes{
		ID:         prod.ID,
		Name:       prod.Name,
		CatID:      prod.CatID,
		Amount:     prod.Amount,
		ProviderID: prod.ProviderID,
		Status:     prod.Status,
		Parent:     prod.Parent,
	}
}
