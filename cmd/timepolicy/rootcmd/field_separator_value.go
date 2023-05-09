package rootcmd

import (
	"encoding/csv"
	"regexp"

	"github.com/alecthomas/kong"
	"github.com/dpb587/timepolicy"
)

type FieldSeparatorValue struct {
	csvApplier func(r *csv.Reader)
	f          timepolicy.EntryFieldSplitterFunc
}

var _ kong.MapperValue = &FieldSeparatorValue{}

func (v *FieldSeparatorValue) Decode(ctx *kong.DecodeContext) error {
	var raw string

	err := ctx.Scan.PopValueInto("string", &raw)
	if err != nil {
		return err
	}

	if raw == "spaces" {
		v.f = timepolicy.SpacesEntryFieldSplitter

		return nil
	} else if raw == "csv" {
		v.csvApplier = func(r *csv.Reader) {
			r.Comma = ','
		}

		return nil
	} else if raw == "tsv" {
		v.csvApplier = func(r *csv.Reader) {
			r.Comma = '\t'
		}

		return nil
	}

	parsed, err := regexp.Compile(raw)
	if err != nil {
		return err
	}

	v.f = func(v string, limit int) ([]string, error) {
		return parsed.Split(v, limit), nil
	}

	return nil
}
