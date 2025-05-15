package methods

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type weightedServer struct {
	Host          string
	Weight        int
	CurrentWeight int
}

type WeightedRoundRobin struct {
	servers []*weightedServer
	mu      sync.Mutex
}

func NewWeightedRoundRobin(hosts []string, weights []int) (*WeightedRoundRobin, error) {
	if len(weights) < len(hosts) {
		return nil, ErrWeightsIsLessThenHosts
	}

	servers := make([]*weightedServer, len(hosts))
	for i, host := range hosts {
		servers[i] = &weightedServer{
			Host:   host,
			Weight: weights[i],
		}
	}

	return &WeightedRoundRobin{servers: servers}, nil
}

func (w *WeightedRoundRobin) Balance(
	method string,
	u *url.URL,
	query url.Values,
	request []byte,
	headers http.Header,
	cookies []*http.Cookie,
) (int, []byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.servers) == 0 {
		return 0, nil, fmt.Errorf("no servers available")
	}

	bestIndex, best := w.selectBestServer()
	if best == nil {
		return 0, nil, fmt.Errorf("unable to select server")
	}

	u.Host = best.Host
	return bestIndex, nil, nil
}

func (w *WeightedRoundRobin) selectBestServer() (int, *weightedServer) {
	var (
		totalWeight int
		best        *weightedServer
		bestIndex   = -1
	)

	for i, s := range w.servers {
		s.CurrentWeight += s.Weight
		totalWeight += s.Weight

		if best == nil || s.CurrentWeight > best.CurrentWeight {
			best = s
			bestIndex = i
		}
	}

	if best != nil {
		best.CurrentWeight -= totalWeight
	}

	return bestIndex, best
}

var (
	ErrWeightsIsLessThenHosts = errors.New("len of weights is less then hosts len")
)
