package app

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
)

func TestDownload(t *testing.T) {
	expectedContent := `welcome to my kingdom`

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(expectedContent))
	}))
	defer server.Close()

	expectedURL := server.URL

	URL, content, err := download(*server.Client(), expectedURL)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if string(content) != expectedContent {
		t.Errorf(`Content # expected: "%s", got: "%s"`, expectedContent, content)
	}
	if URL != expectedURL {
		t.Errorf(`URL # expected: %s, got: %s`, expectedURL, URL)
	}
}

func TestPool(t *testing.T) {
	// dummy server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))
	defer server.Close()
	client := server.Client()

	testData := []struct {
		urls []string
		name string
	}{
		{
			urls: []string{},
			name: "empty",
		},
		{
			urls: []string{
				"http://google.com",
			},
			name: "one url",
		},
		{
			urls: []string{
				"http://google.com",
				"http://google.de",
				"http://google.org",
				"http://google.ru",
				"http://google.gb",
				"http://google",
				"https://google",
			},
			name: "many urls",
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			var b bytes.Buffer
			p := pool{
				client:   *client,
				log:      log.New(os.Stdout, "", log.LstdFlags),
				urls:     td.urls,
				parallel: 3,
				w:        &b,
			}
			var i int32
			dummyJob := func(client http.Client, u string) (string, []byte, error) {
				atomic.AddInt32(&i, 1)
				return u, []byte(u), nil
			}
			err := p.exec(dummyJob)
			if err != nil {
				t.Errorf("unexpected err: %s", err)
			}
			if len(td.urls) != int(i) {
				t.Errorf("expected exec number: %d, got: %d", len(td.urls), i)
			}
			for _, u := range td.urls {
				hash := md5.Sum([]byte(u))
				if !strings.Contains(b.String(), fmt.Sprintf("%x", hash)) {
					t.Errorf("output doesn't contain expected hash: %x", hash)
				}
			}
		})
	}
}
