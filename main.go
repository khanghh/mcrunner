package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/khanghh/mcrunner/internal/file"
	"github.com/khanghh/mcrunner/internal/handlers"
	"github.com/khanghh/mcrunner/internal/mccmd"
	"github.com/khanghh/mcrunner/internal/params"
	"github.com/khanghh/mcrunner/internal/service"
	pb "github.com/khanghh/mcrunner/pkg/proto"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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
		Aliases: []string{"d"},
		Usage:   "File manager root directory",
	}
	inputFifoFlag = &cli.StringFlag{
		Name:    "fifo",
		Aliases: []string{"f"},
		Usage:   "Path to input FIFO file for sending commands to the Minecraft server",
	}
	grpcListenFlag = &cli.StringFlag{
		Name:    "grpc",
		Aliases: []string{"g"},
		Usage:   "gRPC server listen address (host:port)",
		Value:   ":50051",
	}
	httpListenFlag = &cli.StringFlag{
		Name:    "http",
		Aliases: []string{"l"},
		Usage:   "HTTP server listen address (host:port)",
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
		grpcListenFlag,
		httpListenFlag,
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

func fifoInputLoop(mcserverCmd *mccmd.MCServerCmd, fifoPath string) {
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
			if n > 0 && mcserverCmd.GetStatus() == mccmd.StatusRunning {
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

func parseServerCmd(commandStr string) (string, []string) {
	parts := strings.Fields(commandStr)
	cmdPath := parts[0]
	if len(parts) > 1 {
		return cmdPath, parts[1:]
	}
	return commandStr, []string{}
}

func initListeners(grpcAddr, httpAddr string) (net.Listener, net.Listener, error) {
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen on gRPC address %s: %v", grpcAddr, err)
	}

	if grpcAddr == httpAddr {
		return grpcListener, grpcListener, nil
	}

	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		grpcListener.Close()
		return nil, nil, fmt.Errorf("failed to listen on HTTP address %s: %v", httpAddr, err)
	}
	return grpcListener, httpListener, nil
}

// run is the main entry point for the CLI application.
// It initializes and starts the Minecraft server command, sets up HTTP API routes,
// and handles graceful shutdown on receiving termination signals.
func run(cli *cli.Context) error {
	rootDir := cli.String(rootDirFlag.Name)
	gprcListenAddr := cli.String(grpcListenFlag.Name)
	httpListenAddr := cli.String(httpListenFlag.Name)
	serverCmd := cli.String(commandFlag.Name)
	if serverCmd == "" {
		return fmt.Errorf("server command must not be empty")
	}

	absRootDir := mustResolveRootDir(rootDir)
	localFilesSvc := file.NewLocalFileService(absRootDir)

	cmdPath, cmdArgs := parseServerCmd(serverCmd)
	mcserverCmd := mccmd.NewMCServerCmd(cmdPath, cmdArgs, rootDir, os.Stdout)
	if fifoPath := cli.String(inputFifoFlag.Name); fifoPath != "" {
		go fifoInputLoop(mcserverCmd, fifoPath)
	}

	// handlers
	mcrunnerHandler := handlers.NewMCRunnerHandler(mcserverCmd)
	fsHandler := handlers.NewFSHandler(localFilesSvc)

	// setup HTTP server and routes
	router := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          handlers.ErrorHandler,
		BodyLimit:             256 * 1024 * 1024,
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
	router.Get("/api/mc/state", mcrunnerHandler.GetState)
	router.Post("/api/mc/command", mcrunnerHandler.PostCommand)
	router.Post("/api/mc/start", mcrunnerHandler.PostStartServer)
	router.Post("/api/mc/stop", mcrunnerHandler.PostStopServer)
	router.Post("/api/mc/restart", mcrunnerHandler.PostRestartServer)
	router.Post("/api/mc/kill", mcrunnerHandler.PostKillServer)
	router.Get("/readyz", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	mcrunnerSvc := service.NewMCRunnerService(mcserverCmd)
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 0,
			Time:              1 * time.Minute,  // ping every 60s
			Timeout:           10 * time.Second, // wait up to 10s for pong
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second, // clients must wait at least this between pings
			PermitWithoutStream: true,             // allow pings even with no active RPC
		}),
	)
	pb.RegisterMCRunnerServer(server, mcrunnerSvc)

	// Handle signals: first triggers graceful shutdown, second forces exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		go func() {
			mcserverCmd.Stop()
			server.GracefulStop()
			router.Shutdown()
			close(sigCh)
		}()
		<-sigCh
		os.Exit(143)
	}()

	// start the mcserver command
	if err := mcserverCmd.Start(); err != nil {
		return fmt.Errorf("failed to start Minecraft server command: %v", err)
	}
	go func() {
		io.Copy(mcserverCmd, os.Stdin)
	}()

	grpcListener, httpListener, err := initListeners(gprcListenAddr, httpListenAddr)
	if err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		fmt.Printf("Listening gRPC at %s\n", gprcListenAddr)
		if err := server.Serve(grpcListener); err != nil {
			errCh <- fmt.Errorf("gRPC server error: %v", err)
		}
	}()
	go func() {
		fmt.Printf("Listening HTTP at %s\n", httpListenAddr)
		if err := router.Listener(httpListener); err != nil {
			errCh <- fmt.Errorf("HTTP server error: %v", err)
		}
	}()
	return <-errCh
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
