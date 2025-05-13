package application

import (
	"balancer/internal/constants"
	"balancer/internal/healthcheck"
	"balancer/internal/methods"
	"balancer/internal/multiplexer"
	"balancer/pkg/logger"
	"context"
	"errors"
	"fmt"
)

func Serve(
	ctx context.Context,
	host string,
	port int,
	balanceMethod string,
	servers []string,
	weights []string,
	healthCheckInterval int,
	healthCheckTimeout int,
	withLog bool,
) error {

	logger.Init()

	m := constants.BalanceMethod(balanceMethod)

	switch m {
	case constants.RoundRobin:
		method := &methods.RoundRobin{}
		multiplexer.SetBalanceMethod(method)
	case constants.WeightedRoundRobin:
		method := &methods.WeightedRoundRobin{}
		multiplexer.SetBalanceMethod(method)
	case constants.LeastConnections:
		method := &methods.LeastConnections{}
		multiplexer.SetBalanceMethod(method)
	default:
		return ErrInvalidBalanceMethod
	}

	hc := healthcheck.New(
		servers,
		healthCheckInterval,
		healthCheckTimeout,
	)

	go hc.Start(ctx)

	if err := multiplexer.Multiplex(
		host,
		port,
		withLog,
	); err != nil {
		return fmt.Errorf("multiplexing failed: %w", err)
	}

	return nil
}

var (
	ErrInvalidBalanceMethod = errors.New("invalid balance method")
)
