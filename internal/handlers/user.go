package handlers

import (
	"context"
	"go_template/gen/models"
	"go_template/gen/restapi/operations/authentication"
	utilsHash "go_template/internal/utils/hash"
	utilsJwt "go_template/internal/utils/jwt"
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"
)

func (h *handler) Register(ctx context.Context, req authentication.RegisterParams) (*uint64, error) {
	logger := h.runtime.Logger.With().
		Interface("req", req.Data).
		Logger()

	existUsername, err := h.userRepo.UsernameExist(ctx, *req.Data.Username)
	if err != nil {
		logger.Error().Err(err).Msg("error userRepo.UsernameExist")
		return nil, err
	}
	if existUsername {
		return nil, h.runtime.SetError(http.StatusBadRequest, "username already exist")
	}

	now := time.Now().UTC()
	nowStrfmt := strfmt.DateTime(now)
	user, err := h.userRepo.Create(ctx, nil, &models.User{
		ModelTrackTime: models.ModelTrackTime{
			CreatedAt: &nowStrfmt,
		},
		UserData: models.UserData{
			Username: *req.Data.Username,
			Password: utilsHash.HashSha256(*req.Data.Password),
			Note: models.CustomFields{
				"testing": 123,
			},
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("error userRepo.Create")
		return nil, err
	}

	return user.ID, nil
}

func (h *handler) Login(ctx context.Context, req *authentication.LoginParams) (token, expiredAt *string, err error) {
	logger := h.runtime.Logger.With().
		Interface("req", req.Data).
		Logger()

	user, err := h.userRepo.FindBySingleColumn(ctx, "username", req.Data.Username, false)
	if err != nil {
		logger.Error().Err(err).Msg("error userRepo.FindBySingleColumn")
		return
	}

	if user.Password != utilsHash.HashSha256(*req.Data.Password) {
		return nil, nil, h.runtime.SetError(http.StatusForbidden, "wrong password")
	}

	secret := h.runtime.Cfg.JwtSecret
	maker, err := utilsJwt.NewJWTMaker(secret)
	if err != nil {
		logger.Error().Err(err).Msg("error NewJWTMaker")
		return
	}

	tempToken, payload, err := maker.CreateToken(*user.ID, user.Username, time.Duration(h.runtime.Cfg.JwtExp)*time.Second)
	if err != nil {
		logger.Error().Err(err).Msg("error CreateToken")
		return
	}
	token = &tempToken
	tempExp := payload.ExpiredAt.Format(time.RFC3339)
	expiredAt = &tempExp

	return
}
