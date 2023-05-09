package internal

import (
	"fmt"

	"github.com/google/cel-go/cel"
)

var InputExpressionEnv *cel.Env

func init() {
	env, err := cel.NewEnv(
		cel.Variable("entry", cel.StringType),
		cel.Variable("ts", cel.TimestampType),
		cel.Variable("fields", cel.ListType(cel.StringType)),
	)
	if err != nil {
		panic(fmt.Errorf("evaluation expression env: %w", err))
	}

	InputExpressionEnv = env
}
