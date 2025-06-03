package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type RequestOption func(*http.Request)

func WithHeaders(headers map[string]string) RequestOption {
	return func(req *http.Request) {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
}

func WithQueryParams(params url.Values) RequestOption {
	return func(req *http.Request) {
		if params != nil {
			req.URL.RawQuery = params.Encode()
		}
	}
}

func WithJSONBody(body []byte) RequestOption {
	return func(req *http.Request) {
		req.Body = io.NopCloser(&bytes.Buffer{})

		if body != nil {
			req.Header.Set("Content-Type", "application/json")
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}
}
