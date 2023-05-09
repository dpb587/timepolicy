package timepolicy

import (
	"fmt"
	"time"

	"github.com/google/cel-go/cel"
)

type PolicySpec struct {
	raw string

	cutoff    time.Time
	condition cel.Program
	bucket    func(e *Entry) (string, error)
	oldest    bool
	max       int

	name    string
	comment string
}

func (ps *PolicySpec) String() string {
	return ps.raw
}

func (ps *PolicySpec) MatchEntry(e *Entry) (bool, error) {
	if e.Time.Before(ps.cutoff) {
		return false, nil
	}

	if ps.condition != nil {
		val, _, err := e.Eval(ps.condition)
		if err != nil {
			return false, fmt.Errorf("evaluating condition: %v", err)
		} else if !val.Value().(bool) {
			return false, nil
		}
	}

	return true, nil
}
