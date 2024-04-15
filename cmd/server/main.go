package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/ilya372317/pass-keeper/internal/server/app"
	_ "github.com/ilya372317/pass-keeper/proto"
)

func main() {
	a, err := app.New("config/config.yaml")
	if err != nil {
		panic(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	if err = a.StartGRPCServer(ctx); err != nil {
		panic(err)
	}
}
