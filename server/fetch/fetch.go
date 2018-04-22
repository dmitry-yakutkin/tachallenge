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

// Fetcher is supposed to be used to retreive data from external services.
type Fetcher interface {
	Get(link string) (*http.Response, error)
}

// NewHTTPFetcher creates Fetcher instances.
func NewHTTPFetcher() Fetcher {
	return httpFetcher{
		maxRequestProcessingDuration: maxRequestProcessingDuration,
		httpClient:                   http.Client{Timeout: maxExternalNumbersFetchingDuration},
	}
}

type httpFetcher struct {
	maxRequestProcessingDuration time.Duration
	httpClient                   http.Client
}

func (f httpFetcher) Get(link string) (*http.Response, error) {
	resp, err := f.httpClient.Get(link)
	return resp, err
}
