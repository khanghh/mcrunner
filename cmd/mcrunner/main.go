package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/khanghh/mcrunner/internal/core"
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

func run(cli *cli.Context) error {
	serverCmd := cli.String(commandFlag.Name)
	if serverCmd == "" {
		return fmt.Errorf("server command must not be empty")
	}

	rootDir := cli.String(rootDirFlag.Name)
	listenAddr := cli.String(listenFlag.Name)

	unixSockPath := "/tmp/mcrunner.sock"
	unixLogWriter, err := core.NewUnixLogWriter(unixSockPath)
	if err != nil {
		return err
	}

	stdoutWriter := io.MultiWriter(os.Stdout, unixLogWriter)
	mcserver, err := core.RunMinecraftServer(serverCmd, []string{}, rootDir, stdoutWriter)
	if err != nil {
		return err
	}

	// serve http server in background
	go serveHttp(listenAddr, mcserver)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigCh {
			mcserver.Signal(sig)
		}
	}()

	return mcserver.Wait()
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
