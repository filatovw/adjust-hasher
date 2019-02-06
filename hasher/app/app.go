package app

import (
	"io"
	"log"
	"net/http"
	"time"
)

type App struct {
	pool pool
	log  *log.Logger
	w    io.Writer
}

func New(w io.Writer, log *log.Logger, p Params) App {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	pooler := pool{
		log:      log,
		urls:     p.URLs,
		parallel: p.Parallel,
		w:        w,
		client:   http.Client{Transport: tr},
	}
	return App{
		pool: pooler,
		log:  log,
	}
}

func (a App) Run() error {
	return a.pool.exec(download)
}
