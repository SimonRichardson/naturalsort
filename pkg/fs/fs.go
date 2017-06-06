package fs

import (
	"io"
)

type Filesystem interface {
	Create(path string) (File, error)
	Open(path string) (File, error)
	Exists(path string) bool
}

type File interface {
	io.Reader
	io.Writer
	io.Closer
}
