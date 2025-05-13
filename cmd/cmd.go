package main

import (
	"balancer/internal/application"
	_ "balancer/pkg/logger"
	"context"
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

	application.Serve(ctx)

	select {}
}
