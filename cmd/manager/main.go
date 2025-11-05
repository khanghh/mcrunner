package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2/middleware/logger"
	api "github.com/khanghh/mcrunner/internal/api/v1"
	"github.com/khanghh/mcrunner/internal/core"
	"github.com/khanghh/mcrunner/internal/params"

	_ "embed"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
)

var (
	app       *cli.App
	gitCommit string
	gitDate   string
	gitTag    string
)

var (
	rootDir = &cli.StringFlag{
		Name:    "rootdir",
		Aliases: []string{"d"},
		Usage:   "Minecraft server root directory to serve web editor files from",
	}
	staticDir = &cli.StringFlag{
		Name:    "staticdir",
		Aliases: []string{"s"},
		Usage:   "static folder to serve web assets",
		Value:   "./dist",
	}
	listenFlag = &cli.StringFlag{
		Name:    "listen",
		Aliases: []string{"l"},
		Usage:   "HTTP server listen address",
		Value:   ":8080",
	}
	mcrunnerAPIFlag = &cli.StringFlag{
		Name:    "server",
		Aliases: []string{"mc"},
		Usage:   "URL of the mcrunner API",
	}
)

func init() {
	app = cli.NewApp()
	app.EnableBashCompletion = true
	app.Usage = ""
	app.Flags = []cli.Flag{
		rootDir,
		staticDir,
		listenFlag,
		mcrunnerAPIFlag,
	}
	app.Commands = []*cli.Command{
		{
			Name:   "version",
			Action: printVersion,
		},
	}
	app.Action = run
}

func printVersion(cli *cli.Context) error {
	fmt.Println(cli.App.Name)
	fmt.Printf(" Version:\t%s\n", params.Version)
	fmt.Printf(" Commit:\t%s\n", gitCommit)
	fmt.Printf(" Built Time:\t%s\n", gitDate)
	return nil
}

func run(cli *cli.Context) error {
	rootDir := cli.String(rootDir.Name)
	if rootDir == "" {
		return fmt.Errorf("root directory must not be empty")
	}
	staticDir := cli.String(staticDir.Name)
	listenAddr := cli.String(listenFlag.Name)

	lfs := core.NewLocalFileService(rootDir)

	router := fiber.New(fiber.Config{
		CaseSensitive: true,
		BodyLimit:     params.ServerBodyLimit,
		IdleTimeout:   params.ServerIdleTimeout,
		ReadTimeout:   params.ServerReadTimeout,
		WriteTimeout:  params.ServerWriteTimeout,
	})
	router.Use(logger.New())
	router.Static("/", staticDir)
	if err := api.SetupRoutes(router.Group("/api"), lfs); err != nil {
		slog.Error("Failed to setup routes", "error", err)
		return err
	}

	return router.Listen(listenAddr)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
