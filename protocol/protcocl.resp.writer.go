package protocol

import (
	"fmt"
	"io"
)

// Commands to RESP: OK -> +OK\r\n

type Writer struct{
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer{
	return &Writer{
		writer: w,
	}
}

func (w *Writer) WriterSimpleString(value string) error {
	_, err := fmt.Fprintf(w.writer, "+%s\r\n", value)
	return err
}

func (w *Writer) WriteError(value string) error {
	_, err := fmt.Fprintf(w.writer, "-%s\r\n", value)
	return err
}

func (w *Writer) WriteInteger(value int64) error {
	_, err := fmt.Fprintf(w.writer, ":%d\r\n", value)
	return err
}

func (w *Writer) WriteBulkString(value string) error {
	_, err := fmt.Fprintf(w.writer, "$%d\r\n%s\r\n", len(value), value)
	return err
}

func (w *Writer) WriteNull() error {
	_, err := fmt.Fprintf(w.writer, "$-1\r\n")
	return err
}

