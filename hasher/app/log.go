package app

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

func newLogger(debug bool) *log.Logger {
	var w io.Writer
	w = ioutil.Discard
	if debug {
		w = os.Stdout
	}
	return log.New(w, "hasher ### ", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
}
