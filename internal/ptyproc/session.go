package ptyproc

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
)

// PTYSession manages a single process attached to a PTY.
// It provides Start/Stop and simple attach (stdin/stdout piping) plus resize.
type PTYSession struct {
	// immutable after Start
	name string
	cols int
	rows int

	// command configuration
	cmdPath string
	cmdArgs []string
	env     []string
	dir     string

	// runtime
	cmd   *exec.Cmd
	ptmx  *os.File
	mu    sync.Mutex
	alive bool

	// buffer for recent output
	buffer *ringBuffer

	// optional: close signaling
	doneOnce sync.Once
	doneCh   chan struct{}
}

func (s *PTYSession) Name() string {
	return s.name
}

func (s *PTYSession) Size() (int, int) {
	return s.cols, s.rows
}

// NewPTYSession creates a new PTY session configuration. You must call Start.
// If cmdPath is empty, it defaults to "/bin/bash".
func NewPTYSession(name string, cmdPath string, args []string, env []string, dir string, cols, rows int) *PTYSession {
	if cols <= 0 {
		cols = 80
	}
	if rows <= 0 {
		rows = 24
	}
	if cmdPath == "" {
		cmdPath = "/bin/bash"
	}
	return &PTYSession{
		name:    name,
		cols:    cols,
		rows:    rows,
		cmdPath: cmdPath,
		cmdArgs: append([]string(nil), args...),
		env:     append([]string(nil), env...),
		dir:     dir,
		buffer:  newRingBuffer(1 << 20), // 1 MiB buffer
		doneCh:  make(chan struct{}),
	}
}

// Start launches the configured command attached to a PTY with initial size.
func (s *PTYSession) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.alive {
		return errors.New("pty session already started")
	}
	cmd := exec.Command(s.cmdPath, s.cmdArgs...)
	if len(s.env) > 0 {
		cmd.Env = append(os.Environ(), s.env...)
	}
	if s.dir != "" {
		cmd.Dir = s.dir
	}
	// create PTY with size
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Cols: uint16(s.cols), Rows: uint16(s.rows)})
	if err != nil {
		return err
	}
	s.cmd = cmd
	s.ptmx = ptmx
	s.alive = true

	// monitor process exit to close doneCh
	go func() {
		// Wait will return when process exits
		_ = cmd.Wait()
		s.mu.Lock()
		s.alive = false
		s.mu.Unlock()
		s.doneOnce.Do(func() { close(s.doneCh) })
		s.closePTY()
	}()
	return nil
}

// Restart stops the process (TERM -> KILL fallback) and starts it again with same config.
func (s *PTYSession) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}
	return s.Start()
}

// Stop attempts a graceful termination (SIGTERM) and waits briefly without killing.
func (s *PTYSession) Stop() error {
	s.mu.Lock()
	cmd := s.cmd
	alive := s.alive
	s.mu.Unlock()
	if !alive || cmd == nil || cmd.Process == nil {
		return nil
	}
	// send SIGTERM
	_ = cmd.Process.Signal(syscall.SIGTERM)
	// wait with timeout, but do not force kill

	select {
	case <-s.doneCh:
		s.closePTY()
		return nil
	case <-time.After(1 * time.Minute):
		return errors.New("timeout waiting for process to exit")
	}
}

// Kill forcefully terminates the process (SIGKILL).
func (s *PTYSession) Kill() error {
	// Capture the underlying process pointer and liveness under lock to avoid races
	s.mu.Lock()
	alive := s.alive
	var proc *os.Process
	if s.cmd != nil {
		proc = s.cmd.Process
	}
	s.mu.Unlock()
	if !alive || proc == nil {
		return nil
	}
	if err := proc.Kill(); err != nil {
		return err
	}
	// Proactively close PTY to unblock any readers/writers.
	s.closePTY()
	return nil
}

// PTY returns the session PTY as an io.ReadWriter.
// Returns an error if the session is not running.
func (s *PTYSession) PTY() (io.ReadWriter, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ptmx == nil || !s.alive {
		return nil, errors.New("pty not started")
	}
	return s.ptmx, nil
}

func (s *PTYSession) Wait() error {
	<-s.doneCh
	return nil
}

// Attach wires the session PTY to provided stdout/stdin until either side closes.
// Caller is responsible for terminal modes if attaching to a real TTY.
// Attach functionality removed; PTYSession runs headless.

// Resize sets the PTY window size.
func (s *PTYSession) Resize(cols, rows int) error {
	s.mu.Lock()
	ptmx := s.ptmx
	s.cols, s.rows = cols, rows
	s.mu.Unlock()
	if ptmx == nil {
		return errors.New("pty not started")
	}
	return pty.Setsize(ptmx, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

// closePTY safely closes s.ptmx once and nils the reference.
func (s *PTYSession) closePTY() {
	s.mu.Lock()
	pt := s.ptmx
	s.ptmx = nil
	s.mu.Unlock()
	if pt != nil {
		_ = pt.Close()
	}
}
