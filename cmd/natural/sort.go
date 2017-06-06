package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/SimonRichardson/naturalsort/pkg/fs"
	"github.com/SimonRichardson/naturalsort/pkg/group"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

const (
	defaultSeparator  = " "
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

	level.Debug(logger).Log("type", "input", "file", inputFile, "gzip", inputGzip)
	level.Debug(logger).Log("type", "output", "file", outputFile, "gzip", outputGzip)

	// Validate the separator
	sepRune, size := utf8.DecodeLastRuneInString(*separator)
	if size == 0 {
		return errorFor(flagset, "sort [flags]", errors.Errorf("no valid separator (separator: %q)", separator))
	}
	splitFn := splitOn(sepRune)

	// Validate that we either have an input or a input.file. If neither are
	// valid then bail out.
	in, inf := strings.TrimSpace(*input), strings.TrimSpace(*inputFile)
	if in == "" && inf == "" {
		return errorFor(flagset, "sort [flags]", errors.Errorf("no valid input (input: %q, file: %q)", in, inf))
	}

	// Create the file system
	fsys := fs.NewRealFilesystem()

	// Execution group.
	var g group.Group
	{
		g.Add(func() error {
			// Read the file in a execution group, in case the file is huge.
			// That way the cmd is still responsive for potential feedback.
			var reader io.Reader
			if fsys.Exists(inf) {
				file, err := fsys.Open(inf)
				if err != nil {
					return err
				}
				defer file.Close()

				reader = file
			} else {
				reader = strings.NewReader(in)
			}

			// Scan everything!
			scanner := bufio.NewScanner(reader)
			scanner.Split(splitFn)

			var buf []string
			for scanner.Scan() {
				buf = append(buf, scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				return err
			}

			// Perform the sorting
			sort.Strings(buf)

			// Create a buffer so that writing to sources becomes more natural
			out := bytes.NewBufferString(strings.Join(buf, *separator))

			// Work out where to write to.
			var writer io.Writer
			if outf := strings.TrimSpace(*outputFile); outf != "" {
				file, err := fsys.Create(outf)
				defer file.Close()
				if err != nil {
					return err
				}

				writer = file
			} else {
				writer = os.Stdout
			}

			// Write the output
			_, err := out.WriteTo(writer)
			return err
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
