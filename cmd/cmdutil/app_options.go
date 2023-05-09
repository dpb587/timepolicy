package cmdutil

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/dpb587/timepolicy/internal/version"
)

type AppOptions struct {
	Quiet   bool        `name:"quiet" short:"q" help:"Suppress normal output."`
	Version versionFlag `name:"version" help:"Show version."`
}

type versionFlag bool

func (v versionFlag) AfterApply(app *kong.Kong) error {
	if v {
		fmt.Fprintf(
			app.Stdout,
			"%s %s build-commit=%s build-clean=%s build-time=%s\n",
			app.Model.Name,
			version.Name,
			version.BuildCommit,
			version.BuildClean,
			version.BuildTime,
		)

		return NewErrorWithExitCode(nil, 0)
	}

	return nil
}
