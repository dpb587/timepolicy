package internal

import "regexp"

var DollarFieldRegExp = regexp.MustCompile(`^\$\d+$`)
