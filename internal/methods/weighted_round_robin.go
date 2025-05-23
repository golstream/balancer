package methods

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type weightedServer struct {
	Host           string
	Weight         int
	RequestCounter atomic.Int32
}

type WeightedRoundRobin struct {
	servers       []*weightedServer
	totalRequests atomic.Int32
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
	totalReq := w.totalRequests.Add(1)

	var (
		bestIndex = -1
		minOver   = 1.1
	)

	totalWeight := 0
	for _, s := range w.servers {
		totalWeight += s.Weight
	}

	for i, s := range w.servers {
		serverReq := float64(s.RequestCounter.Load())
		expectedShare := float64(s.Weight) / float64(totalWeight)
		actualShare := serverReq / float64(totalReq)

		over := actualShare - expectedShare
		if over <= 0 {
			s.RequestCounter.Add(1)
			u.Host = s.Host
			return i, nil, nil
		}

		if over < minOver {
			bestIndex = i
			minOver = over
		}
	}

	if bestIndex != -1 {
		s := w.servers[bestIndex]
		s.RequestCounter.Add(1)
		u.Host = s.Host
		return bestIndex, nil, nil
	}

	return 0, nil, errors.New("no servers available")
}

var (
	ErrWeightsIsLessThenHosts = errors.New("len of weights is less than hosts len")
)
