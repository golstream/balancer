package main

import (
	"balancer/internal/application"
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	ctx, cancel := signal.
		NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)
	defer cancel()

	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("invalid port %w", err)
	}
	method := os.Getenv("METHOD")
	servers := strings.Split(os.Getenv("SERVERS"), ",")
	weights := strings.Split(os.Getenv("WEIGHTS"), ",")
	healthCheckInterval, err := strconv.Atoi(os.Getenv("HEALTH_CHECK_INTERVAL"))
	if err != nil {
		log.Fatalf("invalid health check interval %w", err)
	}
	healthCheckTimeout, err := strconv.Atoi(os.Getenv("HEALTH_CHECK_TIMEOUT"))
	if err != nil {
		log.Fatalf("invalid health check timeout %w", err)
	}

	withLog, err := strconv.ParseBool(os.Getenv("WITH_LOG"))
	if err != nil {
		log.Fatalf("invalid with log %w", err)
	}

	go func() {
		log.Fatalf("service is down %w",
			application.Serve(ctx, host, port, method, servers, weights, healthCheckInterval, healthCheckTimeout, withLog))
	}()

	select {}
}
