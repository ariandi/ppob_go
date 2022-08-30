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

type ProviderInterface interface {
	CreateProviderService(ctx *gin.Context, in dto.CreateProviderReq) (dto.ProviderRes, error)
	GetProviderService(ctx *gin.Context, in dto.GetProviderReq) (dto.ProviderRes, error)
	ListProviderService(ctx *gin.Context, in dto.ListProviderRequest) ([]dto.ProviderRes, error)
	UpdateProviderService(ctx *gin.Context, in dto.UpdateProviderRequest) (dto.ProviderRes, error)
	SoftDeleteProviderService(ctx *gin.Context, in dto.UpdateInactiveProviderRequest) error
	setCreateProvider(arg db.CreateProviderParams, in dto.CreateProviderReq) (db.CreateProviderParams, error)
	setUpdateProvider(arg db.UpdateProviderParams, in dto.UpdateProviderRequest) (db.UpdateProviderParams, error)
	ProviderRes(provider db.Provider) dto.ProviderRes
}

// ProviderService is
type ProviderService struct {
	store db.Store
}

var providerService *ProviderService

// GetProviderService is
func GetProviderService(store db.Store) *ProviderService {
	if providerService == nil {
		providerService = &ProviderService{
			store: store,
		}
	}
	return providerService
}

func (o *ProviderService) CreateProviderService(ctx *gin.Context, in dto.CreateProviderReq) (dto.ProviderRes, error) {
	logrus.Println("[ProviderService CreateProviderService] start.")
	var out dto.ProviderRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.CreateProviderParams{
		Name:      in.Name,
		User:      in.User,
		Secret:    in.Secret,
		AddInfo1:  in.AddInfo1,
		AddInfo2:  in.AddInfo2,
		BaseUrl:   sql.NullString{},
		Method:    sql.NullString{},
		Inq:       sql.NullString{},
		Pay:       sql.NullString{},
		Adv:       sql.NullString{},
		Cmt:       sql.NullString{},
		Rev:       sql.NullString{},
		Status:    in.Status,
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg, err = o.setCreateProvider(arg, in)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return out, err
	}

	provider, err := o.store.CreateProvider(ctx, arg)
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

	out = o.ProviderRes(provider)
	return out, nil
}

