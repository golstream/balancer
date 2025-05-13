package server

import (
	"balancer/pkg/utils"
	"fmt"
	"log/slog"
	"net/http"
)

func Multiplex(
	host string,
	port int,
	balancer balancer,
	withLog bool) error {

	serv := server{http.NewServeMux()}
	addr := fmt.Sprintf("%s:%d", host, port)
	serv.regHandler(withLog)
	return serv.listen(addr, balancer)
}

func (s *server) regHandler(withLog bool) {
	h := utils.Ternary(withLog, loggedProxyHandler, defaultProxyHandler)
	s.HandleFunc("/", h)
}

func (s *server) listen(host string, balancer balancer) error {
	balanceAlgorithms = balancer
	return http.ListenAndServe(host, s)
}

var (
	defaultProxyHandler = func(w http.ResponseWriter, r *http.Request) {
		balanceAlgorithms.Balance()
	}

	loggedProxyHandler = func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "request", "method", r.Method, "url", r.URL.String(), "header", r.Header)
		balanceAlgorithms.Balance()
		slog.InfoContext(r.Context(), "response", "method", r.Method, "url", r.URL.String(), "header", r.Header)
	}
)
