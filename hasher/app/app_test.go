package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownload(t *testing.T) {
	expectedContent := `welcome to my kingdom`

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(expectedContent))
	}))
	defer server.Close()

	expectedURL := server.URL

	URL, content, err := download(server.Client(), expectedURL)
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
