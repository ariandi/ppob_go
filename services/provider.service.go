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

// ProviderService is
type ProviderService struct {
}

var providerService *ProviderService

// GetProviderService is
func GetProviderService() *ProviderService {
	if providerService == nil {
		providerService = new(ProviderService)
	}
	return providerService
}

func (o *ProviderService) CreateProviderService(req dto.CreateProviderReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.ProviderRes, error) {
	logrus.Println("[ProviderService CreateProviderService] start.")
	var result dto.ProviderRes

	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.CreateProviderParams{
		Name:      req.Name,
		User:      req.User,
		Secret:    req.Secret,
		AddInfo1:  req.AddInfo1,
		AddInfo2:  req.AddInfo2,
		BaseUrl:   sql.NullString{},
		Method:    sql.NullString{},
		Inq:       sql.NullString{},
		Pay:       sql.NullString{},
		Adv:       sql.NullString{},
		Cmt:       sql.NullString{},
		Rev:       sql.NullString{},
		Status:    "",
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg, err = o.setCreateProvider(arg, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return result, err
	}

	provider, err := store.CreateProvider(ctx, arg)
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

	result = ProviderRes(provider)
	return result, nil
}

func (o *ProviderService) GetProviderService(req dto.GetProviderReq, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.ProviderRes, error) {
	logrus.Println("[ProviderService GetProviderService] start.")
	var result dto.ProviderRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	prov, err := store.GetProvider(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return result, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	result = ProviderRes(prov)
	return result, nil
}

func (o *ProviderService) ListProviderService(req dto.ListProviderRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) ([]dto.ProviderRes, error) {
	logrus.Println("[ProviderService ListProviderService] start.")
	var result []dto.ProviderRes
	_, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.ListProviderParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	providers, err := store.ListProvider(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return result, err
	}

	for _, provider := range providers {
		u := ProviderRes(provider)
		result = append(result, u)
	}

	return result, nil
}

func (o *ProviderService) UpdateProviderService(req dto.UpdateProviderRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) (dto.ProviderRes, error) {
	logrus.Println("[ProviderService UpdateProviderService] start.")
	var result dto.ProviderRes
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return result, errors.New("error in user validator")
	}

	arg := db.UpdateProviderParams{
		ID:        req.ID,
		Name:      req.Name,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	arg, err = o.setUpdateProvider(arg, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse(err))
		return result, err
	}

	prov, err := store.UpdateProvider(ctx, arg)
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

	result = ProviderRes(prov)

	return result, nil
}

func (o *ProviderService) SoftDeleteProviderService(req dto.UpdateInactiveProviderRequest, authPayload *token.Payload, ctx *gin.Context, store db.Store) error {
	logrus.Println("[ProviderService SoftDeleteProviderService] start.")
	userValid, err := validator(store, ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveProviderParams{
		ID:        req.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = store.UpdateInactiveProvider(ctx, arg)
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

func (o *ProviderService) setCreateProvider(arg db.CreateProviderParams, req dto.CreateProviderReq) (db.CreateProviderParams, error) {
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
		toDate, err := time.Parse(layoutFormat, req.ValidTo)
		if err != nil {
			return arg, errors.New("to date is not format")
		}

		arg.ValidTo = sql.NullTime{
			Time:  toDate,
			Valid: true,
		}
	}

	if req.BaseUrl != "" {
		arg.BaseUrl = sql.NullString{
			String: req.BaseUrl,
			Valid:  true,
		}
	}

	if req.Method != "" {
		arg.Method = sql.NullString{
			String: req.Method,
			Valid:  true,
		}
	}

	if req.Inq != "" {
		arg.Inq = sql.NullString{
			String: req.Inq,
			Valid:  true,
		}
	}

	if req.Pay != "" {
		arg.Pay = sql.NullString{
			String: req.Pay,
			Valid:  true,
		}
	}

	if req.Adv != "" {
		arg.Adv = sql.NullString{
			String: req.Adv,
			Valid:  true,
		}
	}

	if req.Adv != "" {
		arg.Adv = sql.NullString{
			String: req.Adv,
			Valid:  true,
		}
	}

	if req.Rev != "" {
		arg.Rev = sql.NullString{
			String: req.Rev,
			Valid:  true,
		}
	}

	if req.Status == "" {
		arg.Status = "active"
	}

	return arg, nil
}

func (o *ProviderService) setUpdateProvider(arg db.UpdateProviderParams, req dto.UpdateProviderRequest) (db.UpdateProviderParams, error) {
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
		toDate, err := time.Parse(layoutFormat, req.ValidTo)
		if err != nil {
			return arg, errors.New("to date is not format")
		}

		arg.ValidTo = sql.NullTime{
			Time:  toDate,
			Valid: true,
		}
	}

	if req.BaseUrl != "" {
		arg.BaseUrl = sql.NullString{
			String: req.BaseUrl,
			Valid:  true,
		}
	}

	if req.Method != "" {
		arg.Method = sql.NullString{
			String: req.Method,
			Valid:  true,
		}
	}

	if req.Inq != "" {
		arg.Inq = sql.NullString{
			String: req.Inq,
			Valid:  true,
		}
	}

	if req.Pay != "" {
		arg.Pay = sql.NullString{
			String: req.Pay,
			Valid:  true,
		}
	}

	if req.Adv != "" {
		arg.Adv = sql.NullString{
			String: req.Adv,
			Valid:  true,
		}
	}

	if req.Adv != "" {
		arg.Adv = sql.NullString{
			String: req.Adv,
			Valid:  true,
		}
	}

	if req.Rev != "" {
		arg.Rev = sql.NullString{
			String: req.Rev,
			Valid:  true,
		}
	}

	if req.Status == "" {
		arg.Status = "active"
	}

	return arg, nil
}

func ProviderRes(provider db.Provider) dto.ProviderRes {
	return dto.ProviderRes{
		ID:        provider.ID,
		Name:      provider.Name,
		User:      provider.User,
		Secret:    provider.Secret,
		AddInfo1:  provider.AddInfo2,
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
