package methods

import (
	"balancer/internal/constants"
	"errors"
	"net/url"
	"testing"
)

func TestRoundRobin_NoHosts_ReturnsError(t *testing.T) {
	mockHosts = []string{}

	r := NewRoundRobin(mockHosts)
	u, _ := url.Parse("http://placeholder")

	_, _, err := r.Balance("GET", u, nil, nil, nil, nil)
	if !errors.Is(err, constants.ErrNoHosts) {
		t.Errorf("expected ErrNoHosts, got %v", err)
	}
}

func TestRoundRobin_CyclesThroughHosts(t *testing.T) {
	mockHosts = []string{"host1", "host2", "host3"}

	r := NewRoundRobin(mockHosts)
	u := &url.URL{}

	calls := 6
	expectedHosts := []string{"host2", "host3", "host1", "host2", "host3", "host1"}

	for i := 0; i < calls; i++ {
		urlCopy := *u
		_, _, err := r.Balance("GET", &urlCopy, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if urlCopy.Host != expectedHosts[i] {
			t.Errorf("unexpected host: got %s, want %s", urlCopy.Host, expectedHosts[i])
		}
	}
}

func TestRoundRobin_IndexNeverNegative(t *testing.T) {
	mockHosts = []string{"host1", "host2"}

	r := NewRoundRobin(mockHosts)
	r.index.Store(-100)

	u := &url.URL{}
	_, _, err := r.Balance("GET", u, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Host == "" {
		t.Errorf("expected host to be set, got empty")
	}
}

var (
	mockHosts []string
)
