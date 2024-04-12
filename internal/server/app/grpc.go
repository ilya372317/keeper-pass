package app

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/handler/grpc"
)

// StartGRPCServer starting grpc server. Block gorutine before grpc server will stopped by ctx parameter.
func (a *App) StartGRPCServer(ctx context.Context) error {
	server := grpc.New(a.conf)
	err := server.StartAndListen(ctx)
	if err != nil {
		return fmt.Errorf("failed start grpc server: %w", err)
	}
	return nil
}
