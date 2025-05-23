package methods

import (
	"net/url"
	"testing"
)

func TestLeastConnections_NoHosts_ReturnsError(t *testing.T) {
	lb := NewLeastConnections([]string{})

	u := &url.URL{}
	_, _, err := lb.Balance("GET", u, nil, nil, nil, nil)
	if err == nil || err.Error() != "no hosts available" {
		t.Errorf("expected error 'no hosts available', got %v", err)
	}
}

func TestLeastConnections_HostsAndConnectionsMismatch_ReturnsError(t *testing.T) {
	lb := NewLeastConnections([]string{"host1", "host2"})
	lb.connections = []int{1}

	u := &url.URL{}
	_, _, err := lb.Balance("GET", u, nil, nil, nil, nil)
	if err == nil || err.Error() != "internal error: hosts and connections mismatch" {
		t.Errorf("expected mismatch error, got %v", err)
	}
}

func TestLeastConnections_ChoosesHostWithLeastConnections(t *testing.T) {
	lb := NewLeastConnections([]string{"host1", "host2", "host3"})
	lb.connections = []int{5, 2, 3}

	u := &url.URL{}
	idx, _, err := lb.Balance("GET", u, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if idx != 1 || u.Host != "host2" {
		t.Errorf("expected host2 at index 1, got index=%d, host=%s", idx, u.Host)
	}

	conns := lb.connections
	if conns[1] != 3 {
		t.Errorf("expected incremented connection count at index 1, got %d", conns[1])
	}
}

func TestLeastConnections_Release_DecrementsConnections(t *testing.T) {
	lb := NewLeastConnections([]string{"host1"})
	lb.connections = []int{2}

	lb.Release(0)

	if conns := lb.connections; conns[0] != 1 {
		t.Errorf("expected connection count to be 1, got %d", conns[0])
	}
}

func TestLeastConnections_Release_IgnoresInvalidIndex(t *testing.T) {
	lb := NewLeastConnections([]string{"host1"})
	lb.connections = []int{2}

	lb.Release(10)
	lb.Release(-1)
	lb.Release(0)

	if conns := lb.connections; conns[0] != 1 {
		t.Errorf("expected connection count to be 1, got %d", conns[0])
	}
}
