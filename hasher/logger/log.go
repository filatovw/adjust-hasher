package logger

import (
	"io"
	"log"
)

func New(w io.Writer) *log.Logger {
	return log.New(w, "hasher ### ", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
}
