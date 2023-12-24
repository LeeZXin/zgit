package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"unsafe"
)

type PipeResultCloser interface {
	ClosePipe()
}

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
	o      sync.Once
}

func (r *ReadPipeResult) Reader() io.Reader {
	return r.reader
}

func (r *ReadPipeResult) ClosePipe() {
	r.o.Do(func() {
		r.reader.Close()
	})
}

func (r *ReadPipeResult) RangeStringLines(rangeFn func(int, string) (bool, error)) error {
	defer r.ClosePipe()
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
		shouldContinue, err := rangeFn(i, strings.TrimSpace(string(line)))
		if err != nil {
			return err
		}
		if !shouldContinue {
			return nil
		}
	}
}

type ReadWritePipeResult struct {
	reader *io.PipeReader
	writer *io.PipeWriter
	o      sync.Once
}

func (r *ReadWritePipeResult) Reader() io.Reader {
	return r.reader
}

func (r *ReadWritePipeResult) Writer() io.Writer {
	return r.writer
}

func (r *ReadWritePipeResult) ClosePipe() {
	r.o.Do(func() {
		r.reader.Close()
		r.writer.Close()
	})
}

func (r *ReadWritePipeResult) RangeStringLines(rangeFn func(int, string) (bool, error)) error {
	defer r.ClosePipe()
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
		shouldContinue, err := rangeFn(i, string(line))
		if err != nil {
			return err
		}
		if !shouldContinue {
			return nil
		}
	}
}
