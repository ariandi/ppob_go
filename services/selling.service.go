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

// SellingService is
type SellingService struct {
}

var sellingService *SellingService

// GetSellingService is
func GetSellingService() *SellingService {
	if sellingService == nil {
		sellingService = new(SellingService)
	}
	return sellingService
}

func (o *SellingService) CreateSellingService(req dto.CreateSellingReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.SellingRes, error) {
	logrus.Println("[SellingService CreateSellingService] start.")
	var result dto.SellingRes

	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.CreateSellingParams{
		PartnerID: sql.NullInt64{
			Int64: req.PartnerID,
			Valid: true,
		},
		CategoryID: sql.NullInt64{
			Int64: req.CategoryID,
			Valid: true,
		},
		Amount: sql.NullString{
			String: req.Amount,
			Valid:  true,
		},
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	sell, err := store.CreateSelling(ctx, arg)
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

	result = o.SellingRes(sell)
	return result, nil
}

func (o *SellingService) GetSellingService(req dto.GetSellingReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.SellingRes, error) {
	logrus.Println("[SellingService GetSellingService] start.")
	var result dto.SellingRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	prod, err := store.GetSelling(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return result, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	result = o.SellingRes(prod)
	return result, nil
}

func (o *SellingService) ListSellingService(req dto.ListSellingRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) ([]dto.SellingRes, error) {
	logrus.Println("[SellingService ListSellingService] start.")
	var result []dto.SellingRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.ListSellingByParamsParams{
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
		IsPartner:  false,
		PartnerID:  sql.NullInt64{},
		IsCategory: false,
		CategoryID: sql.NullInt64{},
	}

	if req.PartnerID > 0 {
		arg.IsPartner = true
	}

	if req.CategoryID > 0 {
		arg.IsCategory = true
	}

	sellings, err := store.ListSellingByParams(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	for _, selling := range sellings {
		u := o.SellingRes(selling)
		result = append(result, u)
	}

	return result, nil
}

func (o *SellingService) UpdateSellingService(req dto.UpdateSellingRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.SellingRes, error) {
	logrus.Println("[SellingService UpdateSellingService] start.")
	var result dto.SellingRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	var arg = db.UpdateSellingParams{
		SetPartnerID: true,
		PartnerID: sql.NullInt64{
			Int64: req.PartnerID,
			Valid: true,
		},
		SetCategoryID: true,
		CategoryID: sql.NullInt64{
			Int64: req.CategoryID,
			Valid: true,
		},
		SetAmount: true,
		Amount: sql.NullString{
			String: req.Amount,
			Valid:  true,
		},
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
		ID:        req.ID,
	}

	prod, err := store.UpdateSelling(ctx, arg)
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

	result = o.SellingRes(prod)

	return result, nil
}

func (o *SellingService) SoftDeleteSellingService(req dto.UpdateInactiveSellingRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) error {
	logrus.Println("[SellingService SoftDeleteSellingService] start.")
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveSellingParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = store.UpdateInactiveSelling(ctx, arg)
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

func (o *SellingService) SellingRes(sell db.Selling) dto.SellingRes {
	return dto.SellingRes{
		ID:         sell.ID,
		PartnerID:  sell.PartnerID.Int64,
		CategoryID: sell.CategoryID.Int64,
		Amount:     sell.Amount.String,
	}
}
