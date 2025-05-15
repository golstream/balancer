package methods

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
)

type LeastConnections struct {
	hosts       []string
	connections []int
	mu          sync.Mutex
}

func NewLeastConnections(hosts []string) *LeastConnections {
	return &LeastConnections{
		hosts:       append([]string{}, hosts...),
		connections: make([]int, len(hosts)),
	}
}

func (l *LeastConnections) Balance(
	method string,
	u *url.URL,
	query url.Values,
	request []byte,
	headers http.Header,
	cookies []*http.Cookie) (int, []byte, error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.hosts) == 0 {
		return 0, nil, errors.New("no hosts available")
	}
	if len(l.connections) != len(l.hosts) {
		return 0, nil, errors.New("internal error: hosts and connections mismatch")
	}

	minIndex := 0
	for i := 1; i < len(l.connections); i++ {
		if l.connections[i] < l.connections[minIndex] {
			minIndex = i
		}
	}

	l.connections[minIndex]++

	u.Host = l.hosts[minIndex]

	return minIndex, nil, nil
}

func (l *LeastConnections) Release(index int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index >= 0 && index < len(l.connections) && l.connections[index] > 0 {
		l.connections[index]--
	}
}
