package timepolicy

import "io"

type EntryWriter interface {
	EntriesWritten() int64
	WriteEntry(e *Entry) error
}

//

type rawEntryWriter struct {
	c int64
	w io.Writer
}

func NewEntryWriter(w io.Writer) EntryWriter {
	return &rawEntryWriter{
		w: w,
	}
}

func (w *rawEntryWriter) EntriesWritten() int64 {
	return w.c
}

func (w *rawEntryWriter) WriteEntry(e *Entry) error {
	w.c++
	_, err := w.w.Write([]byte(e.Raw + "\n"))

	return err
}

//

type builtinEntryFieldWriter struct {
	c     int64
	w     io.Writer
	field int
}

func NewEntryFieldWriter(w io.Writer, field int) EntryWriter {
	return &builtinEntryFieldWriter{
		w:     w,
		field: field,
	}
}

func (w *builtinEntryFieldWriter) EntriesWritten() int64 {
	return w.c
}

func (w *builtinEntryFieldWriter) WriteEntry(e *Entry) error {
	var err error

	w.c++

	if len(e.Fields) < w.field {
		_, err = w.w.Write([]byte("\n"))
	} else {
		_, err = w.w.Write([]byte(e.Fields[w.field] + "\n"))
	}

	return err
}

//

type discardEntryWriter struct {
	c int64
}

func NewDiscardEntryWriter() EntryWriter {
	return &discardEntryWriter{}
}

func (w *discardEntryWriter) EntriesWritten() int64 {
	return w.c
}

func (w *discardEntryWriter) WriteEntry(e *Entry) error {
	w.c++

	return nil
}
