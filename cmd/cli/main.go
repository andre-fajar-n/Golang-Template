package main

import (
	"fmt"
	"go_template/cmd/server"
	"go_template/runtime"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	appName  = "Golang Template"
	appUsage = "CLI to run this apps"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage
	app.Commands = []*cli.Command{
		api, migration,
	}
	app.CommandNotFound = commanNotFound

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func commanNotFound(ctx *cli.Context, s string) {
	fmt.Println("Command not found. Please use command \"help\" or \"h\" to show help")
}

var (
	api = &cli.Command{
		Name:  "api",
		Usage: "Run API Server",
		Action: func(ctx *cli.Context) error {
			server.Main()
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "port"},
			&cli.StringFlag{Name: "host"},
		},
	}

	migration = &cli.Command{
		Name:  "migration",
		Usage: "Run migration needs",
		Subcommands: []*cli.Command{
			migrateUp, migrateDown, migrateCreateFile, migrateForce,
		},
	}

	migrateUp = &cli.Command{
		Name:  "up",
		Usage: "Migrate up migrations",
		Action: func(ctx *cli.Context) error {
			rt := runtime.NewRuntime()

			rt.MigrateUp()

			return nil
		},
	}

	migrateDown = &cli.Command{
		Name:  "down",
		Usage: "Migrate down migrations",
		Action: func(ctx *cli.Context) error {
			rt := runtime.NewRuntime()

			rt.MigrateDown()

			return nil
		},
	}

	migrateForce = &cli.Command{
		Name:  "force",
		Usage: "Fix latest dirty migration",
		Action: func(ctx *cli.Context) error {
			rt := runtime.NewRuntime()

			rt.ForceLastestVersion()

			return nil
		},
	}

	migrateCreateFile = &cli.Command{
		Name:  "create",
		Usage: "Create new file migrations",
		Action: func(ctx *cli.Context) error {
			rt := runtime.NewRuntime()

			args := ctx.Args()

			if args.Len() < 1 {
				return fmt.Errorf("please insert filename")
			}

			filename := args.First()
			rt.CreateFileMigration(filename)

			return nil
		},
	}
)
