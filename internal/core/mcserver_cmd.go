package core

import (
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type MCServerCmd struct {
	// command configuration
	path string
	args []string

	// runtime
	cmd    *exec.Cmd
	stdout io.Writer

	// stdin writer to send commands
	stdinPipe io.WriteCloser

	mu   sync.Mutex
	done chan struct{}
	err  error
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
	return m.Signal(syscall.SIGTERM)
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

// RunMinecraftServer starts a Minecraft server process with the given command and arguments.
func RunMinecraftServer(cmdPath string, cmdArgs []string, runDir string, stdout io.Writer) (*MCServerCmd, error) {
	m := &MCServerCmd{
		path:   cmdPath,
		args:   cmdArgs,
		stdout: stdout,
		done:   make(chan struct{}),
	}
	cmd := exec.Command(m.path, m.args...)
	if runDir != "" {
		cmd.Dir = runDir
	}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stdout = m.stdout
	cmd.Stderr = m.stdout

	if err := cmd.Start(); err != nil {
		_ = stdinPipe.Close()
		return nil, err
	}

	m.cmd = cmd
	m.stdinPipe = stdinPipe
	m.done = make(chan struct{})

	go func() {
		mErr := cmd.Wait()
		m.mu.Lock()
		m.err = mErr
		close(m.done)
		m.mu.Unlock()
	}()

	return m, nil
}
