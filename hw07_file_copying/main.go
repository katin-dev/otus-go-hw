package main

import (
	"flag"

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

type AppProgressBar struct {
	bar *pb.ProgressBar
}

func NewProgressBar() *AppProgressBar {
	return &AppProgressBar{pb.New(0)}
}

func (pb *AppProgressBar) Init(total int) {
	pb.bar.SetTotal(total)
}

func (pb *AppProgressBar) Increment() {
	pb.bar.Increment()
}

func (pb *AppProgressBar) Finish() {
	pb.bar.Finish()
}

func main() {
	flag.Parse()
	Copy(from, to, offset, limit, NewProgressBar())
}
