package multiplexer

import (
	"balancer/pkg/utils"
	"fmt"
	"log/slog"
	"net/http"
)

func Multiplex(
	host string,
	port int,
	withLog bool) error {

	serv := server{http.NewServeMux()}
	addr := fmt.Sprintf("%s:%d", host, port)
	serv.regHandler(withLog)
	return serv.listen(addr)
}

func (s *server) regHandler(withLog bool) {
	h := utils.Ternary(withLog, loggedProxyHandler, defaultProxyHandler)
	s.HandleFunc("/", h)
}

func (s *server) listen(host string) error {
	return http.ListenAndServe(host, s)
}

var (
	defaultProxyHandler = func(w http.ResponseWriter, r *http.Request) {
		method.Balance()
	}

	loggedProxyHandler = func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "Request", "method", r.Method, "URL", r.URL.String(), "Header", r.Header)
		method.Balance()
		slog.InfoContext(r.Context(), "Response", "method", r.Method, "URL", r.URL.String(), "Header", r.Header)
	}
)
