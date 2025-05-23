package multiplexer

import (
	"net/http"
	"net/url"
	"sync/atomic"
)

type server struct {
	*http.ServeMux
}

var (
	Urls atomic.Value
)

var (
	method balancer
)

type balancer interface {
	Balance(string, *url.URL, url.Values, []byte, http.Header, []*http.Cookie) (int, []byte, error)
}

func SetBalanceMethod(balancer balancer) {
	method = balancer
}

func GetSliceOfURLs() []string {
	return Urls.Load().([]string)
}
