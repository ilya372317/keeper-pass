package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/ilya372317/pass-keeper/internal/server/app"
)

var configFilePath string

func main() {
	fmt.Println(configFilePath)
	a, err := app.New(configFilePath)
	if err != nil {
		panic(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	if err = a.StartGRPCServer(ctx); err != nil {
		panic(err)
	}
}
