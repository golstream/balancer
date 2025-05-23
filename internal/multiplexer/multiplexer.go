package multiplexer

import (
	httputils "balancer/pkg/httputil"
	"balancer/pkg/utils"
	"fmt"
	"log/slog"
	"net/http"
)

func Multiplex(
	host string,
	port int64,
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
	defaultProxyHandler = func(
		w http.ResponseWriter,
		r *http.Request) {

		body, err := httputils.ReadBody(r.Body)
		if err != nil {
			return
		}
		if err = r.Body.Close(); err != nil {
			return
		}

		statusCode, respBody, err :=
			method.
				Balance(
					r.Method,
					r.URL,
					r.URL.Query(),
					body,
					r.Header,
					r.Cookies(),
				)
		if err != nil {
			return
		}

		w.WriteHeader(statusCode)

		w.Header().Set(
			httputils.ContentTypeHeader,
			httputils.MIMEApplicationJSON,
		)

		if _, err = w.Write(respBody); err != nil {
			return
		}
	}

	loggedProxyHandler = func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "Request", "Method", r.Method, "URL", r.URL.String(), "Header", r.Header)
		defaultProxyHandler(w, r)
		slog.InfoContext(r.Context(), "Response", "Method", r.Method, "URL", r.URL.String(), "Header", r.Header)
	}
)
