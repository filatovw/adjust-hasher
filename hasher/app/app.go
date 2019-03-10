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

func New(w io.Writer, p Params) App {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	logger := newLogger(p.Debug)
	pooler := pool{
		log:      logger,
		urls:     p.URLs,
		parallel: p.Parallel,
		client:   http.Client{Transport: tr},
		w:        w,
	}
	return App{
		pool: pooler,
		log:  logger,
	}
}

func (a App) Run() error {
	return a.pool.exec(download)
}
