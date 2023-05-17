package rest

import (
	"go_template"
	"go_template/gen/restapi/operations"
	"go_template/gen/restapi/operations/health"
	"go_template/internal/handlers"

	"github.com/go-openapi/runtime/middleware"
)

func Route(rt *go_template.Runtime, api *operations.ServerAPI, apiHandler handlers.Handler) {
	//  health
	api.HealthHealthHandler = health.HealthHandlerFunc(func(hp health.HealthParams) middleware.Responder {
		return health.NewHealthOK().WithPayload(&health.HealthOKBody{
			Message: "Server up",
		})
	})

}
