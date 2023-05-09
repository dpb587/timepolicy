package cmdutil

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func Exit(err error, name string, stderr io.Writer, quiet bool) {
	var exitCode = 0

	defer func() {
		if exitCode > 0 && !quiet {
			fmt.Fprintf(stderr, "%s: exit status %d\n", name, exitCode)
		}

		os.Exit(exitCode)
	}()

	if err == nil {
		return
	}

	exitCode = 2

	var ee ErrorWithExitCode

	if errors.As(err, &ee) {
		exitCode = ee.ErrorExitCode()

		if exitCode == 0 {
			return
		}
	}

	if !quiet {
		fmt.Fprintf(stderr, "%s: error: %s\n", name, err.Error())
	}
}
