package fs

import (
	"bytes"
	"os"
	"sync"
)

type virtualFilesystem struct {
	mutex sync.RWMutex
	files map[string]*virtualFile
}

// NewVirtualFilesystem yields an in-memory filesystem.
func NewVirtualFilesystem() Filesystem {
	return &virtualFilesystem{
		files: map[string]*virtualFile{},
	}
}

func (v *virtualFilesystem) Create(path string) (File, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	// os.Create truncates any existing files. So we do, too.
	f := &virtualFile{
		name: path,
	}
	v.files[path] = f
	return f, nil
}

func (v *virtualFilesystem) Open(path string) (File, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	f, ok := v.files[path]
	if !ok {
		return nil, os.ErrNotExist
	}
	return f, nil
}

func (v *virtualFilesystem) Exists(path string) bool {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	_, ok := v.files[path]
	return ok
}

type virtualFile struct {
	name  string
	mutex sync.Mutex
	buf   bytes.Buffer
}

func (v *virtualFile) Read(p []byte) (int, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	return v.buf.Read(p)
}

func (v *virtualFile) Write(p []byte) (int, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	return v.buf.Write(p)
}

func (v *virtualFile) Close() error { return nil }
