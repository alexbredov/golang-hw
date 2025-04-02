package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "target.txt", 10000, 1)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})
	t.Run("no original file", func(t *testing.T) {
		err := Copy("testdata/no_input.txt", "target.txt", 0, 0)
		require.Truef(t, errors.Is(err, ErrFileOpen), "actual error %q", err)
	})
	t.Run("creation failure", func(t *testing.T) {
		err := Copy("testdata/input.txt", "target.jpg", 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})
	t.Run("negative offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "target.txt", -1, 0)
		require.Truef(t, errors.Is(err, ErrNegativeOffset), "actual error %q", err)
	})
	t.Run("negative limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", "target.txt", 0, -1)
		require.Truef(t, errors.Is(err, ErrNegativeLimit), "actual error %q", err)
	})
}
