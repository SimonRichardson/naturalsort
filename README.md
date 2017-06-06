# Natural Sort

## Command Natural

 - [Getting started](#getting-started)
 - [Introduction](#introduction)
 - [Sort](#sort)
 - [Tests](#tests)

### Getting started

The natural sort command expects to have some things pre-installed via `go get`,
so that you can build the project.

 - `go get github.com/Masterminds/glide`

-----

Quick guide to getting started, this assumes you've got the `$GOPATH` setup
correctly and the gopath bin folder is in your `$PATH`:

```
glide install
make clean all
cd dist

./natural sort -input="0001,02,2,1,4,99,00,0,3,1"
```

The following should output:

```
0,00,1,1,0001,2,02,3,4,99
```

### Introduction

The natural sort CLI is broken down into one distinctive command `sort`.

### Sort

The `sort` command essentially takes an input and performs a natural sort on it
and outputs it. The command has a multitude of options for reading files which
can optionally encoded/decoded in gzip/base64 formats.

```
natural sort -input="your input here"
```

Also available is a comprehensive `-help` section:

```
/natural sort -help
USAGE
  sort [flags]

FLAGS
  -debug false          debug logging
  -input                input for natural sorting
  -input.base64 false   decode base 64 input
  -input.file           file required to perform natural sorting on
  -input.gzip false     decode gzip input
  -output.base64 false  encode base64 output
  -output.file          output file for action performed
  -output.gzip false    encode gzip output
  -separator ,          separation value
```

### Tests

Tests can be run using the following command, it also includes a series of
benchmarking tests:

```
 go test -v -bench=. $(glide nv)
```
