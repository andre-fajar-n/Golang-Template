package rest

import (
	"go_template/gen/restapi/operations"
	"go_template/gen/restapi/operations/health"
	"go_template/internal/handlers"
	"go_template/runtime"

	"github.com/go-openapi/runtime/middleware"
)

func Route(rt *runtime.Runtime, api *operations.ServerAPI, apiHandler handlers.Handler) {
	//  health
	api.HealthHealthHandler = health.HealthHandlerFunc(func(hp health.HealthParams) middleware.Responder {
		return health.NewHealthOK().WithPayload(&health.HealthOKBody{
			Message: "Server up",
		})
	})

}
