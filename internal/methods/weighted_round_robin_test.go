package methods

import (
	"errors"
	"fmt"
	"net/url"
	"testing"
)

func TestNewWeightedRoundRobin_Success(t *testing.T) {
	hosts := []string{"host1", "host2"}
	weights := []int{2, 3}

	wrr, err := NewWeightedRoundRobin(hosts, weights)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(wrr.servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(wrr.servers))
	}
}

func TestNewWeightedRoundRobin_ErrorWhenWeightsTooShort(t *testing.T) {
	hosts := []string{"host1", "host2"}
	weights := []int{1}

	_, err := NewWeightedRoundRobin(hosts, weights)
	if !errors.Is(err, ErrWeightsIsLessThenHosts) {
		t.Errorf("expected ErrWeightsIsLessThenHosts, got %v", err)
	}
}

func TestWeightedRoundRobin_BalanceReturnsErrorWhenNoServers(t *testing.T) {
	wrr := &WeightedRoundRobin{}
	u, _ := url.Parse("http://example.com")

	_, _, err := wrr.Balance("GET", u, nil, nil, nil, nil)
	if err == nil || err.Error() != "no servers available" {
		t.Errorf("expected 'no servers available' error, got %v", err)
	}
}

func TestWeightedRoundRobin_BalanceSelection(t *testing.T) {
	hosts := []string{"host1", "host2"}
	weights := []int{1, 2}

	wrr, err := NewWeightedRoundRobin(hosts, weights)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	counts := map[string]int{}
	u := &url.URL{}
	for i := 0; i < 100; i++ {
		urlCopy := *u
		_, _, err := wrr.Balance("GET", &urlCopy, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error on balance: %v", err)
		}
		counts[urlCopy.Host]++
	}

	if counts["host2"] <= counts["host1"] {
		t.Errorf("unexpected balance distribution: %v", counts)
	} else {
		fmt.Printf("Distribution: %v\n", counts)
	}
}

func TestSelectBestServer_AdjustsWeightsCorrectly(t *testing.T) {
	hosts := []string{"host1", "host2"}
	weights := []int{5, 1}

	wrr, _ := NewWeightedRoundRobin(hosts, weights)

	index, best := wrr.selectBestServer()
	if best.Host != "host1" || index != 0 {
		t.Errorf("expected host1 as best, got %v", best.Host)
	}

	expectedWeight := 5 - (5 + 1)

	if wrr.servers[0].CurrentWeight != expectedWeight {
		t.Errorf("expected CurrentWeight %d, got %d", expectedWeight, wrr.servers[0].CurrentWeight)
	}
}
