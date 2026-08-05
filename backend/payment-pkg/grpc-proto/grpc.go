package grpcproto

import (
	"google.golang.org/grpc"

	paymentgrpc "alwis.dev/selectify/payment-pkg/grpc"
	"alwis.dev/selectify/payment-pkg/grpc-proto/payment"
)

func RegisterGrpc(server *grpc.Server) {
	payment.RegisterPaymentServiceServer(server, paymentgrpc.NewController())
}
