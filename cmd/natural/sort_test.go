package main

import (
	"compress/gzip"
	"testing"

	"io"

	"encoding/base64"

	"bytes"

	"github.com/SimonRichardson/naturalsort/pkg/fs"
)

func TestRead(t *testing.T) {
	t.Parallel()

	t.Run("input", func(t *testing.T) {
		fsys := fs.NewVirtualFilesystem()

		content := "0,0001,0,23,5,a3,43123"
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
	})

	t.Run("input file", func(t *testing.T) {
		content := "0,0001,0,23,5,a3,43123"
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