func (o *ProviderService) GetProviderService(ctx *gin.Context, in dto.GetProviderReq) (dto.ProviderRes, error) {
	logrus.Println("[ProviderService GetProviderService] start.")
	var out dto.ProviderRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	prov, err := o.store.GetProvider(ctx, in.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.ProviderRes(prov)
	return out, nil
}

func (o *ProviderService) ListProviderService(ctx *gin.Context, in dto.ListProviderRequest) ([]dto.ProviderRes, error) {
	logrus.Println("[ProviderService ListProviderService] start.")
	var out []dto.ProviderRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.ListProviderParams{
		Limit:  in.PageSize,
		Offset: (in.PageID - 1) * in.PageSize,
	}

	providers, err := o.store.ListProvider(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	for _, provider := range providers {
		u := o.ProviderRes(provider)
		out = append(out, u)
	}

	return out, nil
}

func (o *ProviderService) UpdateProviderService(ctx *gin.Context, in dto.UpdateProviderRequest) (dto.ProviderRes, error) {
	logrus.Println("[ProviderService UpdateProviderService] start.")
	var out dto.ProviderRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.UpdateProviderParams{
		ID:         in.ID,
		Name:       in.Name,
		UserParams: in.User,
		Secret:     in.Secret,
		AddInfo1:   in.AddInfo1,
		AddInfo2:   in.AddInfo2,
		UpdatedBy:  sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg, err = o.setUpdateProvider(arg, in)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return out, err
	}

	prov, err := o.store.UpdateProvider(ctx, arg)
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

	out = o.ProviderRes(prov)

	return out, nil
}

func (o *ProviderService) SoftDeleteProviderService(ctx *gin.Context, in dto.UpdateInactiveProviderRequest) error {
	logrus.Println("[ProviderService SoftDeleteProviderService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveProviderParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = o.store.UpdateInactiveProvider(ctx, arg)
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

func (o *ProviderService) setCreateProvider(arg db.CreateProviderParams, in dto.CreateProviderReq) (db.CreateProviderParams, error) {
	layoutFormat := "2006-01-02 15:04:05"
	if in.ValidFrom != "" {
		fromDate, err := time.Parse(layoutFormat, in.ValidFrom)
		if err != nil {
			return arg, errors.New("from date is not format")
		}

		arg.ValidFrom = sql.NullTime{
			Time:  fromDate,
			Valid: true,
		}
	}

	if in.ValidTo != "" {
		toDate, err := time.Parse(layoutFormat, in.ValidTo)
		if err != nil {
			return arg, errors.New("to date is not format")
		}

		arg.ValidTo = sql.NullTime{
			Time:  toDate,
			Valid: true,
		}
	}

	if in.BaseUrl != "" {
		arg.BaseUrl = sql.NullString{
			String: in.BaseUrl,
			Valid:  true,
		}
	}

	if in.Method != "" {
		arg.Method = sql.NullString{
			String: in.Method,
			Valid:  true,
		}
	}

	if in.Inq != "" {
		arg.Inq = sql.NullString{
			String: in.Inq,
			Valid:  true,
		}
	}

	if in.Pay != "" {
		arg.Pay = sql.NullString{
			String: in.Pay,
			Valid:  true,
		}
	}

	if in.Adv != "" {
		arg.Adv = sql.NullString{
			String: in.Adv,
			Valid:  true,
		}
	}

	if in.Cmt != "" {
		arg.Cmt = sql.NullString{
			String: in.Adv,
			Valid:  true,
		}
	}

	if in.Rev != "" {
		arg.Rev = sql.NullString{
			String: in.Rev,
			Valid:  true,
		}
	}

	if in.Status == "" {
		arg.Status = "Active"
	}

	return arg, nil
}

func (o *ProviderService) setUpdateProvider(arg db.UpdateProviderParams, in dto.UpdateProviderRequest) (db.UpdateProviderParams, error) {
	layoutFormat := "2006-01-02 15:04:05"
	if in.ValidFrom != "" {
		fromDate, err := time.Parse(layoutFormat, in.ValidFrom)
		if err != nil {
			return arg, errors.New("from date is not format")
		}

		arg.ValidFrom = sql.NullTime{
			Time:  fromDate,
			Valid: true,
		}
		arg.SetValidFrom = true
	}

	if in.ValidTo != "" {
		toDate, err := time.Parse(layoutFormat, in.ValidTo)
		if err != nil {
			return arg, errors.New("to date is not format")
		}

		arg.ValidTo = sql.NullTime{
			Time:  toDate,
			Valid: true,
		}
		arg.SetValidTo = true
	}

	if in.Name != "" {
		arg.SetName = true
	}

	if in.User != "" {
		arg.SetUser = true
	}

	if in.Secret != "" {
		arg.SetSecret = true
	}

	if in.AddInfo1 != "" {
		arg.SetAddInfo1 = true
	}

	if in.AddInfo2 != "" {
		arg.SetAddInfo2 = true
	}

	if in.BaseUrl != "" {
		arg.BaseUrl = sql.NullString{
			String: in.BaseUrl,
			Valid:  true,
		}
		arg.SetBaseUrl = true
	}

	if in.Method != "" {
		arg.Method = sql.NullString{
			String: in.Method,
			Valid:  true,
		}
		arg.SetMethod = true
	}

	if in.Inq != "" {
		arg.Inq = sql.NullString{
			String: in.Inq,
			Valid:  true,
		}
		arg.SetInq = true
	}

	if in.Pay != "" {
		arg.Pay = sql.NullString{
			String: in.Pay,
			Valid:  true,
		}
		arg.SetPay = true
	}

	if in.Adv != "" {
		arg.Adv = sql.NullString{
			String: in.Adv,
			Valid:  true,
		}
		arg.SetAdv = true
	}

	if in.Rev != "" {
		arg.Rev = sql.NullString{
			String: in.Rev,
			Valid:  true,
		}
		arg.SetRev = true
	}

	if in.Cmt != "" {
		arg.Cmt = sql.NullString{
			String: in.Rev,
			Valid:  true,
		}
		arg.SetCmt = true
	}

	if in.Cmt != "" {
		arg.Cmt = sql.NullString{
			String: in.Cmt,
			Valid:  true,
		}
		arg.SetCmt = true
	}

	if in.Status != "" {
		arg.Status = in.Status
		arg.SetStatus = true
	}

	return arg, nil
}

func (o *ProviderService) ProviderRes(provider db.Provider) dto.ProviderRes {
	return dto.ProviderRes{
		ID:        provider.ID,
		Name:      provider.Name,
		User:      provider.User,
		Secret:    provider.Secret,
		AddInfo1:  provider.AddInfo1,
		AddInfo2:  provider.AddInfo2,
		ValidFrom: provider.ValidFrom.Time.Format("2006-01-02 15:04:05"),
		ValidTo:   provider.ValidTo.Time.Format("2006-01-02 15:04:05"),
		BaseUrl:   provider.BaseUrl.String,
		Method:    provider.Method.String,
		Inq:       provider.Inq.String,
		Pay:       provider.Pay.String,
		Adv:       provider.Adv.String,
		Cmt:       provider.Cmt.String,
		Rev:       provider.Rev.String,
		Status:    provider.Status,
	}
}
