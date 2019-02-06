package main

import (
	"fmt"
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
	resultsOutput := os.Stdout

	a := app.New(resultsOutput, stdlog, params)
	if err := a.Run(); err != nil {
		fmt.Fprintf(resultsOutput, "%s", err)
		os.Exit(1)
	}
	os.Exit(0)
}
