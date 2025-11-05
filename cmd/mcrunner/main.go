package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/khanghh/mcrunner/internal/params"
	"github.com/khanghh/mcrunner/internal/ptyproc"
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
	serverCmd := cli.String("command")
	if serverCmd == "" {
		return fmt.Errorf("server command must not be empty")
	}

	rootDir := cli.String("rootdir")
	listenAddr := cli.String("listen")

	mcserver := ptyproc.NewPTYSession(ptyproc.Options{
		Name:    "mcserver",
		Command: serverCmd,
		Dir:     rootDir,
		Stdout:  os.Stdout,
		Cols:    80,
		Rows:    24,
	})
	if err := mcserver.Start(); err != nil {
		panic(err)
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
