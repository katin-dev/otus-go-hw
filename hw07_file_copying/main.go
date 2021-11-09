package main

import (
	"flag"
	"io"
	"math"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	// 1. Открыть файл, если могу на чтение
	// 2. Открыть файл, если могу, на запись
	// 3. Проверить, что offset не выходит за размер файла
	// 4. Оказывается, у файла может быть неизвестная длина! Тогда fatal!
	// 5. Скопировать файл и отметить progress

	if from == "" {
		panic("-from must not be empty")
	}
	if to == "" {
		panic("-to must not be empty")
	}

	source, err := os.Open(from)
	if err != nil {
		panic("failed to open source file: " + err.Error())
	}
	defer source.Close()

	sourceInfo, err := source.Stat()
	if err != nil {
		panic("failed to stat source file: " + err.Error())
	}

	if offset != 0 && offset > sourceInfo.Size() {
		panic("offset is out of bound:")
	}

	dest, err := os.Create(to)
	if err != nil {
		panic("Failed to open dest file")
	}
	defer dest.Close()

	source.Seek(offset, io.SeekStart)

	start := offset
	stop := sourceInfo.Size()
	if limit != 0 {
		stop = int64(math.Min(float64(start+limit), float64(sourceInfo.Size())))
	}

	bytesToCopy := stop - start

	var bufferSize int = 4
	var writtenBytes int64

	// Рассчитаем progress bar - сколько буферов потребуется для копирования файла
	iterations := int(math.Ceil(float64(bytesToCopy) / float64(bufferSize)))

	bar := pb.StartNew(iterations)
	for writtenBytes < bytesToCopy {
		n := int64(math.Min(float64(bufferSize), float64(bytesToCopy-writtenBytes)))
		io.CopyN(dest, source, n)

		writtenBytes += n

		bar.Increment()
	}

	bar.Finish()
}
