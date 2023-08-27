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
	readers    *[]*TestingReader
	writers    *[]*TestingWriter
	bufferList *[]*bytes.Buffer
}

// Buffers returns each buffer written by the saver in the order it occurred.
func (t *TestingLoaderSaver) Buffers() []*bytes.Buffer {
	return *t.bufferList
}

// CurrentReader returns the index of the current reader and can be used as the
// key into ReadersClosed to determine which was or was not closed.
func (t *TestingLoaderSaver) CurrentReader() int {
	return len(*t.readers) - 1
}

// ReadersClosed returns a slice of booleans indicating whether each reader
// created was closed.
func (t *TestingLoaderSaver) ReadersClosed() []bool {
	rs := make([]bool, len(*t.readers))
	for i, r := range *t.readers {
		rs[i] = r.Closed
	}
	return rs
}

// CurrentWriter returns the index of the current writer and can be used as the
// key into WritersClosed to determine which was or was not closed.
func (t *TestingLoaderSaver) CurrentWriter() int {
	return len(*t.writers) - 1
}

// WritersClosed returns a slice of booleans indicating whether each writer
// created was closed.
func (t *TestingLoaderSaver) WritersClosed() []bool {
	ws := make([]bool, len(*t.writers))
	for i, w := range *t.writers {
		ws[i] = w.Closed
	}
	return ws
}

// NewTestingLoaderSaver returns a new, blank TestingLoaderSaver.
func NewTestingLoaderSaver() *TestingLoaderSaver {
	rs := make([]*TestingReader, 0)
	ws := make([]*TestingWriter, 0)

	var buf *bytes.Buffer
	loader := func() (io.ReadCloser, error) {
		if buf == nil {
			return nil, os.ErrNotExist
		}

		r := &TestingReader{bytes.NewReader(buf.Bytes()), false}
		rs = append(rs, r)
		return r, nil
	}

	bufferList := make([]*bytes.Buffer, 0)
	saver := func() (io.WriteCloser, error) {
		inProgress := &bytes.Buffer{}
		w := &TestingWriter{inProgress, false, &buf, &bufferList}
		ws = append(ws, w)
		return w, nil
	}

	ls := TestingLoaderSaver{BasicLoaderSaver{loader, saver}, &rs, &ws, &bufferList}

	return &ls
}
