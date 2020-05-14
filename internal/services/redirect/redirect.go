package redirect

import (
	"bytes"
	"net/http"
)

// Redirect : Hub
type Redirect struct {
	BaseURL string
}

// New : New
func New() *Redirect {
	return &Redirect{
		BaseURL: "http://localhost:3333/api",
	}
}

// Do : Do
func (h *Redirect) Do(url string, body []byte, r *http.Request) (*http.Response, error) {
	z := bytes.NewBuffer(body)
	proxyReq, err := http.NewRequest(r.Method, h.BaseURL+url, z)
	if err != nil {
		return nil, err
	}

	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	return client.Do(proxyReq)
}
