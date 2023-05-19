package handlers

import (
	"context"
	"go_template/gen/restapi/operations/authentication"
	"go_template/internal/repositories"
	"go_template/runtime"
)

type (
	handler struct {
		runtime  runtime.Runtime
		userRepo repositories.User
	}

	Handler interface {
		userHandler
	}

	userHandler interface {
		Register(ctx context.Context, req authentication.RegisterParams) (*uint64, error)
		Login(ctx context.Context, req *authentication.LoginParams) (token, expiredAt *string, err error)
	}
)

func NewHandler(
	rt runtime.Runtime,
	userRepo repositories.User,
) Handler {
	return &handler{
		rt,
		userRepo,
	}
}
