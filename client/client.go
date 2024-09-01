package client

import (
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	GetEndpoint() string
}

var _ Client = (*httpsClient)(nil)

func NewClient(endpoint string) Client {
	return &httpsClient{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}

type httpsClient struct {
	endpoint string
	client   *http.Client
}

func (c *httpsClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newReq, err := createNewRequest(c.endpoint, r)
	if err != nil {
		// TODO: Handle errors.
		return
	}

	res, err := c.client.Do(newReq)
	if err != nil {
		// TODO: Handle errors.
		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			// TODO: Handle errors.
		}
	}(res.Body)

	fillResponseHeader(w, res)
	if _, err := io.Copy(w, res.Body); err != nil {
		// TODO: Handle errors.
		return
	}
}

func (c *httpsClient) GetEndpoint() string {
	return c.endpoint
}

func createNewRequest(endpoint string, r *http.Request) (*http.Request, error) {
	url := fmt.Sprintf("%s%s%s", scheme, endpoint, r.URL.Path)

	newReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		// TODO: Handle errors.
		return nil, err
	}

	newReq.Header = r.Header
	newReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

	return newReq, nil
}

func fillResponseHeader(w http.ResponseWriter, res *http.Response) {
	w.WriteHeader(res.StatusCode)
	for k, v := range res.Header {
		w.Header()[k] = v
	}
}

const scheme = "https://"
