package core

import (
	"io"
	"os"
)

type LocalFileService interface {
	Stat(relPath string) (os.FileInfo, error)
	List(relPath string) ([]os.FileInfo, error)
	Open(relPath string) (*os.File, os.FileInfo, error)
	ReadFile(relPath string) ([]byte, error)
	WriteFile(relPath string, data []byte, create bool) error
	SaveStream(relPath string, reader io.Reader, overwrite bool) error
	Delete(relPath string) error
	DeleteRecursive(relPath string) error
	MkdirAll(relPath string) error
	Rename(oldRelPath, newRelPath string, overwrite bool) error
	DetectMIMEType(relPath string) (string, error)
}

type MCManagerService interface {
	StartServer() error
	StopServer() error
	KillServer() error
	RestartServer() error
	ServerStatus() *ServerState
	AttachConsole() (io.ReadWriteCloser, error)
}
