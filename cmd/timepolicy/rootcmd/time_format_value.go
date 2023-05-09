package rootcmd

import (
	"strconv"
	"time"

	"github.com/alecthomas/kong"
	"github.com/dpb587/timepolicy"
)

func timeFormatValueLayout(layout string) timepolicy.TimeParserFunc {
	return func(v string) (time.Time, error) {
		return time.Parse(layout, v)
	}
}

var timeFormatValueEnums = map[string]timepolicy.TimeParserFunc{
	"ANSIC":       timeFormatValueLayout(time.ANSIC),
	"UnixDate":    timeFormatValueLayout(time.UnixDate),
	"RubyDate":    timeFormatValueLayout(time.RubyDate),
	"RFC822":      timeFormatValueLayout(time.RFC822),
	"RFC822Z":     timeFormatValueLayout(time.RFC822Z),
	"RFC850":      timeFormatValueLayout(time.RFC850),
	"RFC1123":     timeFormatValueLayout(time.RFC1123),
	"RFC1123Z":    timeFormatValueLayout(time.RFC1123Z),
	"RFC3339":     timeFormatValueLayout(time.RFC3339),
	"RFC3339Nano": timeFormatValueLayout(time.RFC3339Nano),
	"Stamp":       timeFormatValueLayout(time.Stamp),
	"StampMilli":  timeFormatValueLayout(time.StampMilli),
	"StampMicro":  timeFormatValueLayout(time.StampMicro),
	"StampNano":   timeFormatValueLayout(time.StampNano),

	//

	"Unix": func(v string) (time.Time, error) {
		vInt, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}, err
		}

		return time.Unix(vInt, 0), nil
	},
	"UnixMilli": func(v string) (time.Time, error) {
		vInt, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}, err
		}

		return time.Unix(vInt, 0), nil
	},

	//

	"YYYY-MM-DD": timeFormatValueLayout("2006-01-02"),
}

type TimeFormatValue struct {
	raw string
	f   timepolicy.TimeParserFunc
}

var _ kong.MapperValue = &TimeFormatValue{}

func (v *TimeFormatValue) Decode(ctx *kong.DecodeContext) error {
	var raw string

	err := ctx.Scan.PopValueInto("string", &raw)
	if err != nil {
		return err
	}

	enumFunc, ok := timeFormatValueEnums[raw]
	if ok {
		v.raw = raw
		v.f = enumFunc

		return nil
	}

	v.raw = raw
	v.f = timeFormatValueLayout(raw)

	return nil
}
