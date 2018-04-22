package fetch

import (
	"net/http"
	"time"
)

const (
	maxRequestProcessingDuration       = time.Millisecond * 500
	maxExternalNumbersFetchingDuration = time.Millisecond * 400
	linkKey                            = "u"
)

type Fetcher interface {
	Get(link string) (*http.Response, error)
}

type httpFetcher struct {
	maxRequestProcessingDuration time.Duration
	httpClient                   http.Client
}

func (f httpFetcher) Get(link string) (*http.Response, error) {
	resp, err := f.httpClient.Get(link)
	return resp, err
}

func New() Fetcher {
	return httpFetcher{
		maxRequestProcessingDuration: maxRequestProcessingDuration,
		httpClient:                   http.Client{Timeout: maxExternalNumbersFetchingDuration},
	}
}
