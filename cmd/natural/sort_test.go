package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
	"testing/quick"

	"github.com/SimonRichardson/naturalsort/pkg/fs"
)

func TestRead(t *testing.T) {
	t.Parallel()

	t.Run("input", func(t *testing.T) {
		fn := func(a []string) bool {
			content := strings.Join(a, ",")

			fsys := fs.NewVirtualFilesystem()
			reader, err := read(fsys, content, "", false, false)
			if err != nil {
				t.Fatal(err)
			}

			buf := make([]byte, len(content))
			if _, err := io.ReadFull(reader, buf); err != nil {
				t.Fatal(err)
			}

			if expected, actual := content, string(buf); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}
			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("input file", func(t *testing.T) {
		fn := func(a []string) bool {
			content := strings.Join(a, ",")

			fsys := fs.NewVirtualFilesystem()
			path := "tmpfile"
			file, err := fsys.Create(path)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := file.Write([]byte(content)); err != nil {
				t.Fatal(err)
			}

			reader, err := read(fsys, content, "", false, false)
			if err != nil {
				t.Fatal(err)
			}

			buf := make([]byte, len(content))
			if _, err := io.ReadFull(reader, buf); err != nil {
				t.Fatal(err)
			}

			if expected, actual := content, string(buf); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}
			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("input base64", func(t *testing.T) {
		fsys := fs.NewVirtualFilesystem()

		content := "0,0001,0,23,5,a3,43123"
		baseContent := base64.StdEncoding.EncodeToString([]byte(content))
		reader, err := read(fsys, baseContent, "", false, true)
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, len(content))
		if _, err := io.ReadFull(reader, buf); err != nil {
			t.Fatal(err)
		}

		if expected, actual := content, string(buf); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})

	t.Run("input gzip", func(t *testing.T) {
		fsys := fs.NewVirtualFilesystem()

		content := "0,0001,0,23,5,a3,43123"

		var b bytes.Buffer
		w := gzip.NewWriter(&b)
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
		w.Close()

		reader, err := read(fsys, b.String(), "", true, false)
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, len(content))
		if _, err := io.ReadFull(reader, buf); err != nil {
			t.Fatal(err)
		}

		if expected, actual := content, string(buf); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})
}

func TestWrite(t *testing.T) {
	t.Parallel()

	t.Run("output file", func(t *testing.T) {
		fn := func(a []string) bool {
			fsys := fs.NewVirtualFilesystem()
			path := fmt.Sprintf("tmpfile-%d", rand.Intn(99999))

			content := strings.Join(a, ",")
			buf := bytes.NewBufferString(content)
			if err := write(fsys, path, false, false)(buf); err != nil {
				t.Fatal(err)
			}

			file, err := fsys.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			r := make([]byte, len(content))
			if _, err := io.ReadFull(file, r); err != nil {
				t.Fatal(err)
			}

			if expected, actual := content, string(r); expected != actual {
				t.Errorf("expected: %q, actual: %q", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("output gzip", func(t *testing.T) {
		fsys := fs.NewVirtualFilesystem()
		path := fmt.Sprintf("tmpfile-%d", rand.Intn(99999))

		content := "a,m,1,3,b,f12,12c,41,e"
		buf := bytes.NewBufferString(content)
		if err := write(fsys, path, true, false)(buf); err != nil {
			t.Fatal(err)
		}

		file, err := fsys.Open(path)
		if err != nil {
			t.Fatal(err)
		}

		reader, err := gzip.NewReader(file)
		if err != nil {
			t.Fatal(err)
		}

		r := make([]byte, len(content))
		if _, err := io.ReadFull(reader, r); err != nil {
			t.Fatal(err)
		}

		if expected, actual := content, string(r); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})

	t.Run("output base64", func(t *testing.T) {
		fsys := fs.NewVirtualFilesystem()
		path := fmt.Sprintf("tmpfile-%d", rand.Intn(99999))

		content := "a,m,1,3,b,f12,12c,41,e"
		buf := bytes.NewBufferString(content)
		if err := write(fsys, path, false, true)(buf); err != nil {
			t.Fatal(err)
		}

		file, err := fsys.Open(path)
		if err != nil {
			t.Fatal(err)
		}

		reader := base64.NewDecoder(base64.StdEncoding, file)

		r := make([]byte, len(content))
		if _, err := io.ReadFull(reader, r); err != nil {
			t.Fatal(err)
		}

		if expected, actual := content, string(r); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
	})
}
