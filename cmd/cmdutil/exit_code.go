package cmdutil

type ErrorWithExitCode interface {
	error
	ErrorExitCode() int
}

//

type errorWithExitCode struct {
	error
	exitCode int
}

var _ error = errorWithExitCode{}

func (err errorWithExitCode) Unwrap() error {
	return err.error
}

func (err errorWithExitCode) Error() string {
	return err.error.Error()
}

func (err errorWithExitCode) ErrorExitCode() int {
	return err.exitCode
}

//

func NewErrorWithExitCode(err error, exitCode int) ErrorWithExitCode {
	return errorWithExitCode{
		error:    err,
		exitCode: exitCode,
	}
}
