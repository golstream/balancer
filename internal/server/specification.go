package server

import "net/http"

type server struct {
	*http.ServeMux
}

type balancer interface {
	Balance()
}

var (
	balanceAlgorithms balancer
)
