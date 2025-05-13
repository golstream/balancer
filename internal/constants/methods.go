package constants

type (
	BalanceMethod string
)

const (
	RoundRobin         BalanceMethod = "round_robin"
	WeightedRoundRobin BalanceMethod = "weighted_round_robin"
	LeastConnections   BalanceMethod = "least_connections"
)
