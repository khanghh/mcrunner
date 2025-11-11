package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khanghh/mcrunner/internal/core"
	"github.com/khanghh/mcrunner/internal/handlers"
	"github.com/khanghh/mcrunner/internal/params"
	"github.com/khanghh/mcrunner/internal/websocket"
	"github.com/khanghh/mcrunner/pkg/gen"
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
	inputFifoFlag = &cli.StringFlag{
		Name:    "fifo",
		Aliases: []string{"f"},
		Usage:   "Path to input FIFO file for sending commands to the Minecraft server",
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
		inputFifoFlag,
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

func ensureFifoExist(fifoPath string) error {
	if _, statErr := os.Stat(fifoPath); errors.Is(statErr, os.ErrNotExist) {
		if mkErr := syscall.Mkfifo(fifoPath, 0666); mkErr != nil && !os.IsExist(mkErr) {
			return fmt.Errorf("mkfifo %s: %v", fifoPath, mkErr)
		}
	} else if statErr != nil {
		return fmt.Errorf("stat %s: %v", fifoPath, statErr)
	}
	return nil
}

func fifoInputLoop(mcserverCmd *core.MCServerCmd, fifoPath string) {
	if err := ensureFifoExist(fifoPath); err != nil {
		panic(err)
	}
	for {
		// Attempt to open FIFO for reading (will block until a writer opens it if it exists)
		fifoFile, err := os.OpenFile(fifoPath, os.O_RDONLY, 0600)
		if err != nil {
			if os.IsNotExist(err) {
				ensureFifoExist(fifoPath)
			}
			time.Sleep(time.Second)
			continue
		}

		buf := make([]byte, 4096)
		for {
			n, readErr := fifoFile.Read(buf)
			if n > 0 && mcserverCmd.GetStatus() == core.StateRunning {
				if _, wErr := mcserverCmd.Write(buf[:n]); wErr != nil {
					fmt.Fprintf(os.Stderr, "write stdin failed: %v\n", wErr)
				}
			}
			if readErr != nil {
				if readErr != io.EOF {
					fmt.Fprintf(os.Stderr, "read fifo error: %v\n", readErr)
				}
				break
			}
		}
		fifoFile.Close()
	}
}

func mustResolveRootDir(rootDir string) string {
	absPath, err := filepath.Abs(rootDir)
	if err != nil {
		panic("failed to resolve absolute path: " + err.Error())
	}

	resolved, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		panic("failed to resolve symlinks: " + err.Error())
	}

	info, err := os.Stat(resolved)
	if err != nil {
		if os.IsNotExist(err) {
			panic("rootdir does not exist")
		}
		panic("failed to stat rootdir: " + err.Error())
	}

	if !info.IsDir() {
		panic("rootdir must be a directory")
	}

	return resolved
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

	absRootDir := mustResolveRootDir(rootDir)
	localFilesSvc := core.NewLocalFileService(absRootDir)
	mcserverCmd := core.NewMCServerCmd(serverCmd, []string{}, rootDir, os.Stdout)
	if fifoPath := cli.String(inputFifoFlag.Name); fifoPath != "" {
		go fifoInputLoop(mcserverCmd, fifoPath)
	}

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
	fsHandler := handlers.NewFSHandler(localFilesSvc)

	// setup websocket server
	wsServer := websocket.NewServer()
	wsServer.StartBroadcast(mcrunnerHandler.WSBroadcast)
	wsServer.OnConnect(mcrunnerHandler.WSOnClientConnect)
	wsServer.OnMessage(gen.MessageType_PTY_INPUT, mcrunnerHandler.WSHandlePTYInput)

	// setup HTTP server and routes
	router := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          handlers.ErrorHandler,
	})
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "*",
	}))
	router.Get("/api/fs/*", fsHandler.Get)
	router.Post("/api/fs/*", fsHandler.Post)
	router.Put("/api/fs/*", fsHandler.Put)
	router.Patch("/api/fs/*", fsHandler.Patch)
	router.Delete("/api/fs/*", fsHandler.Delete)
	router.Get("/api/mc/status", mcrunnerHandler.GetStatus)
	router.Post("/api/mc/command", mcrunnerHandler.PostCommand)
	router.Post("/api/mc/start", mcrunnerHandler.PostStartServer)
	router.Post("/api/mc/stop", mcrunnerHandler.PostStopServer)
	router.Post("/api/mc/restart", mcrunnerHandler.PostRestartServer)
	router.Post("/api/mc/kill", mcrunnerHandler.PostKillServer)
	router.Get("/ws", wsUpgradeRequired, wsServer.ServeFiberWS())

	// start the mcserver command
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
			_ = wsServer.Shutdown()
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
