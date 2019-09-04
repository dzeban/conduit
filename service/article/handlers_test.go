package article

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testSecret = "test"

func TestHandleArticleList(t *testing.T) {
	s := New(newMockStore(), testSecret)

	ts := httptest.NewServer(s.router)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("failed to make a request: %s", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status, want %d, got %d", http.StatusOK, resp.StatusCode)
		t.Error(string(body))
	}
}
