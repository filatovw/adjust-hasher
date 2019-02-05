package logger

import (
	"io"
	"log"
	"os"
	"strings"
)

func New(w io.Writer) *log.Logger {
	stdlog := log.New(w, "hasher ### ", log.Lmicroseconds|log.LstdFlags)
	if debug, ok := os.LookupEnv("DEBUG"); ok {
		if strings.ToLower(debug) == "true" {
			stdlog.SetFlags(log.Llongfile)
		}
	}
	return stdlog
}
