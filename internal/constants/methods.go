package constants

import "errors"

type (
	BalanceMethod string
)

const (
	RoundRobin         BalanceMethod = "round_robin"
	WeightedRoundRobin BalanceMethod = "weighted_round_robin"
	LeastConnections   BalanceMethod = "least_connections"
)

var (
	ErrNoHosts = errors.New("no hosts available")
)
