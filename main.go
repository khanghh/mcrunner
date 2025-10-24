package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	api "github.com/khanghh/mcrunner/internal/api/v1"
	"github.com/khanghh/mcrunner/internal/core"
	"github.com/khanghh/mcrunner/internal/params"

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
	configFileFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "YAML config file",
		Value: "config.yaml",
	}
	debugFlag = &cli.BoolFlag{
		Name:  "debug",
		Usage: "Enable debug logging",
	}
)

func init() {
	app = cli.NewApp()
	app.EnableBashCompletion = true
	app.Usage = ""
	app.Flags = []cli.Flag{
		configFileFlag,
		debugFlag,
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
	fmt.Printf(" Version:\t%s\n", gitTag)
	fmt.Printf(" Commit:\t%s\n", gitCommit)
	fmt.Printf(" Built Time:\t%s\n", gitDate)
	return nil
}

func mustInitLogger(debug bool) {
	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))
}

// Mock server runner for demo purposes
type mockServerRunner struct{}

func (m *mockServerRunner) Start() error    { return nil }
func (m *mockServerRunner) Stop() error     { return nil }
func (m *mockServerRunner) IsRunning() bool { return true }
func (m *mockServerRunner) SendCommand(cmd string) error {
	slog.Info("Mock send command", "cmd", cmd)
	return nil
}
func (m *mockServerRunner) GetOutputChannel() <-chan string { return nil }
func (m *mockServerRunner) GetErrorChannel() <-chan string  { return nil }

func run(cli *cli.Context) error {
	config, err := core.LoadConfig(cli.String(configFileFlag.Name))
	if err != nil {
		slog.Error("Could not load config file.", "error", err)
		return err
	}

	mustInitLogger(config.Debug || cli.IsSet(debugFlag.Name))

	router := fiber.New(fiber.Config{
		CaseSensitive: true,
		BodyLimit:     params.ServerBodyLimit,
		IdleTimeout:   params.ServerIdleTimeout,
		ReadTimeout:   params.ServerReadTimeout,
		WriteTimeout:  params.ServerWriteTimeout,
	})

	mockServer := &mockServerRunner{}
	lfs := core.NewLocalFileService(config.RootDir)

	if err := api.SetupRoutes(router, lfs, mockServer); err != nil {
		slog.Error("Failed to setup routes", "error", err)
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		slog.Info("Stopping server...")
		router.Shutdown()
	}()

	return router.Listen(config.ListenAddr)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
