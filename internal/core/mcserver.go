package core

import (
	"os"

	"github.com/khanghh/mcrunner/internal/ptyproc"
)

type ServerStatus int

const (
	StatusStopped ServerStatus = iota
	StatusRunning
	StatusStopping
)

type MCServer struct {
	proc *ptyproc.PTYSession
}

func (s *MCServer) PTYSession() *ptyproc.PTYSession {
	return s.proc
}

func (s *MCServer) Start() error {
	return s.proc.Start()
}

func (s *MCServer) Restart() error {
	return s.proc.Restart()
}

func (s *MCServer) Stop() error {
	return s.proc.Stop()
}

func (s *MCServer) SendCommand(cmd string) error {
	ptmx, err := s.proc.PTY()
	if err != nil {
		return err
	}
	_, _ = ptmx.Write([]byte(cmd + "\n"))
	return nil
}

func (s *MCServer) Status() ServerStatus {
	return StatusRunning
}

func RunMCServer(ptyManager *ptyproc.PTYManager, serverDir string, serverCmd string) (*MCServer, error) {
	opts := ptyproc.Options{
		Name:    "mcserver",
		Command: serverCmd,
		Dir:     serverDir,
		Stdout:  os.Stdout,
		Stdin:   os.Stdin,
	}
	proc, err := ptyManager.NewSession(opts)
	if err != nil {
		return nil, err
	}
	return &MCServer{proc: proc}, nil
}
