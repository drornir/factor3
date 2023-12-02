package cmd

import (
	"fmt"
	"os"
)

// EnvVars captures values set by `go:generate`, and others
/*
	$GOARCH
		The execution architecture (arm, amd64, etc.)
	$GOOS
		The execution operating system (linux, windows, etc.)
	$GOFILE
		The base name of the file.
	$GOLINE
		The line number of the directive in the source file.
	$GOPACKAGE
		The name of the package of the file containing the directive.
	$GOROOT
		The GOROOT directory for the 'go' command that invoked the
		generator, containing the Go toolchain and standard library.
	$DOLLAR
		A dollar sign.
	$PATH
		The $PATH of the parent process, with $GOROOT/bin
		placed at the beginning. This causes generators
		that execute 'go' commands to use the same 'go'
		as the parent 'go generate' command.
*/
type EnvVars struct {
	WorkingDirectory string // os.Getwd()

	GoArch    string // $GOARCH
	GoOS      string // $GOOS
	GoFile    string // $GOFILE
	GoLine    string // $GOLINE
	GoPackage string // $GOPACKAGE
	GoRoot    string // $GOROOT
	Dollar    string // $DOLLAR
	Path      string // $PATH
}

func LoadEnv() (EnvVars, error) {
	var e EnvVars

	var err error
	if e.WorkingDirectory, err = os.Getwd(); err != nil {
		return e, fmt.Errorf("os.Getwd() returned an error: %w", err)
	}

	e.GoArch = os.Getenv("GOARCH")
	e.GoOS = os.Getenv("GOOS")
	e.GoFile = os.Getenv("GOFILE")
	e.GoLine = os.Getenv("GOLINE")
	e.GoPackage = os.Getenv("GOPACKAGE")
	e.GoRoot = os.Getenv("GOROOT")
	e.Dollar = os.Getenv("DOLLAR")
	e.Path = os.Getenv("PATH")

	return e, nil
}
