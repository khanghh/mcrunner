package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/internal/core"
	"github.com/khanghh/mcrunner/internal/handlers"
	"github.com/khanghh/mcrunner/internal/params"
	"github.com/urfave/cli/v2"
)

var (
	app       *cli.App
	gitCommit string
	gitDate   string
	gitTag    string
)

var (
	commandFlag = &cli.StringFlag{
		Name:    "command",
		Aliases: []string{"cmd", "c"},
		Usage:   "Minecraft server command to run",
	}
	rootDirFlag = &cli.StringFlag{
		Name:    "rootdir",
		Aliases: []string{"dir", "d"},
		Usage:   "Minecraft server root directory",
	}
	listenFlag = &cli.StringFlag{
		Name:    "listen",
		Aliases: []string{"l"},
		Usage:   "HTTP server listen address",
		Value:   ":3000",
	}
)

func init() {
	app = cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "Minecraft server runner"
	app.Usage = ""
	app.Flags = []cli.Flag{
		commandFlag,
		rootDirFlag,
		listenFlag,
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

// run is the main entry point for the CLI application.
// It initializes and starts the Minecraft server command, sets up HTTP API routes,
// and handles graceful shutdown on receiving termination signals.
func run(cli *cli.Context) error {
	rootDir := cli.String(rootDirFlag.Name)
	listenAddr := cli.String(listenFlag.Name)
	serverCmd := cli.String(commandFlag.Name)
	if serverCmd == "" {
		return fmt.Errorf("server command must not be empty")
	}

	mcserverCmd := core.NewMCServerCmd(serverCmd, []string{}, rootDir, os.Stdout)

	// middlewares
	var (
		wsUpgradeRequired = func(ctx *fiber.Ctx) error {
			if !fiberws.IsWebSocketUpgrade(ctx) {
				return fiber.ErrUpgradeRequired
			}
			return ctx.Next()
		}
	)

	// handlers
	mcrunnerHandler := handlers.NewMCRunnerHandler(mcserverCmd)

	router := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          handlers.ErrorHandler,
	})

	router.Get("/status", mcrunnerHandler.GetStatus)
	router.Post("/command", mcrunnerHandler.PostCommand)
	router.Post("/start", mcrunnerHandler.PostStartServer)
	router.Post("/stop", mcrunnerHandler.PostStopServer)
	router.Post("/kill", mcrunnerHandler.PostKillServer)
	router.Get("/ws", wsUpgradeRequired, mcrunnerHandler.WebsocketHandler())

	// start the mcserver command and serve http API
	if err := mcserverCmd.Start(); err != nil {
		return fmt.Errorf("failed to start Minecraft server command: %v", err)
	}

	// Handle signals: first triggers graceful shutdown, second forces exit
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		go func() {
			_ = mcserverCmd.Stop()
			_ = router.Shutdown()
		}()
		<-sigCh
		os.Exit(1)
	}()

	return router.Listen(listenAddr)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
