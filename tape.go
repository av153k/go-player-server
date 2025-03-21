package poker

import (
	"io"
	"os"
)

type Tape struct {
	file *os.File
}

func (t *Tape) Write(p []byte) (n int, e error) {
	t.file.Truncate(0)
	t.file.Seek(0, io.SeekStart)
	return t.file.Write(p)
}
