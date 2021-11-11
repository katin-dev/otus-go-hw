package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	source = "testdata/input.txt"
	dest   = "testdata/dest.txt"
)

type ProgressBarStub struct{}

func (pb *ProgressBarStub) Init(total int) {
}

func (pb *ProgressBarStub) Increment() {
}

func (pb *ProgressBarStub) Finish() {
}

func TestCopy(t *testing.T) {
	pb := &ProgressBarStub{}

	t.Run("Empty Source File Name", func(t *testing.T) {
		err := Copy("", "", 0, 0, pb)
		require.Equal(t, ErrEmptySourceFilename, err)
	})

	t.Run("Empty Dest File Name", func(t *testing.T) {
		err := Copy(source, "", 0, 0, pb)
		require.Equal(t, ErrEmptyDestFilename, err)
	})

	t.Run("Source file is NOT exists", func(t *testing.T) {
		err := Copy("noname.txt", dest, 0, 0, pb)
		require.Equal(t, "failed to open source file noname.txt: open noname.txt: no such file or directory", err.Error())
	})

	t.Run("Offset exceed file size", func(t *testing.T) {
		err := Copy(source, dest, 7000, 0, pb)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("Offset exceed file size", func(t *testing.T) {
		err := Copy(source, dest, 7000, 0, pb)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("Offset 0 Limit 0", func(t *testing.T) {
		err := Copy(source, dest, 0, 0, pb)
		require.Nil(t, err)
		AssertFilesEqual(t, "testdata/out_offset0_limit0.txt", dest)
	})

	t.Run("Offset 0 Limit 10", func(t *testing.T) {
		err := Copy(source, dest, 0, 10, pb)
		require.Nil(t, err)
		AssertFilesEqual(t, "testdata/out_offset0_limit10.txt", dest)
	})

	t.Run("Offset 0 Limit 1000", func(t *testing.T) {
		err := Copy(source, dest, 0, 1000, pb)
		require.Nil(t, err)
		AssertFilesEqual(t, "testdata/out_offset0_limit1000.txt", dest)
	})

	t.Run("Offset 0 Limit 10000", func(t *testing.T) {
		err := Copy(source, dest, 0, 10000, pb)
		require.Nil(t, err)
		AssertFilesEqual(t, "testdata/out_offset0_limit10000.txt", dest)
	})

	t.Run("Offset 100 Limit 1000", func(t *testing.T) {
		err := Copy(source, dest, 100, 1000, pb)
		require.Nil(t, err)
		AssertFilesEqual(t, "testdata/out_offset100_limit1000.txt", dest)
	})

	t.Run("Offset 6000 Limit 1000", func(t *testing.T) {
		err := Copy(source, dest, 6000, 1000, pb)
		require.Nil(t, err)
		AssertFilesEqual(t, "testdata/out_offset6000_limit1000.txt", dest)
	})
}

func AssertFilesEqual(t *testing.T, expected, actual string) {
	t.Helper()
	c1, _ := ioutil.ReadFile(expected)
	c2, _ := ioutil.ReadFile(actual)
	require.Equal(t, c1, c2)
}
