package ym

import "net/http"

type (
	httpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)
