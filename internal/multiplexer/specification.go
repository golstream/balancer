package multiplexer

import (
	"net/http"
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
	Balance()
}

func SetBalanceMethod(balancer balancer) {
	method = balancer
}

func GetSliceOfURLs() []string {
	return Urls.Load().([]string)
}
