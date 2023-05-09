package rootcmd

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/dpb587/timepolicy"
	"github.com/dpb587/timepolicy/internal"
)

type WriteFormatValue struct {
	builder func(w io.Writer) timepolicy.EntryWriter
}

var _ kong.MapperValue = &WriteFormatValue{}

func (v *WriteFormatValue) Decode(ctx *kong.DecodeContext) error {
	var raw string

	err := ctx.Scan.PopValueInto("string", &raw)
	if err != nil {
		return err
	}

	if internal.DollarFieldRegExp.MatchString(raw) {
		fieldColumn, err := strconv.ParseInt(strings.TrimPrefix(raw, "$"), 10, 64)
		if err != nil {
			return fmt.Errorf("parsing field number: %v", err)
		} else if fieldColumn < 1 {
			return errors.New("expected field number to be greater than 0")
		}

		v.builder = func(w io.Writer) timepolicy.EntryWriter {
			return timepolicy.NewEntryFieldWriter(w, int(fieldColumn-1))
		}

		return nil
	}

	return errors.New("unsupported value")
}
