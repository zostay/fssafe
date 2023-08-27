package fssafe

import (
	"bytes"
	"io"
	"os"
)

// TestingReader provides an in-memory reader for testing.
type TestingReader struct {
	*bytes.Reader      // the in-memory reader
	Closed        bool // whether the reader has been closed
}

// Close marks the reader as closed.
func (t *TestingReader) Close() error { t.Closed = true; return nil }

// TestingWriter provides an in-memory writer for testing.
type TestingWriter struct {
	*bytes.Buffer      // the in-memory writer
	Closed        bool // whether the writer has been closed
	final         **bytes.Buffer
	bufferList    *[]*bytes.Buffer
}

// Close marks the writer as closed.
func (t *TestingWriter) Close() error {
	*t.final = t.Buffer
	*t.bufferList = append(*t.bufferList, t.Buffer)
	t.Closed = true
	return nil
}

// TestingLoaderSaver provides an in-memory loader and saver for testing.
type TestingLoaderSaver struct {
	BasicLoaderSaver
	bufferList *[]*bytes.Buffer
}

func (t *TestingLoaderSaver) Buffers() []*bytes.Buffer {
	return *t.bufferList
}

// NewTestingLoaderSaver returns a new, blank TestingLoaderSaver.
func NewTestingLoaderSaver() *TestingLoaderSaver {
	var buf *bytes.Buffer
	loader := func() (io.ReadCloser, error) {
		if buf == nil {
			return nil, os.ErrNotExist
		}

		r := &TestingReader{bytes.NewReader(buf.Bytes()), false}
		return r, nil
	}

	bufferList := []*bytes.Buffer{}
	saver := func() (io.WriteCloser, error) {
		inProgress := &bytes.Buffer{}
		w := &TestingWriter{inProgress, false, &buf, &bufferList}
		return w, nil
	}

	ls := TestingLoaderSaver{BasicLoaderSaver{loader, saver}, &bufferList}

	return &ls
}
