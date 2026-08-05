.PHONY: test
test:
	go test ./... -tags="test" -count=1

build_proto:
	 protoc \
      --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      payment-pkg/grpc-proto/payment/payment.proto
