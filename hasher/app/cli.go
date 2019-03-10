package app

import (
	"errors"
	"flag"
)

type Params struct {
	Debug    bool
	Parallel int
	URLs     []string
}

func ReadParams() (Params, error) {
	p := Params{}
	flag.IntVar(&p.Parallel, "parallel", 10, "number of parallel processes")
	flag.BoolVar(&p.Debug, "debug", false, "debug mode")
	flag.Parse()

	p.URLs = flag.Args()
	if len(p.URLs) == 0 {
		return p, errors.New("no urls provided")
	}
	return p, nil
}
