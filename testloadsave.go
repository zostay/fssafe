package fssafe

import (
	"bytes"
	"io"
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
}

// Close marks the writer as closed.
func (t *TestingWriter) Close() error { t.Closed = true; return nil }

// TestingLoaderSaver provides an in-memory loader and saver for testing.
type TestingLoaderSaver struct {
	BasicLoaderSaver
	Readers []*TestingReader
	Writers []*TestingWriter
}

// NewTestingLoaderSaver returns a new, blank TestingLoaderSaver.
func NewTestingLoaderSaver() *TestingLoaderSaver {
	rs := make([]*TestingReader, 0)
	ws := make([]*TestingWriter, 0)

	var buf *bytes.Buffer
	loader := func() (io.ReadCloser, error) {
		r := &TestingReader{bytes.NewReader(buf.Bytes()), false}
		rs = append(rs, r)
		return r, nil
	}

	saver := func() (io.WriteCloser, error) {
		buf = new(bytes.Buffer)
		w := &TestingWriter{buf, false}
		ws = append(ws, w)
		return w, nil
	}

	ls := TestingLoaderSaver{BasicLoaderSaver{loader, saver}, rs, ws}

	return &ls
}
