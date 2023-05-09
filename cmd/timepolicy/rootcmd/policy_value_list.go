package rootcmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/dpb587/timepolicy"
)

type PolicyValueList struct {
	values []*timepolicy.PolicySpec
}

var _ kong.MapperValue = &PolicyValueList{}

func (v *PolicyValueList) Decode(ctx *kong.DecodeContext) error {
	var raw string

	err := ctx.Scan.PopValueInto("string", &raw)
	if err != nil {
		return err
	}

	parsed, err := timepolicy.ParsePolicySpecString(fmt.Sprintf("policy-%d", len(v.values)), raw)
	if err != nil {
		return err
	}

	v.values = append(v.values, parsed)

	return nil
}
