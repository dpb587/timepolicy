package timepolicy

import (
	"regexp"
	"time"
)

type EntryScanner interface {
	Scan() bool
	Err() error
	Entry() *Entry
	EntryOffset() int
}

type EntryFieldSplitterFunc func(v string, limit int) ([]string, error)

//

var reEntryFieldSplitterSpaces = regexp.MustCompile(`[[:space:]]+`)

var SpacesEntryFieldSplitter = EntryFieldSplitterFunc(func(v string, limit int) ([]string, error) {
	return reEntryFieldSplitterSpaces.Split(v, limit), nil
})

//

type TimeParserFunc func(v string) (time.Time, error)

//

//
