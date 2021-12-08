package transport

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMGetOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.URL.Path))
		},
	))

	paths := makePaths()
	urls := urlsFromPaths(ts.URL, paths)
	expected := make(map[string]FetchResult, len(urls))
	for i, url := range urls {
		expected[url] = FetchResult{[]byte(paths[i]), nil}
	}

	actual := MGet(urls)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}

func TestMGetConnectionError(t *testing.T) {
	ts := httptest.NewServer(nil)
	ts.Close()

	_, err := http.Get(ts.URL)
	expected := map[string]FetchResult{ts.URL: FetchResult{nil, err}}

	actual := MGet([]string{ts.URL})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}

func TestMGetStatusCodeNotOK(t *testing.T) {
	badPath := fmt.Sprintf("/%d", rand.Intn(8))
	badResp := "Bad request"

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != badPath {
				w.Write([]byte(r.URL.Path))
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(badResp))
			}
		},
	))

	paths := makePaths()
	urls := urlsFromPaths(ts.URL, paths)

	badUrl := ts.URL + badPath
	err := fmt.Errorf("Status code: %d", http.StatusBadRequest)
	expected := make(map[string]FetchResult, len(urls))
	for i, url := range urls {
		var res FetchResult
		if url == badUrl {
			res = FetchResult{[]byte(badResp), err}
		} else {
			res = FetchResult{[]byte(paths[i]), nil}
		}
		expected[url] = res
	}

	actual := MGet(urls)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}

func makePaths() []string {
	paths := make([]string, 8)
	for i, _ := range paths {
		paths[i] = fmt.Sprintf("/%d", i)
	}

	return paths
}

func urlsFromPaths(base string, paths []string) []string {
	urls := make([]string, len(paths))
	for i, path := range paths {
		urls[i] = base + path
	}

	return urls
}
