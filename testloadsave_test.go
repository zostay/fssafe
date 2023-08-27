package fssafe

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestLoaderSaver(t *testing.T) {
	t.Parallel()

	// setup
	k := NewTestingLoaderSaver()

	// testing the loader, before first saver
	r, err := k.Loader()
	assert.Error(t, err, "loader should fail before saver")
	assert.Nil(t, r, "loader should fail before saver")

	expectBuffers := make([][]byte, 0, 3)
	for i := 1; i <= 3; i++ {
		// testing the saver
		w, err := k.Saver()
		require.NoError(t, err, "save creates file")

		assert.Len(t, k.Buffers(), i-1, "no buffer yet")

		s := RandString(33)
		expectBuffers = append(expectBuffers, []byte(s))
		_, _ = io.WriteString(w, s)

		assert.Len(t, k.Buffers(), i-1, "still not buffer yet")

		// saver finalizes the buffer for testing
		err = w.Close()
		require.NoError(t, err, "close should create file")

		assert.Len(t, k.Buffers(), i, "got a buffer now")

		// testing the loader
		r, err := k.Loader()
		require.NoError(t, err, "buffer exists, so reading it should be fine")
		require.NotNil(t, r, "buffer exists, so we get a reader")

		sr, err := io.ReadAll(r)
		require.NoError(t, err, "reading file should not have an err")

		assert.Equal(t, []byte(s), sr, "found the same string from loader that we wrote")
	}

	assert.Len(t, k.Buffers(), 3, "got three buffers at the end")
	for i, buf := range k.Buffers() {
		assert.Equal(t, expectBuffers[i], buf.Bytes(), "got the expected buffers")
	}
}
