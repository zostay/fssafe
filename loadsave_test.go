package fssafe

import (
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))] //nolint:gosec // test does not need secure random
	}
	return string(b)
}

func TestLoaderSaver(t *testing.T) {
	t.Parallel()

	// get a tempfile to work with
	tmpfile, err := os.CreateTemp(os.TempDir(), "kbdx")
	require.NoError(t, err, "able to get a tempfile")

	// cleanup tooling
	var hasfile, hasnew, hasold bool
	hasfile = true
	fn := tmpfile.Name()
	defer func() {
		if hasfile {
			_ = os.Remove(fn)
		}
		if hasnew {
			_ = os.Remove(fn + ".new")
		}
		if hasold {
			_ = os.Remove(fn + ".old")
		}
	}()

	// writing starter data
	s := RandString(20)
	_, _ = tmpfile.WriteString(s)
	err = tmpfile.Close()
	require.NoError(t, err, "closed initial tmpfile")

	// setup
	k := NewFileSystemLoaderSaver(fn)

	// testing the loader
	r, err := k.Loader()
	require.NoError(t, err, "tempfile already exists, so reading it should be fine")

	sr, err := io.ReadAll(r)
	require.NoError(t, err, "reading file should not have an err")

	assert.Equal(t, []byte(s), sr, "found the same string from loader that we wrote")

	fi, err := os.Stat(fn)
	require.NoError(t, err, "can stat the file")
	size := fi.Size()
	mtime := fi.ModTime()

	// testing the saver
	w, err := k.Saver()
	require.NoError(t, err, "save creates file")

	_, err = os.Stat(fn + ".new")
	require.NoError(t, err, "save created .new")

	hasnew = true

	fi, err = os.Stat(fn)
	require.NoError(t, err, "save preserved orig")

	// backup file is okay
	assert.Equal(t, size, fi.Size(), "orig size is same")
	assert.Equal(t, mtime, fi.ModTime(), "orig mtime is same")

	s = RandString(33)
	_, _ = io.WriteString(w, s)

	err = w.(*safeWriter).w.Sync()
	require.NoError(t, err, "sync worked")

	fi, err = os.Stat(fn + ".new")
	require.NoError(t, err, "stat new file while writing is ok")

	newsize := fi.Size()
	newmtime := fi.ModTime()

	// saver close does some renaming
	err = w.Close()
	require.NoError(t, err, "close should create file")

	_, err = os.Stat(fn + ".new")
	require.True(t, os.IsNotExist(err), ".new is gone")

	fi, err = os.Stat(fn + ".old")
	require.NoError(t, err, "orig is now .old")

	hasold = true

	assert.Equal(t, size, fi.Size(), ".old is same size as old orig")
	assert.Equal(t, mtime, fi.ModTime(), ".old is same mtime as old orig")

	fi, err = os.Stat(fn)
	require.NoError(t, err, ".new is now main")

	assert.Equal(t, newsize, fi.Size(), "size matches what was .new")
	assert.Equal(t, newmtime, fi.ModTime(), "mtime matches what was .new")

	err = os.Remove(fn)
	require.NoError(t, err, "failed to remove main")

	hasfile = false

	_, err = k.loader()
	require.Error(t, err, "tempfile was deleted, so loader should fail")
}

func TestLoaderSaverWithoutFile(t *testing.T) {
	t.Parallel()

	// get a tempfile to work with
	tmpfile, err := os.CreateTemp(os.TempDir(), "kbdx")
	require.NoError(t, err, "able to get a tempfile")

	err = os.Remove(tmpfile.Name())
	require.NoError(t, err)

	// cleanup tooling
	var hasfile, hasnew bool
	hasfile = true
	fn := tmpfile.Name()
	defer func() {
		if hasfile {
			_ = os.Remove(fn)
		}
		if hasnew {
			_ = os.Remove(fn + ".new")
		}
	}()

	// setup
	k := NewFileSystemLoaderSaver(fn)

	// testing the loader
	r, err := k.Loader()
	assert.Error(t, err, "tempfile does not exist")
	assert.Nil(t, r, "tempfile does not exist, so no reader")

	// testing the saver
	w, err := k.Saver()
	require.NoError(t, err, "save creates file")

	_, err = os.Stat(fn + ".new")
	require.NoError(t, err, "save created .new")

	hasnew = true

	fi, err := os.Stat(fn)
	require.ErrorIs(t, err, os.ErrNotExist, "final not exists yet")
	assert.Nil(t, fi, "no file info yet")

	s := RandString(33)
	_, _ = io.WriteString(w, s)

	err = w.(*safeWriter).w.Sync()
	require.NoError(t, err, "sync worked")

	fi, err = os.Stat(fn + ".new")
	require.NoError(t, err, "stat new file while writing is ok")

	newsize := fi.Size()
	newmtime := fi.ModTime()

	// saver close does some renaming
	err = w.Close()
	require.NoError(t, err, "close should finalize the file")

	hasnew = false

	_, err = os.Stat(fn + ".new")
	require.ErrorIs(t, err, os.ErrNotExist, ".new is gone")

	fi, err = os.Stat(fn + ".old")
	require.ErrorIs(t, err, os.ErrNotExist, ".old not created because there was no original")
	assert.Nil(t, fi, "no file info for .old")

	fi, err = os.Stat(fn)
	require.NoError(t, err, ".new is now main")

	assert.Equal(t, newsize, fi.Size(), "size matches what was .new")
	assert.Equal(t, newmtime, fi.ModTime(), "mtime matches what was .new")

	err = os.Remove(fn)
	require.NoError(t, err, "can remove file")

	hasfile = false

	_, err = k.loader()
	require.Error(t, err, "tempfile was deleted, so loader should fail")
}
