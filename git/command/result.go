package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unsafe"
)

type Result struct {
	stdOut *bytes.Buffer
}

func (r *Result) ReadAsBytes() []byte {
	return r.stdOut.Bytes()
}

func (r *Result) ReadAsString() string {
	return bytesToString(r.ReadAsBytes())
}

func stdErrorResult(err error, stdErr string) error {
	return fmt.Errorf("%w - %s", err, stdErr)
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

type ReadPipeResult struct {
	reader *io.PipeReader
}

func (r *ReadPipeResult) Reader() io.Reader {
	return r.reader
}

func (r *ReadPipeResult) ClosePipe() {
	r.reader.Close()
}

func (r *ReadPipeResult) RangeStringLines(rangeFn func(int, string) error) error {
	reader := bufio.NewReader(r.reader)
	for i := 0; ; i++ {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if isPrefix {
			continue
		}
		err = rangeFn(i, string(line))
		if err != nil {
			return err
		}
	}
}

type ReadWritePipeResult struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func (r *ReadWritePipeResult) Reader() io.Reader {
	return r.reader
}

func (r *ReadWritePipeResult) Writer() io.Writer {
	return r.writer
}

func (r *ReadWritePipeResult) ClosePipe() {
	r.reader.Close()
	r.writer.Close()
}
