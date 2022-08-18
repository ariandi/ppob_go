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
	"time"
)

// PartnerService is
type PartnerService struct {
}

var partnerService *PartnerService

// GetPartnerService is
func GetPartnerService() *PartnerService {
	if partnerService == nil {
		partnerService = new(PartnerService)
	}
	return partnerService
}

func (o *PartnerService) CreatePartnerService(req dto.CreatePartnerReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.PartnerRes, error) {
	logrus.Println("[PartnerService CreatePartnerService] start.")
	var result dto.PartnerRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.CreatePartnerParams{
		Name:        req.Name,
		User:        req.User,
		Secret:      req.Secret,
		AddInfo1:    req.AddInfo1,
		AddInfo2:    req.AddInfo2,
		PaymentType: req.PaymentType,
		Status:      req.Status,
		CreatedBy:   sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg, err = o.setFromDateToDate(arg, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return result, err
	}

	partner, err := store.CreatePartner(ctx, arg)
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

	result = o.PartnerResponse(partner)
	return result, nil
}

func (o *PartnerService) GetPartnerService(req dto.GetPartnerReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.PartnerRes, error) {
	logrus.Println("[CategoryService GetCategoryService] start.")
	var result dto.PartnerRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	partner, err := store.GetPartner(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return result, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	result = o.PartnerResponse(partner)
	return result, nil
}

func (o *PartnerService) ListPartnerService(req dto.ListPartnerRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) ([]dto.PartnerRes, error) {
	logrus.Println("[CategoryService GetCategoryService] start.")
	var result []dto.PartnerRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.ListPartnerParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	partners, err := store.ListPartner(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	for _, partner := range partners {
		u := o.PartnerResponse(partner)
		result = append(result, u)
	}

	return result, nil
}

func (o *PartnerService) UpdatePartnerService(req dto.UpdatePartnerRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.PartnerRes, error) {
	logrus.Println("[PartnerService UpdatePartnerService] start.")
	var result dto.PartnerRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.UpdatePartnerParams{
		ID:        req.ID,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg = o.setUpdateParamsService(arg, req)

	partner, err := store.UpdatePartner(ctx, arg)
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

	result = o.PartnerResponse(partner)

	return result, nil
}

func (o *PartnerService) SoftDeletePartnerService(req dto.UpdateInactivePartnerRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) error {
	logrus.Println("[PartnerService SoftDeletePartnerService] start.")
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactivePartnerParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = store.UpdateInactivePartner(ctx, arg)
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

func (o *PartnerService) PartnerResponse(partner db.Partner) dto.PartnerRes {
	return dto.PartnerRes{
		ID:          partner.ID,
		Name:        partner.Name,
		User:        partner.User,
		Secret:      partner.Secret,
		AddInfo1:    partner.AddInfo2,
		AddInfo2:    partner.AddInfo2,
		ValidFrom:   partner.ValidFrom.Time.Format("2006-01-02 15:04:05"),
		ValidTo:     partner.ValidTo.Time.Format("2006-01-02 15:04:05"),
		PaymentType: partner.PaymentType,
		Status:      partner.Status,
	}
}

func (o *PartnerService) setUpdateParamsService(arg db.UpdatePartnerParams, req dto.UpdatePartnerRequest) db.UpdatePartnerParams {
	if req.Name != "" {
		arg.Name = req.Name
		arg.SetName = true
	}

	if req.User != "" {
		arg.UserParams = req.User
		arg.SetUser = true
	}

	if req.Secret != "" {
		arg.Secret = req.Secret
		arg.SetSecret = true
	}

	if req.AddInfo1 != "" {
		arg.AddInfo1 = req.AddInfo1
		arg.SetAddInfo1 = true
	}

	if req.AddInfo2 != "" {
		arg.AddInfo2 = req.AddInfo2
		arg.SetAddInfo2 = true
	}

	if req.PaymentType != "" {
		arg.PaymentType = req.PaymentType
		arg.SetPaymentType = true
	}

	return arg
}

func (o *PartnerService) setFromDateToDate(arg db.CreatePartnerParams, req dto.CreatePartnerReq) (db.CreatePartnerParams, error) {
	layoutFormat := "2006-01-02 15:04:05"

	if req.ValidFrom != "" {
		fromDate, err := time.Parse(layoutFormat, req.ValidFrom)
		if err != nil {
			return arg, errors.New("from date is not format")
		}

		arg.ValidFrom = sql.NullTime{
			Time:  fromDate,
			Valid: true,
		}
	}

	if req.ValidTo != "" {
		toDate, err := time.Parse("2006-01-02 15:04:05", req.ValidTo)
		if err != nil {
			return arg, errors.New("to date is not format")
		}

		arg.ValidTo = sql.NullTime{
			Time:  toDate,
			Valid: true,
		}
	}

	if req.Status == "" {
		arg.Status = "active"
	}

	return arg, nil
}
