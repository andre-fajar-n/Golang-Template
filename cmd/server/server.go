package server

import (
	"fmt"
	"go_template/gen/restapi"
	"go_template/gen/restapi/operations"
	"go_template/internal/handlers"
	"go_template/internal/repositories"
	"go_template/internal/rest"
	"go_template/runtime"
	"log"
	"os"

	"github.com/casualjim/middlewares"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/swag"
	"github.com/jessevdk/go-flags"
	"github.com/justinas/alice"
)

var mainFlags = struct {
	AppConfig string `long:"config" description:"Main application configuration YAML path"`
}{}

func Main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalf("Error : %f", err)
	}

	api := operations.NewServerAPI(swaggerSpec)
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "App Flags",
			LongDescription:  "",
			Options:          &mainFlags,
		},
	}

	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "go template"
	parser.LongDescription = "go template"
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	rt := runtime.NewRuntime()

	userRepo := repositories.Newuser(*rt)

	h := handlers.NewHandler(*rt, userRepo)

	rest.Authorization(rt, api)
	rest.Route(rt, api, h)

	api.Logger = func(s string, i ...interface{}) {
		msg := "Logger: " + s
		if i != nil {
			msg = fmt.Sprintf(msg, i)
		}
		rt.Logger.Info().Msg(msg)
	}

	handler := alice.New(
		middlewares.NewRecoveryMW("golang-template", nil),
		middlewares.NewProfiler,
	).Then(api.Serve(nil))

	server.SetHandler(handler)
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
