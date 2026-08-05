package grpcclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"alwis.dev/selectify/internal/config"
	"alwis.dev/selectify/internal/logger"
)

func NewPaymentClientConnection() *grpc.ClientConn {
	ctx := context.Background()
	cnf, err := config.LoadGrpcPaymentConfig()
	if err != nil {
		logger.Fatal(ctx, err, "failed to load grpc payment config")
	}

	if cnf.Address == "" {
		logger.Fatal(ctx, fmt.Errorf("API_GRPC_PAYMENT_ADDRESS is empty"), "grpc payment address is required")
	}

	conn, err := connectWithRetry(ctx, cnf.Address, 2, 2)
	if err != nil {
		logger.Fatal(ctx, err, "failed to connect to grpc payment service")
	}

	conn.Connect()

	return conn
}

func connectWithRetry(ctx context.Context, address string, maxRetries int, delay time.Duration) (*grpc.ClientConn, error) {
	var lastErr error
	delay = delay * time.Second
	for attempt := 1; attempt <= maxRetries; attempt++ {
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			return conn, nil
		}

		lastErr = err

		select {
		case <-time.After(delay):
			// retry
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, lastErr)
}
