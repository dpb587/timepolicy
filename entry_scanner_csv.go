package timepolicy

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

type csvEntryScanner struct {
	r          *csv.Reader
	timeField  int
	timeParser TimeParserFunc

	err         error
	entry       *Entry
	entryOffset int
}

var _ EntryScanner = &csvEntryScanner{}

func NewCSVEntryScanner(r *csv.Reader, timeField int, timeParser TimeParserFunc) EntryScanner {
	return &csvEntryScanner{
		r:          r,
		timeField:  timeField,
		timeParser: timeParser,
	}
}

func (es *csvEntryScanner) Scan() bool {
	if es.err != nil {
		return false
	}

	fields, err := es.r.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false
		}

		es.err = err

		return false
	}

	es.entryOffset++

	if len(fields)-1 < es.timeField {
		es.err = fmt.Errorf("parsing entry %d: parsing time: missing field %d", es.entryOffset, es.timeField)

		return false
	}

	timeParsed, err := es.timeParser(fields[es.timeField])
	if err != nil {
		es.err = fmt.Errorf("parsing entry %d: parsing time: %v", es.entryOffset, err)

		return false
	}

	// hacky
	raw := bytes.NewBuffer(nil)
	w := csv.NewWriter(raw)
	w.Write(fields)
	w.Flush()

	es.entry = &Entry{
		Raw:    strings.TrimSuffix(raw.String(), "\n"),
		Fields: fields,
		Time:   timeParsed,
	}

	return true
}

func (es *csvEntryScanner) Entry() *Entry {
	return es.entry
}

func (es *csvEntryScanner) EntryOffset() int {
	return es.entryOffset - 1
}

func (es *csvEntryScanner) Err() error {
	return es.err
}
