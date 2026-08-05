package app

import (
	grpcclient "alwis.dev/selectify/internal/grpc-client"
)

type ProtoClient struct {
	PaymentClient grpcclient.PaymentClient
}

func NewGrpcClient() *ProtoClient {
	client := new(ProtoClient)

	client.PaymentClient = grpcclient.NewPaymentClient(appEnv.grpcPayClient)

	return client
}
