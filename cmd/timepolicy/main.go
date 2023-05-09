package main

import (
	"go/doc"
	"os"

	"github.com/alecthomas/kong"
	"github.com/dpb587/timepolicy/cmd/cmdutil"
	"github.com/dpb587/timepolicy/cmd/timepolicy/rootcmd"
)

const mainName = "timepolicy"

type mainOptions struct {
	*cmdutil.AppOptions `group:"Global Flags"`
	*rootcmd.Command    `cmd:""`
}

func main() {
	cli := mainOptions{
		AppOptions: &cmdutil.AppOptions{},
	}

	app, err := kong.New(
		&cli,
		kong.Name(mainName),
		kong.Description("Select entries that match policies based on time and metadata."),
		kong.ConfigureHelp(kong.HelpOptions{
			WrapUpperBound: 120,
		}),
		kong.Help(func(options kong.HelpOptions, ctx *kong.Context) error {
			err := kong.DefaultHelpPrinter(options, ctx)
			if err != nil {
				return err
			}

			ctx.Stdout.Write([]byte("\n"))

			doc.ToText(
				ctx.Stdout,
				`POLICY SPECIFICATIONS

Policies evaluate entries based on (1) being within a Time Range and (2) Optional Qualifiers. An entry is selected while one or more policies apply to it.

Time Range should be in the format of {INT}{unit}. Supported units are: s for seconds, h for hours, d for days, m for months, and y for years.

Optional Qualifiers may be zero or more of the following:

 - if={EXPR} - an expression that must be true for the entry to be considered (in addition to Time Range). Expressions must evaluate to true or false. See ADVANCED EXPRESSIONS for details.
 - by={EXPR} - a method to further segment matching entries. Simple values of year, month (year-month), day, (year-month-day), and hour (year-month-day-hour) are supported; and ADVANCED EXPRESSIONS may be used for complex strategies.
 - oldest or newest - whether the policy prefers older or newer entries. By default, newest entries are preferred.
 - max={INT} - limit the policy to apply to, at most, {INT} entries. By default, if no other qualifier is used, a policy applies to all entries (max=-1); otherwise it defaults to a single entry (max=1).
`, "", "    ", 120)

			ctx.Stdout.Write([]byte("\n"))

			doc.ToText(
				ctx.Stdout,
				`ADVANCED EXPRESSIONS

For complex situations, use Common Expression Language (see https://github.com/google/cel-spec) for supported values. The following context is available:

 - entry - raw input (string)
 - ts - parsed timestamp from the entry
 - fields - parsed field list from the entry (strings)
`, "", "    ", 120)

			ctx.Stdout.Write([]byte("\n"))

			doc.ToText(
				ctx.Stdout,
				`EXIT STATUS

 - 0 - one or more entries were selected
 - 1 - no entries were selected
 - >1 - a general error occurred
`, "", "    ", 120)

			return nil
		}),
	)
	if err != nil {
		cmdutil.Exit(err, mainName, os.Stderr, false)
	}

	appctx, err := app.Parse(os.Args[1:])
	if err != nil {
		cmdutil.Exit(err, mainName, app.Stderr, cli.AppOptions.Quiet)
	}

	cmdutil.Exit(appctx.Run(cli.AppOptions), mainName, app.Stderr, cli.AppOptions.Quiet)
}
