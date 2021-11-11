package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptySourceFilename   = errors.New("source filename is empty")
	ErrEmptyDestFilename     = errors.New("destination filename is empty")
)

type ProgressBar interface {
	Init(limit int)
	Increment()
	Finish()
}

func Copy(fromPath, toPath string, offset, limit int64, pb ProgressBar) error {
	if fromPath == "" {
		return ErrEmptySourceFilename
	}

	if toPath == "" {
		return ErrEmptyDestFilename
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", fromPath, err)
	}
	defer source.Close()

	sourceInfo, err := source.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	if offset != 0 && offset > sourceInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	dest, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to open destination file %s: %w", toPath, err)
	}
	defer dest.Close()

	source.Seek(offset, io.SeekStart)

	start := offset
	stop := sourceInfo.Size()
	if limit != 0 {
		stop = int64(math.Min(float64(start+limit), float64(sourceInfo.Size())))
	}

	bytesToCopy := stop - start

	bufferSize := 4
	var writtenBytes int64

	// Рассчитаем progress bar - сколько буферов потребуется для копирования файла
	iterations := int(math.Ceil(float64(bytesToCopy) / float64(bufferSize)))

	pb.Init(iterations)
	for writtenBytes < bytesToCopy {
		n := int64(math.Min(float64(bufferSize), float64(bytesToCopy-writtenBytes)))
		n, err := io.CopyN(dest, source, n)
		if err != nil {
			return err
		}

		writtenBytes += n

		pb.Increment()
	}

	pb.Finish()
	return nil
}
