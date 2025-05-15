package application

import (
	"balancer/internal/configuration"
	"balancer/internal/constants"
	"balancer/internal/healthcheck"
	"balancer/internal/methods"
	"balancer/internal/multiplexer"
	"balancer/pkg/logger"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

func Serve(ctx context.Context) error {
	logger.Init()

	cfg, err := configuration.Init()
	if err != nil {
		slog.ErrorContext(ctx, "cfg parse failed",
			slog.Any("error", err))

		return err
	}

	go healthcheck.New(
		cfg.Servers,
		cfg.HealthCheckInterval,
		cfg.HealthCheckTimeout,
	).Start(ctx)

	switch cfg.Method {
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

	if err = multiplexer.Multiplex(
		cfg.Host,
		cfg.Port,
		cfg.WithLog,
	); err != nil {
		return fmt.Errorf("multiplexing failed: %w", err)
	}

	return nil
}

var (
	ErrInvalidBalanceMethod = errors.New("invalid balance method")
)
