package app

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type App struct {
	pool   pool
	log    *log.Logger
	w      io.Writer
	params Params
	client Doer
}
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func New(w io.Writer, params Params) App {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	logger := newLogger(params.Debug)
	p := pool{
		log:      logger,
		urls:     params.URLs,
		parallel: params.Parallel,
		client:   http.Client{Transport: tr},
		w:        w,
	}
	return App{
		client: &http.Client{Transport: tr},
		params: params,
		pool:   p,
		log:    logger,
		w:      w,
	}
}

type res struct {
	err  error
	body string
	url  string
}

// Run execute downloading for a bunch of urls
func (a App) Run() error {
	/*
		return a.pool.exec(download)
	*/
	a.log.Printf("config %+v", a.params)
	wg := &sync.WaitGroup{}
	queue := make(chan string)
	results := make(chan res)

	wg.Add(1)
	go func() {
		defer close(queue)
		a.log.Printf("queue sender started")
		for _, url := range a.params.URLs {
			a.log.Printf("send to queue: %s", url)
			queue <- url
		}
		wg.Done()
	}()

	for i := 0; i < a.params.Parallel; i++ {
		i := i
		wg.Add(1)
		go func(queue <-chan string, results chan<- res, wg *sync.WaitGroup) {
			a.log.Printf("worker %d started", i)
			for job := range queue {
				if requestedURL, body, err := download(a.client, job); err != nil {
					a.log.Printf("downloaded with error: %s", err)
					results <- res{
						err: err,
					}
				} else {
					a.log.Printf("%s downloaded successfully", requestedURL)
					results <- res{
						body: fmt.Sprintf("%x", md5.Sum(body)),
						url:  requestedURL,
					}
				}
			}
			wg.Done()
			a.log.Printf("worker done: %d", i)
		}(queue, results, wg)
	}

	go func() {
		wg.Wait()
		a.log.Printf("close results")
		close(results)
	}()

	for result := range results {
		if result.err != nil {
			a.log.Printf("error: %s", result.err)
		} else {
			fmt.Fprintf(a.w, "%s %s\n", result.url, result.body)
		}
	}

	return nil
}
