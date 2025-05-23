package methods

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

func TestNewWeightedRoundRobin_ErrorWhenWeightsLessThanHosts(t *testing.T) {
	_, err := NewWeightedRoundRobin([]string{"a", "b"}, []int{1})
	if !errors.Is(err, ErrWeightsIsLessThenHosts) {
		t.Errorf("expected ErrWeightsIsLessThenHosts, got %v", err)
	}
}

func TestNewWeightedRoundRobin_Success(t *testing.T) {
	wrr, err := NewWeightedRoundRobin([]string{"a", "b"}, []int{1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(wrr.servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(wrr.servers))
	}
}

func TestWeightedRoundRobin_Balance_Distribution(t *testing.T) {
	wrr, _ := NewWeightedRoundRobin([]string{"a", "b"}, []int{1, 3})

	counts := make([]int, 2)
	reqCount := 1000
	var mu sync.Mutex

	for i := 0; i < reqCount; i++ {
		u, _ := url.Parse("http://test")
		ix, _, err := wrr.Balance(http.MethodGet, u, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		mu.Lock()
		counts[ix]++
		mu.Unlock()
	}

	t.Logf("Distribution: %+v", counts)

	expected := []float64{0.25, 0.75}
	for i := range counts {
		actualRatio := float64(counts[i]) / float64(reqCount)
		if actualRatio < expected[i]-0.05 || actualRatio > expected[i]+0.05 {
			t.Errorf("server %d: expected ~%.2f, got %.2f", i, expected[i], actualRatio)
		}
	}
}
