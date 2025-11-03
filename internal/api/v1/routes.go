package api

import (
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/khanghh/mcrunner/internal/core"
)

// Runner represents the subset of server runner functionality handlers need.
type MCServer interface {
	Start() error
	Stop() error
	Restart() error
	Status() core.ServerStatus
	SendCommand(cmd string) error
}

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

func SetupRoutes(router fiber.Router, lfs LocalFileService, mcserver MCServer) error {
	serverHandler := NewServerHandler(mcserver)
	fsHandler := NewFSHandler(lfs)
	api := router.Group("/api/v1")

	// Server management
	api.Get("/server/status", serverHandler.GetStatus)
	api.Post("/server/start", serverHandler.PostStart)
	api.Post("/server/stop", serverHandler.PostStop)
	api.Post("/server/restart", serverHandler.PostRestart)
	api.Post("/server/command", serverHandler.PostCommand)

	// File system
	api.Get("/fs/*", fsHandler.Get)
	api.Post("/fs/*", fsHandler.Post)
	api.Put("/fs/*", fsHandler.Put)
	api.Patch("/fs/*", fsHandler.Patch)
	api.Delete("/fs/*", fsHandler.Delete)
	return nil
}
