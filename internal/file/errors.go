package file

import "errors"

var (
	ErrPathTraversal  = errors.New("invalid path: traversal outside root is not allowed")
	ErrNotFound       = errors.New("path not found")
	ErrIsDirectory    = errors.New("path is a directory")
	ErrNotDirectory   = errors.New("path is not a directory")
	ErrAlreadyExists  = errors.New("already exists")
	ErrDirNotEmpty    = errors.New("directory not empty")
	ErrMissingNewName = errors.New("missing new name")
)
