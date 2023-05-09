package timepolicy

import (
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
)

type Entry struct {
	Raw    string
	Fields []string
	Time   time.Time
}

func (e *Entry) Eval(prg cel.Program) (ref.Val, *cel.EvalDetails, error) {
	return prg.Eval(map[string]interface{}{
		"entry":  e.Raw,
		"ts":     e.Time,
		"fields": e.Fields,
	})
}
