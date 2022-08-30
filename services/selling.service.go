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

type SellingInterface interface {
	CreateSellingService(ctx *gin.Context, in dto.CreateSellingReq) (dto.SellingRes, error)
	GetSellingService(ctx *gin.Context, in dto.GetSellingReq) (dto.SellingRes, error)
	ListSellingService(ctx *gin.Context, in dto.ListSellingRequest) ([]dto.SellingRes, error)
	UpdateSellingService(ctx *gin.Context, in dto.UpdateSellingRequest) (dto.SellingRes, error)
	SoftDeleteSellingService(ctx *gin.Context, in dto.UpdateInactiveSellingRequest) error
	SellingRes(sell db.Selling) dto.SellingRes
}

// SellingService is
type SellingService struct {
	store db.Store
}

var sellingService *SellingService

// GetSellingService is
func GetSellingService(store db.Store) *SellingService {
	if sellingService == nil {
		sellingService = &SellingService{
			store: store,
		}
	}
	return sellingService
}

func (o *SellingService) CreateSellingService(ctx *gin.Context, in dto.CreateSellingReq) (dto.SellingRes, error) {
	logrus.Println("[SellingService CreateSellingService] start.")
	var out dto.SellingRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.CreateSellingParams{
		PartnerID: sql.NullInt64{
			Int64: in.PartnerID,
			Valid: true,
		},
		CategoryID: sql.NullInt64{
			Int64: in.CategoryID,
			Valid: true,
		},
		Amount: sql.NullString{
			String: in.Amount,
			Valid:  true,
		},
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	sell, err := o.store.CreateSelling(ctx, arg)
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

	out = o.SellingRes(sell)
	return out, nil
}

func (o *SellingService) GetSellingService(ctx *gin.Context, in dto.GetSellingReq) (dto.SellingRes, error) {
	logrus.Println("[SellingService GetSellingService] start.")
	var out dto.SellingRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	prod, err := o.store.GetSelling(ctx, in.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.SellingRes(prod)
	return out, nil
}

func (o *SellingService) ListSellingService(ctx *gin.Context, in dto.ListSellingRequest) ([]dto.SellingRes, error) {
	logrus.Println("[SellingService ListSellingService] start.")
	var out []dto.SellingRes
	out = []dto.SellingRes{}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.ListSellingByParamsParams{
		Limit:     in.PageSize,
		Offset:    (in.PageID - 1) * in.PageSize,
		IsPartner: false,
		PartnerID: sql.NullInt64{
			Int64: in.PartnerID,
			Valid: true,
		},
		IsCategory: false,
		CategoryID: sql.NullInt64{
			Int64: in.CategoryID,
			Valid: true,
		},
	}

	if in.PartnerID > 0 {
		arg.IsPartner = true
	}

	if in.CategoryID > 0 {
		arg.IsCategory = true
	}

	sells, err := o.store.ListSellingByParams(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	logrus.Println(sells)

	for _, selling := range sells {
		u := o.SellingRes(selling)
		out = append(out, u)
	}

	return out, nil
}

func (o *SellingService) UpdateSellingService(ctx *gin.Context, in dto.UpdateSellingRequest) (dto.SellingRes, error) {
	logrus.Println("[SellingService UpdateSellingService] start.")
	var out dto.SellingRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	var arg = db.UpdateSellingParams{
		SetPartnerID: true,
		PartnerID: sql.NullInt64{
			Int64: in.PartnerID,
			Valid: true,
		},
		SetCategoryID: true,
		CategoryID: sql.NullInt64{
			Int64: in.CategoryID,
			Valid: true,
		},
		SetAmount: true,
		Amount: sql.NullString{
			String: in.Amount,
			Valid:  true,
		},
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
		ID:        in.ID,
	}

	prod, err := o.store.UpdateSelling(ctx, arg)
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

	out = o.SellingRes(prod)

	return out, nil
}

func (o *SellingService) SoftDeleteSellingService(ctx *gin.Context, in dto.UpdateInactiveSellingRequest) error {
	logrus.Println("[SellingService SoftDeleteSellingService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveSellingParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = o.store.UpdateInactiveSelling(ctx, arg)
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
