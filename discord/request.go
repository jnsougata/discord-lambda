package discord

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Method  string
	Body    []byte
	Path    string
	Token   string
	Headers map[string]string
}

func (r *Request) Do() (*http.Response, error) {
	url := "https://discord.com/api/v10" + r.Path
	req, _ := http.NewRequest(r.Method, url, nil)
	if r.Body != nil {
		req.Body = io.NopCloser(bytes.NewReader(r.Body))
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	if r.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bot %s", r.Token))
	}
	return http.DefaultClient.Do(req)
}
