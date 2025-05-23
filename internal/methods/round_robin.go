package methods

import (
	"balancer/internal/constants"
	"math"
	"net/http"
	"net/url"
	"sync/atomic"
)

type RoundRobin struct {
	hosts []string
	index atomic.Int32
}

func NewRoundRobin(hosts []string) *RoundRobin {
	return &RoundRobin{
		hosts: hosts,
	}
}

func (r *RoundRobin) Balance(
	method string,
	url *url.URL,
	query url.Values,
	request []byte,
	headers http.Header,
	cookies []*http.Cookie) (int, []byte, error) {

	// atomic receiving and casting any to []string
	hosts := r.hosts

	if len(hosts) == 0 {
		return 0, nil, constants.ErrNoHosts
	}

	selected := r.index.Add(1) % int32(len(hosts))
	if selected < 0 {
		selected = int32(math.Abs(float64(selected)))
	}

	url.Host = hosts[selected]

	return int(selected), nil, nil
}
