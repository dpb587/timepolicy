package timepolicy

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dpb587/timepolicy/internal"
	"github.com/google/cel-go/cel"
)

var (
	policySpecQualifierSplit = regexp.MustCompile(`\s*;\s*`)
	policySpecBucketEnums    = map[string]func(e *Entry) (string, error){
		"year": func(e *Entry) (string, error) {
			return e.Time.Format("2006"), nil
		},
		"month": func(e *Entry) (string, error) {
			return e.Time.Format("2006-01"), nil
		},
		"day": func(e *Entry) (string, error) {
			return e.Time.Format("2006-01-02"), nil
		},
		"hour": func(e *Entry) (string, error) {
			return e.Time.Format("2006-01-02T15"), nil
		},
	}
	policySpecRangeUnits = map[byte]func(n int64) time.Time{
		's': func(n int64) time.Time {
			return Now().Add(-1 * time.Second * time.Duration(n))
		},
		'h': func(n int64) time.Time {
			return Now().Add(-1 * time.Hour * time.Duration(n))
		},
		'd': func(n int64) time.Time {
			return Now().AddDate(0, 0, int(-1*n))
		},
		'm': func(n int64) time.Time {
			return Now().AddDate(0, int(-1*n), 0)
		},
		'y': func(n int64) time.Time {
			return Now().AddDate(int(-1*n), 0, 0)
		},
	}
)

var commentSplit = regexp.MustCompile(`\s+//\s+`)

func ParsePolicySpecString(defaultName string, raw string) (*PolicySpec, error) {
	var err error

	ps := &PolicySpec{
		raw:  raw,
		name: defaultName,
		max:  -1,
	}

	rawSplit := commentSplit.Split(raw, 2)
	if len(rawSplit) > 1 {
		ps.comment = rawSplit[1]
	}

	rawStatements := policySpecQualifierSplit.Split(rawSplit[0], -1)
	uniqQualifiers := map[string]struct{}{}

	for rawPieceIdx, rawPiece := range rawStatements {
		rawPieceSplit := strings.SplitN(rawPiece, "=", 2)

		if rawPieceIdx == 0 {
			if len(rawPieceSplit) == 1 {
				ps.cutoff, err = parsePolicySpecRangeCutoff(rawPieceSplit[0])
				if err != nil {
					return nil, fmt.Errorf("parsing range: %v", err)
				}

				continue
			}

			// skipped range; applies to all
		}

		if _, known := uniqQualifiers[rawPieceSplit[0]]; known {
			return nil, fmt.Errorf("parsing qualifier: parsing %s: already configured", rawPieceSplit[0])
		}

		uniqQualifiers[rawPieceSplit[0]] = struct{}{}

		switch rawPieceSplit[0] {
		case "if":
			if len(rawPieceSplit) != 2 {
				return nil, errors.New("parsing qualifier: parsing if: missing value")
			}

			ast, issues := internal.InputExpressionEnv.Compile(rawPieceSplit[1])
			if err := issues.Err(); err != nil {
				return nil, fmt.Errorf("parsing qualifier: parsing if: compiling: %v", err)
			} else if !ast.IsChecked() || ast.OutputType() != cel.BoolType {
				return nil, errors.New("parsing qualifier: parsing if: expression must have boolean result")
			}

			prg, err := internal.InputExpressionEnv.Program(ast)
			if err != nil {
				return nil, fmt.Errorf("parsing qualifier: parsing if: installing: %v", err)
			}

			ps.condition = prg
		case "by":
			if len(rawPieceSplit) != 2 {
				return nil, errors.New("parsing qualifier: parsing by: missing value")
			}

			enumFunc, ok := policySpecBucketEnums[rawPieceSplit[1]]
			if ok {
				ps.bucket = enumFunc

				continue
			}

			ast, issues := internal.InputExpressionEnv.Compile(rawPieceSplit[1])
			if err := issues.Err(); err != nil {
				return nil, fmt.Errorf("parsing qualifier: parsing by: compiling: %v", err)
			} else if !ast.IsChecked() {
				return nil, errors.New("parsing qualifier: parsing by: expression must have a deterministic result")
			}

			prg, err := internal.InputExpressionEnv.Program(ast)
			if err != nil {
				return nil, fmt.Errorf("parsing qualifier: parsing by: installing: %v", err)
			}

			switch ast.OutputType() {
			case cel.StringType:
				ps.bucket = func(e *Entry) (string, error) {
					val, _, err := e.Eval(prg)
					if err != nil {
						return "", err
					}

					return val.Value().(string), nil
				}
			case cel.BoolType:
				ps.bucket = func(e *Entry) (string, error) {
					val, _, err := e.Eval(prg)
					if err != nil {
						return "", err
					} else if val.Value().(bool) {
						return "true", nil
					}

					return "false", nil
				}
			case cel.IntType:
				ps.bucket = func(e *Entry) (string, error) {
					val, _, err := e.Eval(prg)
					if err != nil {
						return "", err
					}

					switch valT := val.Value().(type) {
					case int:
						return strconv.FormatInt(int64(valT), 10), nil
					case int8:
						return strconv.FormatInt(int64(valT), 10), nil
					case int16:
						return strconv.FormatInt(int64(valT), 10), nil
					case int32:
						return strconv.FormatInt(int64(valT), 10), nil
					case int64:
						return strconv.FormatInt(int64(valT), 10), nil
					}

					return "", fmt.Errorf("invalid int data type: %T", val)
				}
			default:
				return nil, fmt.Errorf("parsing qualifier: expression must have a string, bool, or integer result")
			}

			continue
		case "oldest", "newest":
			if len(rawPieceSplit) == 2 {
				return nil, fmt.Errorf("parsing qualifier: parsing %s: unexpected value", rawPieceSplit[0])
			}

			if rawPieceSplit[0] == "oldest" {
				if _, known := uniqQualifiers["newest"]; known {
					return nil, fmt.Errorf("parsing qualifier: parsing %s: either oldest or newest may be configured only once", rawPieceSplit[0])
				}

				ps.oldest = true
			} else {
				if _, known := uniqQualifiers["oldest"]; known {
					return nil, fmt.Errorf("parsing qualifier: parsing %s: either oldest or newest may be configured only once", rawPieceSplit[0])
				}
			}

			continue
		case "max":
			if len(rawPieceSplit) != 2 {
				return nil, errors.New("parsing qualifier: parsing max: missing value")
			}

			maxInt, err := strconv.ParseInt(rawPieceSplit[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parsing qualifier: parsing max: parsing number: %v", err)
			} else if maxInt == 0 || maxInt < -1 {
				return nil, fmt.Errorf("parsing qualifier: parsing max: number must be -1 or greater than 0")
			}

			ps.max = int(maxInt)

			continue
		}

		return nil, fmt.Errorf("parsing qualifier: unexpected input: %s", rawPiece)
	}

	if _, known := uniqQualifiers["max"]; !known {
		if ps.bucket != nil {
			ps.max = 1
		}
	}

	return ps, nil
}

func parsePolicySpecRangeCutoff(value string) (time.Time, error) {
	valueLen := len(value)
	if valueLen < 2 {
		return time.Time{}, errors.New("invalid value: expected `{INT}{UNIT}`")
	}

	valueUnitFunc, ok := policySpecRangeUnits[value[valueLen-1]]
	if !ok {
		return time.Time{}, fmt.Errorf("parsing unit: %s", string(value[valueLen-1]))
	}

	valueNumber, err := strconv.ParseInt(value[0:(valueLen-1)], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing number: %v", err)
	}

	return valueUnitFunc(valueNumber), nil
}
