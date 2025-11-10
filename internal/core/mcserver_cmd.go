package core

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
)

// ServerStatus represents the current server status
type ServerStatus string

const (
	StatusRunning  ServerStatus = "running"
	StatusStopping ServerStatus = "stopping"
	StatusStopped  ServerStatus = "stopped"
)

type MCServerCmd struct {
	// configuration
	cmdPath string
	cmdArgs []string
	cmdDir  string

	// runtime
	cmd  *exec.Cmd
	ptmx *os.File

	stream       *outputStream
	outputWriter io.Writer

	mu        sync.Mutex
	done      chan struct{}
	err       error
	startTime *time.Time
	status    ServerStatus
}

// NewMCServerCmd creates a new MCServerCmd instance with proper initialization.
func NewMCServerCmd(cmdPath string, cmdArgs []string, runDir string, stdout io.Writer) *MCServerCmd {
	stream := newOutputStream(10)
	return &MCServerCmd{
		cmdPath:      cmdPath,
		cmdArgs:      cmdArgs,
		cmdDir:       runDir,
		stream:       stream,
		outputWriter: io.MultiWriter(stdout, stream),
		done:         make(chan struct{}),
		status:       StatusStopped,
	}
}

// SendCommand writes a command to the server stdin. A newline is appended
// if the provided command doesn't already end with one.
func (m *MCServerCmd) SendCommand(cmd string) error {
	if !strings.HasSuffix(cmd, "\n") {
		cmd += "\n"
	}
	_, err := m.Write([]byte(cmd))
	return err
}

// Write writes data to the server command's stdin.
func (m *MCServerCmd) Write(data []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ptmx == nil {
		return 0, os.ErrInvalid
	}
	return m.ptmx.Write(data)
}
func (m *MCServerCmd) Read(p []byte) (int, error) {
	return 0, nil
}

// Wait blocks until the Minecraft server process exits.
func (m *MCServerCmd) Wait() error {
	<-m.done
	return m.err
}

// Stop attempts to gracefully stop the Minecraft server by sending SIGTERM.
func (m *MCServerCmd) Stop() error {
	if err := m.Signal(syscall.SIGTERM); err != nil {
		return err
	}
	m.mu.Lock()
	m.status = StatusStopping
	m.mu.Unlock()
	return m.Wait()
}

// Signal sends a signal to the underlying Minecraft server process.
func (m *MCServerCmd) Signal(sig os.Signal) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cmd == nil || m.cmd.Process == nil {
		return os.ErrInvalid
	}
	return m.cmd.Process.Signal(sig)
}

// Kill forcefully terminates the Minecraft server process.
func (m *MCServerCmd) Kill() error {
	if m.cmd == nil || m.cmd.Process == nil {
		return os.ErrInvalid
	}
	return m.cmd.Process.Kill()
}

func (m *MCServerCmd) OutputStream() io.Reader {
	return m.stream
}

// GetStatus returns the current server status
func (m *MCServerCmd) GetStatus() ServerStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

// GetProcess returns the underlying process
func (m *MCServerCmd) GetProcess() *os.Process {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cmd == nil {
		return nil
	}
	return m.cmd.Process
}

// GetStartTime returns the server start time
func (m *MCServerCmd) GetStartTime() *time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.startTime
}

// Start starts a Minecraft server process using the configured command and arguments.
func (m *MCServerCmd) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cmd := exec.Command(m.cmdPath, m.cmdArgs...)
	if m.cmdDir != "" {
		cmd.Dir = m.cmdDir
	}

	// Start the command with PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	m.cmd = cmd
	m.ptmx = ptmx
	m.status = StatusRunning
	now := time.Now()
	m.startTime = &now
	m.done = make(chan struct{})

	go io.Copy(m.outputWriter, ptmx)

	// Wait for command to finish
	go func() {
		mErr := cmd.Wait()
		m.mu.Lock()
		m.err = mErr
		m.status = StatusStopped
		ptmx.Close()
		close(m.done)
		m.mu.Unlock()
	}()

	return nil
}
