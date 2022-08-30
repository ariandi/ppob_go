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

type PartnerInterface interface {
	CreatePartnerService(ctx *gin.Context, in dto.CreatePartnerReq) (dto.PartnerRes, error)
	GetPartnerService(ctx *gin.Context, in dto.GetPartnerReq) (dto.PartnerRes, error)
	ListPartnerService(ctx *gin.Context, in dto.ListPartnerRequest) ([]dto.PartnerRes, error)
	UpdatePartnerService(ctx *gin.Context, in dto.UpdatePartnerRequest) (dto.PartnerRes, error)
	SoftDeletePartnerService(ctx *gin.Context, in dto.UpdateInactivePartnerRequest) error
	PartnerResponse(partner db.Partner) dto.PartnerRes
	setUpdateParamsService(arg db.UpdatePartnerParams, req dto.UpdatePartnerRequest) db.UpdatePartnerParams
	setFromDateToDate(arg db.CreatePartnerParams, req dto.CreatePartnerReq) (db.CreatePartnerParams, error)
}

// PartnerService is
type PartnerService struct {
	store db.Store
}

var partnerService *PartnerService

// GetPartnerService is
func GetPartnerService(store db.Store) PartnerInterface {
	if partnerService == nil {
		partnerService = &PartnerService{
			store: store,
		}
	}
	return partnerService
}

func (o *PartnerService) CreatePartnerService(ctx *gin.Context, in dto.CreatePartnerReq) (dto.PartnerRes, error) {
	logrus.Println("[PartnerService CreatePartnerService] start.")
	var out dto.PartnerRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.CreatePartnerParams{
		Name:        in.Name,
		User:        in.User,
		Secret:      in.Secret,
		AddInfo1:    in.AddInfo1,
		AddInfo2:    in.AddInfo2,
		PaymentType: in.PaymentType,
		Status:      in.Status,
		CreatedBy:   sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg, err = o.setFromDateToDate(arg, in)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return out, err
	}

	partner, err := o.store.CreatePartner(ctx, arg)
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

	out = o.PartnerResponse(partner)
	return out, nil
}

func (o *PartnerService) GetPartnerService(ctx *gin.Context, in dto.GetPartnerReq) (dto.PartnerRes, error) {
	logrus.Println("[CategoryService GetCategoryService] start.")
	var out dto.PartnerRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	partner, err := o.store.GetPartner(ctx, in.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.PartnerResponse(partner)
	return out, nil
}

func (o *PartnerService) ListPartnerService(ctx *gin.Context, in dto.ListPartnerRequest) ([]dto.PartnerRes, error) {
	logrus.Println("[CategoryService GetCategoryService] start.")
	var out []dto.PartnerRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.ListPartnerParams{
		Limit:  in.PageSize,
		Offset: (in.PageID - 1) * in.PageSize,
	}

	partners, err := o.store.ListPartner(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	for _, partner := range partners {
		u := o.PartnerResponse(partner)
		out = append(out, u)
	}

	return out, nil
}

func (o *PartnerService) UpdatePartnerService(ctx *gin.Context, in dto.UpdatePartnerRequest) (dto.PartnerRes, error) {
	logrus.Println("[PartnerService UpdatePartnerService] start.")
	var out dto.PartnerRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.UpdatePartnerParams{
		ID:        in.ID,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg = o.setUpdateParamsService(arg, in)

	partner, err := o.store.UpdatePartner(ctx, arg)
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

	out = o.PartnerResponse(partner)

	return out, nil
}

func (o *PartnerService) SoftDeletePartnerService(ctx *gin.Context, in dto.UpdateInactivePartnerRequest) error {
	logrus.Println("[PartnerService SoftDeletePartnerService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactivePartnerParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = o.store.UpdateInactivePartner(ctx, arg)
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
