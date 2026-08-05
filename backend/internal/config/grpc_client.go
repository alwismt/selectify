package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"

	"alwis.dev/selectify/internal/program"
)

type GRPCClient struct {
	Address string `envconfig:"address" default:"localhost:3002"`
}

func LoadGrpcPaymentConfig() (*GRPCClient, error) {
	var wrapper struct {
		GRPC GRPCClient `envconfig:"grpc_payment"`
	}

	prefix := program.AppPrefix

	err := envconfig.Process(prefix, &wrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to load database configuration: %w", err)
	}

	return &wrapper.GRPC, nil
}
