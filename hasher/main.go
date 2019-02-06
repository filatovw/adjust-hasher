package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/filatovw/adjust-hasher/hasher/app"
	"github.com/filatovw/adjust-hasher/hasher/logger"
)

func main() {
	logOutput := ioutil.Discard
	if debug, ok := os.LookupEnv("DEBUG"); ok && strings.ToLower(debug) == "true" {
		logOutput = os.Stdout
	}
	stdlog := logger.New(logOutput)
	params, err := app.ReadParams()
	if err != nil {
		stdlog.Fatal(err)
	}

	a := app.New(os.Stdout, stdlog, params)
	if err := a.Run(); err != nil {
		stdlog.Fatal(err)
	}
}
