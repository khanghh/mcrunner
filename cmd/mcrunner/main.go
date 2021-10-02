package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mcrunner/pkg/logger"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/itzg/go-flagsfiller"
)

var (
	stdin       io.Writer
	stdout      io.Reader
	stderr      io.Reader
	cmdPipe     io.Reader
	cmdCh       chan string
	cmdExitChan chan int
)

var (
	AppName    string
	Version    string
	CommitHash string
	BuiltTime  string
	OsArch     string
)

func init() {
	AppName = "Minecraft server runner"
	OsArch = runtime.GOOS + "/" + runtime.GOARCH
}

func printVersion() {
	fmt.Println(AppName)
	fmt.Printf(" Version:\t%s\n", Version)
	fmt.Printf(" Commit:\t%s\n", CommitHash)
	fmt.Printf(" Built Time:\t%s\n", BuiltTime)
	fmt.Printf(" OS/Arch:\t%s\n", OsArch)
}

type Args struct {
	Debug                   bool          `usage:"Enable debug logging"`
	Bootstrap               string        `usage:"Specifies a file with commands to initially send to the server"`
	StopDuration            time.Duration `usage:"Amount of time in Golang duration to wait after sending the 'stop' command."`
	StopServerAnnounceDelay time.Duration `default:"0s" usage:"Amount of time in Golang duration to wait after announcing server shutdown"`
	DetachStdin             bool          `usage:"Don't forward stdin and allow process to be put in background"`
	Shell                   string        `usage:"When set, pass the arguments to this shell"`
	CmdPipe                 string        `usage:"Specifies a fifo file to pipe minecraft command to stdin"`
	Version                 bool          `usage:"Show version information"`
}

func main() {
	sigCh := make(chan os.Signal, 1)
	// docker stop sends a SIGTERM, so intercept that and send a 'stop' command to the server
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	cmdExitChan = make(chan int, 1)
	cmdCh = make(chan string, 1)

	var args Args
	err := flagsfiller.Parse(&args)
	if err != nil {
		log.Fatal(err)
	}

	if args.Version {
		printVersion()
		return
	}

	if args.Debug {
		logger.LogLevel = logger.Debug
	}

	var cmd *exec.Cmd

	if flag.NArg() < 1 {
		logger.Fatalln("Missing executable arguments")
	}

	if args.Shell != "" {
		cmd = exec.Command(args.Shell, flag.Args()...)
	} else {
		cmd = exec.Command(flag.Arg(0), flag.Args()[1:]...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if args.StopDuration == 0 {
		args.StopDuration = 60 * time.Second
	}

	if args.CmdPipe != "" {
		os.Remove(args.CmdPipe)
		err := syscall.Mkfifo(args.CmdPipe, 0666)
		if err != nil {
			logger.Fatalln("Make named pipe file error:", err)
		}
		cmdPipe, err = os.OpenFile(args.CmdPipe,
			syscall.O_CREAT|syscall.O_RDONLY|syscall.O_RDWR, os.ModeNamedPipe)
		if err != nil {
			logger.Errorln("Unable to read named pipe file", err)
		}
	}

	stdin, err = cmd.StdinPipe()
	if err != nil {
		logger.Fatalln("Unable to get stdin", err)
	}

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		logger.Fatalln("Unable to get stdout", err)
	}

	stderr, err = cmd.StderrPipe()
	if err != nil {
		logger.Fatalln("Unable to get stderr", err)
	}

	err = cmd.Start()
	if err != nil {
		logger.Fatalln("Failed to start", err)
	}

	if args.Bootstrap != "" {
		bootstrapContent, err := ioutil.ReadFile(args.Bootstrap)
		if err != nil {
			logger.Fatalln("Failed to read bootstrap commands", err)
		}
		_, err = stdin.Write(bootstrapContent)
		if err != nil {
			logger.Fatalln("Failed to write bootstrap content", err)
		}
	}

	// Relay stdout/stderr from server to runner
	go func() {
		io.Copy(os.Stdout, stdout)
	}()

	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	if !args.DetachStdin {
		go func() {
			for cmd := range cmdCh {
				stdin.Write([]byte(cmd))
			}
		}()
		go pipeCmd(os.Stdin)
		if cmdPipe != nil {
			go pipeCmd(cmdPipe)
		}
	}

	go func() {
		waitErr := cmd.Wait()
		if waitErr != nil {
			if exitErr, ok := waitErr.(*exec.ExitError); ok {
				exitCode := exitErr.ExitCode()
				logger.Warnf("sub-process failed. exitCode: %d\n", exitCode)
				cmdExitChan <- exitCode
			} else {
				logger.Errorln("Command failed abnormally. ", waitErr)
				cmdExitChan <- 1
			}
			return
		} else {
			cmdExitChan <- 0
		}
	}()

	isStopping := false
	for {
		select {
		case <-sigCh:
			if isStopping {
				continue
			}
			isStopping = true
			if args.StopServerAnnounceDelay > 0 {
				sendCmd(fmt.Sprintf("say Server shutting down in %0.f seconds\n", args.StopServerAnnounceDelay.Seconds()))
				logger.Printf("Sleeping %0.f seconds before stopping server\n", args.StopServerAnnounceDelay.Seconds())
				time.Sleep(args.StopServerAnnounceDelay)
			}
			sendCmd("save-all\n")
			sendCmd("stop\n")
			logger.Println("Waiting for server to stopped...")
			time.AfterFunc(args.StopDuration, func() {
				logger.Errorln("Still not stopped, so killing server process")
				err := cmd.Process.Kill()
				if err != nil {
					logger.Errorln("Failed to shutdown server.")
				}
			})
		case exitCode := <-cmdExitChan:
			logger.Println("Server stopped.")
			os.Exit(exitCode)
		}
	}
}

func sendCmd(cmd string) {
	cmdCh <- cmd
}

func pipeCmd(input io.Reader) {
	reader := bufio.NewReader(input)
	for {
		cmd, _ := reader.ReadString('\n')
		if len(cmd) != 0 {
			sendCmd(cmd)
		}
	}
}
