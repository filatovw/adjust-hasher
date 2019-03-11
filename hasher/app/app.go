package app

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// App instance of hasher
type App struct {
	log    *log.Logger
	w      io.Writer
	params Params
	client Doer
}

// Doer does http queries
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// New create application
func New(w io.Writer, params Params) App {
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	return App{
		client: &http.Client{Transport: tr},
		params: params,
		log:    newLogger(params.Debug),
		w:      w,
	}
}

type result struct {
	err  error
	body string
	url  string
}

// Run execute downloading for a bunch of urls
func (a App) Run() error {
	a.log.Printf("config %+v", a.params)
	wg := &sync.WaitGroup{}
	queue := make(chan string)
	results := make(chan result)

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
		go func(queue <-chan string, results chan<- result, wg *sync.WaitGroup) {
			a.log.Printf("worker %d started", i)
			for job := range queue {
				if requestedURL, body, err := download(a.client, job); err != nil {
					a.log.Printf("downloaded with error: %s", err)
					results <- result{
						err: err,
					}
				} else {
					a.log.Printf("%s downloaded successfully", requestedURL)
					results <- result{
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

	for r := range results {
		if r.err != nil {
			a.log.Printf("error: %s", r.err)
		} else {
			if _, err := fmt.Fprintf(a.w, "%s %s\n", r.url, r.body); err != nil {
				return fmt.Errorf("error writing results: %s", err)
			}
		}
	}

	return nil
}

// download URL content
func download(client Doer, u string) (string, []byte, error) {
	up, err := url.Parse(u)
	if err != nil {
		return "", nil, fmt.Errorf("download: %s", err)
	}
	if !up.IsAbs() {
		up.Scheme = "http"
	}
	req, err := http.NewRequest(http.MethodGet, up.String(), nil)
	if err != nil {
		return "", nil, fmt.Errorf("download: %s", err)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	req.Header.Add("accept-language", "en-US,en;q=0.9,ru;q=0.8")

	res, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("download: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return "", nil, fmt.Errorf("download: [%d] %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, fmt.Errorf("download: %s", err)
	}

	return up.String(), body, nil
}
