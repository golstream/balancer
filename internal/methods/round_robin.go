package methods

import (
	"balancer/internal/constants"
	"balancer/internal/multiplexer"
	"net/http"
	"net/url"
	"sync/atomic"
)

type RoundRobin struct {
	index atomic.Int32
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (r *RoundRobin) Balance(
	method string,
	url *url.URL,
	query url.Values,
	request []byte,
	headers http.Header,
	cookies []*http.Cookie) (int, []byte, error) {

	// atomic receiving and casting any to []string
	hosts := multiplexer.GetSliceOfURLs()

	if len(hosts) == 0 {
		return 0, nil, constants.ErrNoHosts
	}

	selected := r.index.Add(1) % int32(len(hosts))
	if selected < 0 {
		selected = -selected
	}

	url.Host = hosts[selected]

	return int(selected), nil, nil
}
