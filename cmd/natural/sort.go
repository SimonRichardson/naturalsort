package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/SimonRichardson/naturalsort/pkg/fs"
	"github.com/SimonRichardson/naturalsort/pkg/group"
	"github.com/SimonRichardson/naturalsort/pkg/natural"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

const (
	defaultSeparator  = ","
	defaultInputGzip  = false
	defaultOutputGzip = false
)

// runSort performs the sorting of the input
func runSort(args []string) error {
	// flags for the sort command
	var (
		flagset = flag.NewFlagSet("sort", flag.ExitOnError)

		debug      = flagset.Bool("debug", false, "debug logging")
		separator  = flagset.String("separator", defaultSeparator, "separation value")
		input      = flagset.String("input", "", "input for natural sorting")
		inputFile  = flagset.String("input.file", "", "file required to perform natural sorting on")
		inputGzip  = flagset.Bool("input.gzip", defaultInputGzip, "decode gzip input")
		outputFile = flagset.String("output.file", "", "output file for action performed")
		outputGzip = flagset.Bool("output.gzip", defaultOutputGzip, "encode gzip output")
	)
	flagset.Usage = usageFor(flagset, "sort [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	// Setup the logger.
	var logger log.Logger
	{
		logLevel := level.AllowInfo()
		if *debug {
			logLevel = level.AllowAll()
		}
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = level.NewFilter(logger, logLevel)
	}

	level.Debug(logger).Log("type", "input", "file", *inputFile, "gzip", *inputGzip)
	level.Debug(logger).Log("type", "output", "file", *outputFile, "gzip", *outputGzip)

	// Validate the separator
	if len(*separator) != 1 {
		return errorFor(flagset, "sort [flags]", errors.Errorf("invalid separator (separator: %q)", *separator))
	}

	sepRune, size := utf8.DecodeLastRuneInString(*separator)
	if size == 0 {
		return errorFor(flagset, "sort [flags]", errors.Errorf("no valid separator (separator: %q)", *separator))
	}
	splitFn := splitOn(sepRune)

	// Validate that we either have an input or a input.file. If neither are
	// valid then bail out.
	in, inf := strings.TrimSpace(*input), strings.TrimSpace(*inputFile)
	if in == "" && inf == "" {
		return errorFor(flagset, "sort [flags]", errors.Errorf("no valid input (input: %q, file: %q)", in, inf))
	}

	// Execution group.
	var g group.Group
	{
		// Create the file system
		fsys := fs.NewRealFilesystem()
		g.Add(func() error {
			// Setup how we're going to read and write.
			reader, err := read(fsys, in, inf, *inputGzip)
			if err != nil {
				return err
			}
			defer reader.Close()

			writer := write(fsys, *outputFile, *outputGzip)

			// Work out how we're going to split then join on the input.
			iso := splitJoin{
				Split: splitFn,
				Join: func(x []string) string {
					return strings.Join(x, *separator)
				},
			}

			return naturalSort(iso, reader, writer)
		}, func(error) {
			// Nothing to close
		})
	}
	{
		// Setup os signal interruptions.
		cancel := make(chan struct{})
		g.Add(func() error {
			return interrupt(cancel)
		}, func(error) {
			close(cancel)
		})
	}
	return g.Run()
}

func read(fsys fs.Filesystem, input, inputFile string, inputGzip bool) (reader io.ReadCloser, err error) {
	// Read the file in a execution group, in case the file is huge.
	// That way the cmd is still responsive for potential feedback.
	if fsys.Exists(inputFile) {
		var file fs.File
		file, err = fsys.Open(inputFile)
		if err != nil {
			return
		}

		reader = file
	} else {
		// File doesn't exist, but the path does
		if inputFile != "" {
			err = errors.Errorf("file does not exist (input.file: %q)", inputFile)
			return
		}
		// No file reader found, so default to `input` flag
		reader = readCloser{strings.NewReader(input)}
	}

	if inputGzip {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return
		}
	}

	return
}

func write(fsys fs.Filesystem, outputFile string, outgzip bool) writeFn {
	return func(buf *bytes.Buffer) (err error) {
		// Work out where to write to.
		var writer io.Writer
		if outFile := strings.TrimSpace(outputFile); outFile != "" {
			var file fs.File
			file, err = fsys.Create(outFile)
			defer file.Close()
			if err != nil {
				return err
			}

			writer = file
		} else {
			writer = os.Stdout
		}

		if outgzip {
			w := gzip.NewWriter(writer)
			defer w.Close()

			writer = w
		}

		// Write the output
		_, err = buf.WriteTo(writer)
		return
	}
}

func naturalSort(iso splitJoin, reader io.Reader, writer writeFn) error {
	// Scan everything!
	scanner := bufio.NewScanner(reader)
	scanner.Split(iso.Split)

	var buf []string
	for scanner.Scan() {
		buf = append(buf, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Remove the last trailing `\n` of some files
	if last := buf[len(buf)-1]; strings.ContainsRune(last, '\n') {
		last = strings.TrimRightFunc(last, func(r rune) bool {
			return r == '\n'
		})
		buf[len(buf)-1] = last
	}

	// Perform the sorting
	natural.Sort(buf)

	// Create a buffer so that writing to sources becomes more natural
	out := bytes.NewBufferString(iso.Join(buf))
	return writer(out)

}

func splitOn(r rune) bufio.SplitFunc {
	if r == ' ' {
		return bufio.ScanWords
	}

	// Note this comes from the stdlib https://golang.org/src/bufio/example_test.go
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == byte(r) {
				return i + 1, data[:i], nil
			}
		}
		return 0, data, bufio.ErrFinalToken
	}
}

type writeFn func(*bytes.Buffer) error

type splitJoin struct {
	Split bufio.SplitFunc
	Join  func([]string) string
}

type readCloser struct {
	io.Reader
}

func (readCloser) Close() error {
	return nil
}
