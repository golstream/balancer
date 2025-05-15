package main

import (
	"balancer/internal/application"
	"context"
	"log"
	"os"
	"os/signal"
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

	go func() {
		log.Fatalf("service is down %w",
			application.Serve(ctx))
	}()

	select {}
}
