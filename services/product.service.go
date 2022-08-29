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

type ProductInterface interface {
	CreateProductService(ctx *gin.Context, in dto.CreateProductReq) (dto.ProductRes, error)
	GetProductService(ctx *gin.Context, in dto.GetProductReq) (dto.ProductRes, error)
	ListProductService(ctx *gin.Context, in dto.ListProductRequest) ([]dto.ProductRes, error)
	UpdateProductService(ctx *gin.Context, in dto.UpdateProductRequest) (dto.ProductRes, error)
	SoftDeleteProductService(ctx *gin.Context, in dto.UpdateInactiveProductRequest) error
	setUpdateProd(arg db.UpdateProductParams, req dto.UpdateProductRequest) db.UpdateProductParams
	ProductRes(prod db.Product) dto.ProductRes
}

// ProductService is
type ProductService struct {
	store db.Store
}

//var productService *ProductService

// GetProductService is
func GetProductService(store db.Store) ProductInterface {
	return &ProductService{
		store: store,
	}
}

func (o *ProductService) CreateProductService(ctx *gin.Context, in dto.CreateProductReq) (dto.ProductRes, error) {
	logrus.Println("[ProductService CreateProductService] start.")
	var out dto.ProductRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.CreateProductParams{
		CatID:        in.CatID,
		Name:         in.Name,
		Amount:       in.Amount,
		ProviderID:   in.ProviderID,
		ProviderCode: in.ProviderCode,
		Status:       in.Status,
		Parent:       in.Parent,
		CreatedBy:    sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	if in.Status == "" {
		arg.Status = "active"
	}

	prod, err := o.store.CreateProduct(ctx, arg)
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

	out = o.ProductRes(prod)
	return out, nil
}

func (o *ProductService) GetProductService(ctx *gin.Context, in dto.GetProductReq) (dto.ProductRes, error) {
	logrus.Println("[ProductService GetProductService] start.")
	var out dto.ProductRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	prod, err := o.store.GetProduct(ctx, in.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.ProductRes(prod)
	return out, nil
}

func (o *ProductService) ListProductService(ctx *gin.Context, in dto.ListProductRequest) ([]dto.ProductRes, error) {
	logrus.Println("[ProductService ListProductService] start.")
	var out []dto.ProductRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.ListProductParams{
		Limit:      in.PageSize,
		Offset:     (in.PageID - 1) * in.PageSize,
		CatID:      in.CatID,
		ProviderID: in.ProviderID,
	}

	if in.CatID > 0 {
		arg.IsCat = true
	}

	if in.ProviderID > 0 {
		arg.IsProv = true
	}

	products, err := o.store.ListProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	for _, prod := range products {
		u := o.ProductRes(prod)
		out = append(out, u)
	}

	return out, nil
}

func (o *ProductService) UpdateProductService(ctx *gin.Context, in dto.UpdateProductRequest) (dto.ProductRes, error) {
	logrus.Println("[ProductService UpdateProductService] start.")
	var out dto.ProductRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	var arg = db.UpdateProductParams{
		ID:           in.ID,
		Name:         in.Name,
		CatID:        in.CatID,
		Amount:       in.Amount,
		ProviderID:   in.ProviderID,
		ProviderCode: in.ProviderCode,
		Status:       in.Status,
		Parent:       in.Parent,
		UpdatedBy:    sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg = o.setUpdateProd(arg, in)

	prod, err := o.store.UpdateProduct(ctx, arg)
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

	out = o.ProductRes(prod)

	return out, nil
}

func (o *ProductService) SoftDeleteProductService(ctx *gin.Context, in dto.UpdateInactiveProductRequest) error {
	logrus.Println("[ProductService SoftDeleteProviderService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveProductParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = o.store.UpdateInactiveProduct(ctx, arg)
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

func (o *ProductService) setUpdateProd(arg db.UpdateProductParams, in dto.UpdateProductRequest) db.UpdateProductParams {

	if in.Name != "" {
		arg.SetName = true
	}
	if in.CatID > 0 {
		arg.SetCat = true
	}
	if in.Amount != "" {
		arg.SetAmount = true
	}
	if in.ProviderID > 0 {
		arg.SetProvider = true
	}
	if in.ProviderCode != "" {
		arg.SetProviderCode = true
	}
	if in.Status != "" {
		arg.SetStatus = true
	}
	if in.Parent > 0 {
		arg.SetParent = true
	}

	return arg
}

func (o *ProductService) ProductRes(prod db.Product) dto.ProductRes {
	return dto.ProductRes{
		ID:           prod.ID,
		Name:         prod.Name,
		CatID:        prod.CatID,
		Amount:       prod.Amount,
		ProviderID:   prod.ProviderID,
		ProviderCode: prod.ProviderCode,
		Status:       prod.Status,
		Parent:       prod.Parent,
	}
}
