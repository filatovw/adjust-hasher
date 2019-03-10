package main

import (
	"log"
	"os"

	"github.com/filatovw/adjust-hasher/hasher/app"
)

func main() {
	params, err := app.ReadParams()
	if err != nil {
		log.Fatal(err)
	}

	a := app.New(os.Stdout, params)
	if err := a.Run(); err != nil {
		log.Print(err)
	}
}
