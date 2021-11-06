package cmdutil

import (
	"errors"
	"fmt"
	"os"
)

// Common cmdutil errors
var (
	ErrFilenameNotFound error = errors.New("Filename not found")
)

// ReadFileArg checks the arguments passed to via the cli and reads the filename provided
// exitOnArg boolean indicates whether to exit if the filename is provided. If false and
// the filename has not been provided, it will return an error instead
func ReadFileArg(exitOnArg bool) ([]byte, error) {
	argsEnough := len(os.Args) > 1
	if !argsEnough && exitOnArg {
		fmt.Println("Expected filename")
		os.Exit(1)
	} else if !argsEnough {
		return nil, ErrFilenameNotFound
	}

	filename := os.Args[1]
	return os.ReadFile(filename)
}

// ExitOnError will exit the application if the provided error is not nil
func ExitOnError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(2)
	}
}
