package app

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type Params struct {
	Parallel int
	URLs     []string
}

func ReadParams() (Params, error) {
	p := Params{}
	flag.IntVar(&p.Parallel, "parallel", 10, "number of parallel processes")
	flag.Parse()
	p.URLs = flag.Args()
	if err := p.validate(); err != nil {
		return p, err
	}

	return p, nil
}

func validateURL(u string) error {
	//_, err := url.ParseRequestURI(u)
	//return err
	return nil
}

func (p Params) validate() error {
	if len(p.URLs) == 0 {
		return errors.New("no urls provided")
	}
	msgs := []string{}
	for _, u := range p.URLs {
		if err := validateURL(u); err != nil {
			msgs = append(msgs, fmt.Sprintf("error: %s", err))
		}
	}
	if len(msgs) > 0 {
		return errors.New(strings.Join(msgs, ";"))
	}
	return nil
}
