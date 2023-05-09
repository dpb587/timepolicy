package rootcmd

import (
	"errors"

	"github.com/dpb587/timepolicy/cmd/cmdutil"
)

var ErrNoEntries error = cmdutil.NewErrorWithExitCode(errors.New("no entries selected"), 1)
