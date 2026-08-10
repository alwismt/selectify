package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/program"
	"alwis.dev/selectify/payment-pkg/app"
	"alwis.dev/selectify/payment-pkg/grpc-proto"
	"alwis.dev/selectify/payment-pkg/routers"
)

func main() {
	ctx := context.Background()
	logger.Init()

	// Set application prefix for environment variable prefixing
	program.AppPrefix = "PAY"

	if err := loadEnv(ctx); err != nil {
		panic(err)
	}

	app.NewAppEnvironment()

	h := routers.CreateHandler()
	srv := &http.Server{
		Addr:    ":3003",
		Handler: h,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal(ctx, err, "Server failed to start")
		}
	}()

	// gRPC
	grpcServer := grpc.NewServer()
	grpcproto.RegisterGrpc(grpcServer)

	go func() {
		lis, err := net.Listen("tcp", ":3002")
		if err != nil {
			logger.Fatal(ctx, err, "failed to listen for gRPC")
		}

		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal(ctx, err, "gRPC server failed")
		}
	}()

	logger.Info(ctx, "gRPC server started")
	logger.Info(ctx, "Server started successfully")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "Shutting down server...")

	// Stop accepting new gRPC requests and finish existing ones.
	grpcServer.GracefulStop()

	if err := srv.Shutdown(ctx); err != nil {
		_ = logger.Errorf(ctx, err, "Server forced to shutdown")
	}

	logger.Info(ctx, "Server exited")
}

func loadEnv(ctx context.Context) error {
	execPath, err := os.Executable()
	if err != nil {
		logger.Fatal(ctx, err, "Failed to get executable path")
	}

	execDir := filepath.Dir(execPath)
	execName := filepath.Base(execPath)
	execName = strings.TrimSuffix(execName, filepath.Ext(execName))
	envPath := filepath.Join(execDir, execName+".env")

	if err = godotenv.Load(envPath); err != nil {
		logger.Info(ctx, fmt.Sprintf("No .env file found at %s, using system environment variables", envPath))
		return err
	}
	return nil
}
