package rest

import (
	"context"
	"go_template/gen/models"
	"go_template/gen/restapi/operations"
	"go_template/gen/restapi/operations/authentication"
	"go_template/gen/restapi/operations/health"
	"go_template/internal/handlers"
	"go_template/runtime"

	"github.com/go-openapi/runtime/middleware"
)

func Route(rt *runtime.Runtime, api *operations.ServerAPI, apiHandler handlers.Handler) {
	//  health
	api.HealthHealthHandler = health.HealthHandlerFunc(func(hp health.HealthParams) middleware.Responder {
		return health.NewHealthOK().WithPayload(&models.Success{
			Message: "Server up",
		})
	})

	api.AuthenticationRegisterHandler = authentication.RegisterHandlerFunc(func(rp authentication.RegisterParams) middleware.Responder {
		userID, err := apiHandler.Register(context.Background(), rp)
		if err != nil {
			errRes := rt.GetError(err)
			return authentication.NewRegisterDefault(int(errRes.Code())).WithPayload(&models.Error{
				Code:    int64(errRes.Code()),
				Message: errRes.Error(),
			})
		}
		return authentication.NewRegisterOK().WithPayload(&models.SuccessRegister{
			Success: models.Success{
				Message: "success register",
			},
			SuccessRegisterAllOf1: models.SuccessRegisterAllOf1{
				UserID: *userID,
			},
		})
	})

	api.AuthenticationLoginHandler = authentication.LoginHandlerFunc(func(lp authentication.LoginParams) middleware.Responder {
		token, expiredAt, err := apiHandler.Login(context.Background(), &lp)
		if err != nil {
			errRes := rt.GetError(err)
			return authentication.NewLoginDefault(int(errRes.Code())).WithPayload(&models.Error{
				Code:    int64(errRes.Code()),
				Message: errRes.Error(),
			})
		}

		return authentication.NewLoginCreated().WithPayload(&models.SuccessLogin{
			Success: models.Success{
				Message: "success login",
			},
			SuccessLoginAllOf1: models.SuccessLoginAllOf1{
				ExpiredAt: *expiredAt,
			},
		}).WithToken(*token)
	})
}
