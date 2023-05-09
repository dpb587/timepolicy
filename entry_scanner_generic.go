package timepolicy

import (
	"bufio"
	"fmt"
)

type genericEntryScanner struct {
	s             *bufio.Scanner
	fieldSplitter EntryFieldSplitterFunc
	fieldsLimit   int
	timeField     int
	timeParser    TimeParserFunc

	err         error
	entry       *Entry
	entryOffset int
}

var _ EntryScanner = &genericEntryScanner{}

func NewGenericEntryScanner(s *bufio.Scanner, fieldSplitter EntryFieldSplitterFunc, fieldsLimit int, timeField int, timeParser TimeParserFunc) EntryScanner {
	return &genericEntryScanner{
		s:             s,
		fieldSplitter: fieldSplitter,
		fieldsLimit:   fieldsLimit,
		timeField:     timeField,
		timeParser:    timeParser,
	}
}

func (es *genericEntryScanner) Scan() bool {
	if es.err != nil {
		return false
	}

	scanned := es.s.Scan()
	if !scanned {
		es.err = es.s.Err()

		return false
	}

	es.entryOffset++

	raw := string(es.s.Bytes())
	fields, err := es.fieldSplitter(raw, es.fieldsLimit)
	if err != nil {
		es.err = fmt.Errorf("parsing entry %d: splitting fields: %v", es.entryOffset, err)

		return false
	}

	if len(fields)-1 < es.timeField {
		es.err = fmt.Errorf("parsing entry %d: parsing time: missing field %d", es.entryOffset, es.timeField)

		return false
	}

	timeParsed, err := es.timeParser(fields[es.timeField])
	if err != nil {
		es.err = fmt.Errorf("parsing entry %d: parsing time: %v", es.entryOffset, err)

		return false
	}

	es.entry = &Entry{
		Raw:    raw,
		Fields: fields,
		Time:   timeParsed,
	}

	return true
}

func (es *genericEntryScanner) Entry() *Entry {
	return es.entry
}

func (es *genericEntryScanner) EntryOffset() int {
	return es.entryOffset - 1
}

func (es *genericEntryScanner) Err() error {
	return es.err
}
