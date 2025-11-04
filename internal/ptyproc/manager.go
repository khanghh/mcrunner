package ptyproc

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

type PTYManager struct {
	mu       sync.RWMutex
	sessions map[string]*PTYSession
}

type Options struct {
	// Optional logical session identifier. If empty, a unique id will be generated.
	Name string

	Command string   // Command to run (required)
	Args    []string // Command arguments (optional)
	Env     []string // Environment variables (optional)
	Dir     string   // Working directory (optional)

	Cols int // Terminal columns (optional, default: 80)
	Rows int // Terminal rows (optional, default: 24)

	Stdout io.Writer // Redirect PTY stdout
	Stdin  io.Reader // Attach PTY stdin
}

// NewSession creates and starts a PTY session per the options and stores it by name.
// If opts.Name is empty, a unique name is generated and returned.
func (m *PTYManager) NewSession(opts Options) (*PTYSession, error) {
	if opts.Command == "" {
		opts.Command = "/bin/bash"
	}
	name := opts.Name
	if name == "" {
		name = fmt.Sprintf("sess-%d", time.Now().UnixNano())
	}

	sess := NewPTYSession(opts)

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sessions == nil {
		m.sessions = make(map[string]*PTYSession)
	}

	if _, exists := m.sessions[name]; !exists {
		m.sessions[name] = sess
	} else {
		return nil, fmt.Errorf("session %q already exists", name)
	}

	if err := sess.Start(); err != nil {
		delete(m.sessions, name)
		return nil, err
	}

	return sess, nil
}

func (m *PTYManager) Shutdown() error {
	m.mu.RLock()
	names := make([]string, 0, len(m.sessions))
	for name := range m.sessions {
		names = append(names, name)
	}
	m.mu.RUnlock()
	var firstErr error
	for _, name := range names {
		log.Printf("Stopping session %s\n", name)
		if err := m.Stop(name); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func NewPTYManager() *PTYManager {
	return &PTYManager{
		sessions: make(map[string]*PTYSession),
	}
}

// Get returns the session by name.
func (m *PTYManager) Get(name string) (*PTYSession, bool) {
	m.mu.RLock()
	s, ok := m.sessions[name]
	m.mu.RUnlock()
	return s, ok
}

// PTY returns the PTY ReadWriter for a session.
func (m *PTYManager) PTY(name string) (io.ReadWriter, error) {
	s, ok := m.Get(name)
	if !ok {
		return nil, fmt.Errorf("session %q not found", name)
	}
	return s.PTY()
}

// Write writes data to the session's PTY stdin.
func (m *PTYManager) Write(name string, p []byte) (int, error) {
	rw, err := m.PTY(name)
	if err != nil {
		return 0, err
	}
	return rw.Write(p)
}

// Resize updates the PTY size for the session.
func (m *PTYManager) Resize(name string, cols, rows int) error {
	s, ok := m.Get(name)
	if !ok {
		return fmt.Errorf("session %q not found", name)
	}
	return s.Resize(cols, rows)
}

// Stop sends SIGTERM and waits for exit with timeout.
func (m *PTYManager) Stop(name string) error {
	s, ok := m.Get(name)
	if !ok {
		return fmt.Errorf("session %q not found", name)
	}
	err := s.Stop()
	// keep session in map in case caller wants exit info; could remove if desired
	return err
}

// Kill forcefully terminates the session.
func (m *PTYManager) Kill(name string) error {
	s, ok := m.Get(name)
	if !ok {
		return fmt.Errorf("session %q not found", name)
	}
	return s.Kill()
}

// Remove removes a stopped session from the manager map.
func (m *PTYManager) Remove(name string) {
	m.mu.Lock()
	delete(m.sessions, name)
	m.mu.Unlock()
}

func (m *PTYManager) ListSessions() []string {
	m.mu.RLock()
	names := make([]string, 0, len(m.sessions))
	for name := range m.sessions {
		names = append(names, name)
	}
	m.mu.RUnlock()
	return names
}
