package transport

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type FetchResult struct {
	Buffer []byte
	Err    error
}

type fetchData struct {
	url string
	FetchResult
}

func MGet(urls []string) map[string]FetchResult {
	length := len(urls)
	c := make(chan fetchData)

	for _, url := range urls {
		go fetchToChan(url, c)
	}

	res := make(map[string]FetchResult, length)
	received := 0

	for received < length {
		data := <-c
		res[data.url] = FetchResult{data.Buffer, data.Err}
		received += 1
	}

	return res
}

func fetchToChan(url string, ch chan fetchData) {
	var buffer []byte

	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()

		buffer, err = ioutil.ReadAll(resp.Body)
		if err == nil && resp.StatusCode/100 != 2 {
			err = fmt.Errorf("Status code: %d", resp.StatusCode)
		}
	}

	res := FetchResult{Buffer: buffer, Err: err}
	data := fetchData{url: url, FetchResult: res}
	ch <- data
}
