package api

import (
	"net/http"
	"time"
)

type FileSystemAPI struct {
	baseURL    string
	httpClient *http.Client
}

func NewFileSystemAPI(baseURL string) *FileSystemAPI {
	return &FileSystemAPI{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
