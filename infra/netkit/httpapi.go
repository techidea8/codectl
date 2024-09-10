package netkit

import (
	"bytes"
	"io"
	"net/http"
)

// post 请求
func Post(url string, query map[string]string, header map[string]string, buf *bytes.Buffer) (body []byte, err error) {
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return
	}

	// appending to existing query args
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	for k, v := range header {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	return responseBody, err
}

// 请求header
func Get(url string, query map[string]string, header map[string]string) (body []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	// appending to existing query args
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	for k, v := range header {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	return responseBody, err
}
