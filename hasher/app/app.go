package app

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type App struct {
	client   http.Client
	log      *log.Logger
	urls     []string
	parallel int
	w        io.Writer
}

func New(w io.Writer, log *log.Logger, p Params) App {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	return App{
		log:      log,
		urls:     p.URLs,
		parallel: p.Parallel,
		w:        w,
		client:   http.Client{Transport: tr},
	}
}

func (a App) Run() error {
	in := make(chan string, a.parallel)
	out := make(chan result, a.parallel)
	for i := a.parallel; i > 0; i-- {
		go worker(i, in, out, a.client)
	}

	go scheduler(a.urls, in)

	j := 0
	for j < len(a.urls) {
		select {
		case r := <-out:
			if r.err != nil {
				a.log.Printf("worker %d | %s", r.id, r.err)
			} else {
				a.log.Printf("%d: %s %x", r.id, r.url, r.hash)
				fmt.Fprintf(a.w, "%s %x\n", r.url, r.hash)
			}
			j++
		}
	}
	return nil
}

func scheduler(urls []string, in chan<- string) {
	for _, u := range urls {
		in <- u
	}
}

type result struct {
	url  string
	hash [md5.Size]byte
	err  error
	id   int
}

func worker(id int, in <-chan string, out chan<- result, client http.Client) {
	timer := time.NewTimer(time.Second * 30)
	for {
		select {
		case url, ok := <-in:
			if !ok {
				timer.Stop()
				goto WORK_END
			}
			realURL, content, err := download(client, url)
			if err != nil {
				out <- result{err: err, id: id}
				continue
			}
			hash := md5.Sum(content)
			out <- result{hash: hash, url: realURL, id: id}
		case <-timer.C:
			timer.Stop()
			goto WORK_END
		}
	}
WORK_END:
}

func download(client http.Client, u string) (string, []byte, error) {
	up, err := url.Parse(u)
	if err != nil {
		return "", nil, fmt.Errorf("download: %s", err)
	}
	if !up.IsAbs() {
		up.Scheme = "http"
	}
	req, err := http.NewRequest("GET", up.String(), nil)
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
