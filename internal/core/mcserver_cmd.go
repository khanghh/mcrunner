package core

import (
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type MCServerCmd struct {
	// configuration
	cmdPath string
	cmdArgs []string
	runDir  string
	stdout  io.Writer

	// runtime
	cmd *exec.Cmd

	// stdin writer to send commands
	stdinPipe io.WriteCloser
	stream    *OutputStream

	mu   sync.Mutex
	done chan struct{}
	err  error
}

// NewMCServerCmd creates a new MCServerCmd instance with proper initialization.
func NewMCServerCmd(cmdPath string, cmdArgs []string, runDir string, stdout io.Writer) *MCServerCmd {
	return &MCServerCmd{
		cmdPath: cmdPath,
		cmdArgs: cmdArgs,
		runDir:  runDir,
		stdout:  stdout,
		stream:  NewOutputStream(10),
		done:    make(chan struct{}),
	}
}

// SendCommand writes a command to the server stdin. A newline is appended
// if the provided command doesn't already end with one.
func (m *MCServerCmd) SendCommand(cmd string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stdinPipe == nil {
		return os.ErrInvalid
	}
	// ensure newline
	if len(cmd) == 0 || cmd[len(cmd)-1] != '\n' {
		cmd = cmd + "\n"
	}
	_, err := m.stdinPipe.Write([]byte(cmd))
	return err
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

// Start starts a Minecraft server process using the configured command and arguments.
func (m *MCServerCmd) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cmd := exec.Command(m.cmdPath, m.cmdArgs...)
	if m.runDir != "" {
		cmd.Dir = m.runDir
	}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	var stdoutWriter io.Writer = m.stream
	if m.stdout != nil {
		stdoutWriter = io.MultiWriter(m.stream, m.stdout)
	}

	cmd.Stdout = stdoutWriter
	cmd.Stderr = stdoutWriter

	if err := cmd.Start(); err != nil {
		_ = stdinPipe.Close()
		return err
	}

	m.cmd = cmd
	m.stdinPipe = stdinPipe

	go func() {
		mErr := cmd.Wait()
		m.mu.Lock()
		m.err = mErr
		m.stream.Close()
		close(m.done)
		m.mu.Unlock()
	}()

	return nil
}
