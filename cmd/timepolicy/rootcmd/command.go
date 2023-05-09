package rootcmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/dpb587/timepolicy"
	"github.com/dpb587/timepolicy/cmd/cmdutil"
)

type Command struct {
	FieldCount     int                  `name:"field-count" help:"Limit the number of fields extracted per entry."`
	FieldSeparator *FieldSeparatorValue `name:"field-separator" short:"F" placeholder:"STRING" help:"Separator used between fields. Value should be a regular expression (see https://pkg.go.dev/regexp/syntax) or a supported alias (csv, spaces, tsv). Default is spaces."`
	Read           io.Reader            `name:"read-from" short:"i" placeholder:"PATH" type:"path" help:"Read entries from file or path. Default is stdin."`
	Write          io.Writer            `name:"write-to" short:"o" placeholder:"PATH" type:"path" help:"Write selected entries to file or path. Default is stdout."`
	WriteFormat    *WriteFormatValue    `name:"write" placeholder:"FORMAT" help:"Write selected entries in a custom format. A single field may be output using dollar + field number, such as $1 for the first field."`
	Policies       PolicyValueList      `name:"policy" short:"p" placeholder:"STRING..." help:"One or more policies to evaluate entries against. See POLICY SPECIFICATIONS."`
	Invert         bool                 `name:"invert" help:"Show entries which are not covered by any policy. Enables streaming mode and entries may be written in a different order than they were read."`
	TimeFormat     *TimeFormatValue     `name:"time" placeholder:"STRING" help:"Format used by the time field. Value should be a custom layout (see https://pkg.go.dev/time#Layout) or a supported alias (ANSIC, UnixDate, RubyDate, RFC822, RFC822Z, RFC850, RFC1123, RFC1123Z, RFC3339, RFC3339Nano, Stamp, StampMilli, StampMicro, StampNano, Unix, UnixMilli, and YYYY-MM-DD). Default is RFC3339."`
	TimeField      int                  `name:"time-field" placeholder:"INT" help:"Field number containing the time, such as 1 for the first field."`
}

func (cmd *Command) BeforeApply() error {
	cmd.FieldSeparator = &FieldSeparatorValue{
		f: timepolicy.SpacesEntryFieldSplitter,
	}
	cmd.Read = os.Stdin
	cmd.Write = os.Stdout
	cmd.WriteFormat = &WriteFormatValue{
		builder: func(w io.Writer) timepolicy.EntryWriter {
			return timepolicy.NewEntryWriter(w)
		},
	}
	cmd.TimeFormat = &TimeFormatValue{
		f: timeFormatValueEnums["RFC3339"],
	}

	return nil
}

func (cmd *Command) Run(app *kong.Kong, appOptions *cmdutil.AppOptions) error {
	var input timepolicy.EntryScanner

	var fieldCount = -1
	if cmd.FieldCount > 0 {
		fieldCount = cmd.FieldCount
	}

	if cmd.FieldSeparator.csvApplier != nil {
		r := csv.NewReader(cmd.Read)
		cmd.FieldSeparator.csvApplier(r)
		r.FieldsPerRecord = fieldCount

		input = timepolicy.NewCSVEntryScanner(
			r,
			cmd.TimeField,
			cmd.TimeFormat.f,
		)
	} else {
		s := bufio.NewScanner(cmd.Read)

		input = timepolicy.NewGenericEntryScanner(
			s,
			cmd.FieldSeparator.f,
			fieldCount,
			cmd.TimeField,
			cmd.TimeFormat.f,
		)
	}

	selectedWriter := cmd.WriteFormat.builder(cmd.Write)
	evictedWriter := timepolicy.NewDiscardEntryWriter()

	if appOptions.Quiet {
		selectedWriter = timepolicy.NewDiscardEntryWriter()
		evictedWriter = timepolicy.NewDiscardEntryWriter()
	}

	if cmd.Invert {
		evictedWriter, selectedWriter = selectedWriter, evictedWriter
	}

	//

	policySelections := timepolicy.NewPolicySelectionSet(cmd.Policies.values, evictedWriter)

	//

	for input.Scan() {
		var entry = input.Entry()

		entrySelected, err := policySelections.EvaluateEntry(entry)
		if err != nil {
			return fmt.Errorf("processing entry %d: %v", input.EntryOffset()+1, err)
		} else if !entrySelected {
			evictedWriter.WriteEntry(entry)
		}
	}

	if err := input.Err(); err != nil {
		return err
	}

	//

	if cmd.Invert {
		if evictedWriter.EntriesWritten() == 0 {
			return ErrNoEntries
		}
	} else {
		for _, entry := range policySelections.Entries() {
			err := selectedWriter.WriteEntry(entry)
			if err != nil {
				return fmt.Errorf("writing selection: %v", err)
			}
		}

		if selectedWriter.EntriesWritten() == 0 {
			return ErrNoEntries
		}
	}

	return nil
}
