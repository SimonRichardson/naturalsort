package fs

import (
	"io"
	"os"
)

type realFilesystem struct{}

// NewRealFilesystem yields a real disk filesystem
func NewRealFilesystem() Filesystem {
	return realFilesystem{}
}

func (realFilesystem) Create(path string) (file File, err error) {
	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return
	}

	return realFile{
		File:   f,
		Reader: f,
		Closer: f,
	}, nil
}

func (fs realFilesystem) Open(path string) (file File, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}

	return realFile{
		File:   f,
		Reader: f,
		Closer: f,
	}, nil
}

func (realFilesystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

type realFile struct {
	*os.File
	io.Reader
	io.Closer
}

func (f realFile) Read(p []byte) (int, error) {
	return f.Reader.Read(p)
}

func (f realFile) Close() error {
	return f.Closer.Close()
}
