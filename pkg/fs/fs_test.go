package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReal(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	t.Run("create", func(t *testing.T) {
		fsys := NewRealFilesystem()
		path := filepath.Join(dir, "tmpfile")
		file, err := fsys.Create(path)
		if err != nil {
			t.Error(err)
		}

		defer file.Close()

		if !fsys.Exists(path) {
			t.Errorf("expected: %q to exist", path)
		}
	})

	t.Run("open", func(t *testing.T) {
		content := []byte("hello world")
		tmpfile, err := ioutil.TempFile(dir, "tmpfile")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name())
		if _, err := tmpfile.Write(content); err != nil {
			log.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}

		fsys := NewRealFilesystem()
		path := tmpfile.Name()
		if !fsys.Exists(path) {
			t.Fatalf("expected: %q to exist", path)
		}

		file, err := fsys.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		buf := make([]byte, len(content))
		if _, err := io.ReadFull(file, buf); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(content, buf) {
			t.Errorf("expected: %v, actual: %v", content, buf)
		}
	})
}

func TestVirtual(t *testing.T) {
	t.Parallel()

	t.Run("create", func(t *testing.T) {
		fsys := NewVirtualFilesystem()
		path := fmt.Sprintf("tmpfile-%d", rand.Intn(1000))
		file, err := fsys.Create(path)
		if err != nil {
			t.Error(err)
		}

		defer file.Close()

		if !fsys.Exists(path) {
			t.Errorf("expected: %q to exist", path)
		}
	})

	t.Run("open", func(t *testing.T) {
		fsys := NewVirtualFilesystem()
		path := fmt.Sprintf("tmpfile-%d", rand.Intn(1000))
		tmpfile, err := fsys.Create(path)
		if err != nil {
			t.Error(err)
		}

		content := []byte("hello world")
		if _, err := tmpfile.Write(content); err != nil {
			t.Fatal(err)
		}

		file, err := fsys.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		buf := make([]byte, len(content))
		if _, err := io.ReadFull(file, buf); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(content, buf) {
			t.Errorf("expected: %v, actual: %v", content, buf)
		}
	})
}
