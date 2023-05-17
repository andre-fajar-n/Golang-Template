package rest

import (
	"go_template/gen/models"
	"go_template/gen/restapi/operations"
	"go_template/internal/utils/jwt"
	"go_template/runtime"
	"time"

	"github.com/go-openapi/errors"
)

func parseToken(rt *runtime.Runtime, token string) (*jwt.Payload, error) {
	secret := rt.Cfg.JwtSecret
	maker, err := jwt.NewJWTMaker(secret)
	if err != nil {
		return nil, err
	}

	payload, err := maker.VerifyToken(token)
	if err != nil {
		return nil, rt.SetError(401, "Unauthorized: invalid API key token: %v", err)
	}

	return payload, nil
}

func verifySingleRole(payload *jwt.Payload, role string) (*models.Principal, error) {
	if payload.Role != role {
		return nil, errors.New(403, "Forbidden: insufficient API key privileges")
	}

	return &models.Principal{
		UserID:    payload.UserID,
		ExpiredAt: payload.ExpiredAt.Format(time.RFC3339),
		Email:     payload.Email,
	}, nil
}

func checkHasRole(rt *runtime.Runtime, role, token string) (*models.Principal, error) {
	payload, err := parseToken(rt, token)
	if err != nil {
		return nil, err
	}

	p, err := verifySingleRole(payload, role)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func Authorization(rt *runtime.Runtime, api *operations.ServerAPI) {

}
