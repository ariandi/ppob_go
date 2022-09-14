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

type MediaInterface interface {
	CreateMediaService(ctx *gin.Context, in dto.CreateMediaReq) (dto.MediaRes, error)
	GetMediaService(ctx *gin.Context, in dto.GetMediaReq) (dto.MediaRes, error)
	ListMediaService(ctx *gin.Context, in dto.ListMediaRequest) ([]dto.MediaRes, error)
	UpdateMediaService(ctx *gin.Context, in dto.UpdateMediaRequest) (dto.MediaRes, error)
	SoftDeleteMediaService(ctx *gin.Context, in dto.UpdateInactiveMediaRequest) error
	MediaResponse(cat db.Madiastorage) dto.MediaRes
}

// MediaService is
type MediaService struct {
	store db.Store
}

var mediaService *MediaService

// GetMediaService is
func GetMediaService(store db.Store) MediaInterface {
	if mediaService == nil {
		mediaService = &MediaService{
			store: store,
		}
	}
	return mediaService
}

func (o *MediaService) CreateMediaService(ctx *gin.Context, in dto.CreateMediaReq) (dto.MediaRes, error) {
	logrus.Println("[MediaService CreateMediaService] start.")
	var out dto.MediaRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.CreateMediaStorageParams{
		SecID:     in.SecID,
		TabID:     in.TabID,
		Name:      in.Name,
		Type:      in.Type,
		Content:   in.Content,
		CreatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	media, err := o.store.CreateMediaStorage(ctx, arg)
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

	out = o.MediaResponse(media)
	return out, nil
}

func (o *MediaService) GetMediaService(ctx *gin.Context, in dto.GetMediaReq) (dto.MediaRes, error) {
	logrus.Println("[MediaService GetMediaService] start.")
	var out dto.MediaRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.GetMediaStorageParams{
		ID:    in.ID,
		SecID: in.SecID,
		TabID: in.TabID,
	}
	if in.ID > 0 {
		arg.IsID = true
	}

	if in.SecID != "" {
		arg.IsSec = true
	}

	if in.TabID != "" {
		arg.IsTab = true
	}

	media, err := o.store.GetMediaStorage(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse(err))
			return out, err
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	out = o.MediaResponse(media)
	return out, nil
}

func (o *MediaService) ListMediaService(ctx *gin.Context, in dto.ListMediaRequest) ([]dto.MediaRes, error) {
	logrus.Println("[MediaService ListMediaService] start.")
	var out []dto.MediaRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	_, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.ListMediaStorageParams{
		Limit:  in.PageSize,
		Offset: (in.PageID - 1) * in.PageSize,
	}

	medias, err := o.store.ListMediaStorage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse(err))
		return out, err
	}

	for _, media := range medias {
		u := o.MediaResponse(media)
		out = append(out, u)
	}

	return out, nil
}

func (o *MediaService) UpdateMediaService(ctx *gin.Context, in dto.UpdateMediaRequest) (dto.MediaRes, error) {
	logrus.Println("[MediaService UpdateMediaService] start.")
	var out dto.MediaRes

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return out, errors.New("error in user validator")
	}

	arg := db.UpdateMediaStorageParams{
		SetName:   true,
		Name:      in.Name,
		SetType:   true,
		Type:      in.Type,
		UpdatedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
		ID:        in.ID,
	}

	media, err := o.store.UpdateMediaStorage(ctx, arg)
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

	out = o.MediaResponse(media)

	return out, nil
}

func (o *MediaService) SoftDeleteMediaService(ctx *gin.Context, in dto.UpdateInactiveMediaRequest) error {
	logrus.Println("[MediaService SoftDeleteMediaService] start.")

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	userValid, err := userService.validator(ctx, authPayload)
	if err != nil {
		return errors.New("error in user validator")
	}

	arg := db.UpdateInactiveMediaStorageParams{
		ID:        in.ID,
		DeletedBy: sql.NullInt64{Int64: userValid.ID, Valid: true},
	}

	_, err = o.store.UpdateInactiveMediaStorage(ctx, arg)
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

func (o *MediaService) MediaResponse(media db.Madiastorage) dto.MediaRes {
	return dto.MediaRes{
		ID:        media.ID,
		SecID:     media.SecID,
		TabID:     media.TabID,
		Name:      media.Name,
		Type:      media.Type,
		Content:   media.Content,
		CreatedAt: media.CreatedAt.Time,
		UpdatedAt: media.UpdatedAt.Time,
		CreatedBy: media.CreatedBy.Int64,
		UpdatedBy: media.UpdatedBy.Int64,
	}
}
