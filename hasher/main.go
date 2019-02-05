package main

import (
	"os"

	"github.com/filatovw/adjust-hasher/hasher/app"
	"github.com/filatovw/adjust-hasher/hasher/logger"
)

func main() {
	stdlog := logger.New(os.Stdout)
	params, err := app.ReadParams()
	if err != nil {
		stdlog.Fatal(err)
	}

	a := app.New(stdlog, params, os.Stdout)
	if err := a.Run(); err != nil {
		stdlog.Fatal(err)
	}
}
