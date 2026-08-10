package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/program"
	"alwis.dev/selectify/selectify-pkg/app"
	"alwis.dev/selectify/selectify-pkg/routers"
)

func main() {
	ctx := context.Background()
	logger.Init()

	// Set application prefix for environment variable prefixing
	program.AppPrefix = "API"

	if err := loadEnv(ctx); err != nil {
		panic(err)
	}

	app.NewAppEnvironment()

	h := routers.CreateHandler()

	srv := &http.Server{
		Addr:    ":3001",
		Handler: h,
	}
	go func() {
		logger.Infof(ctx, "Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal(ctx, err, "Server failed to start")
		}
	}()

	logger.Info(ctx, "Server started successfully")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		_ = logger.Errorf(ctx, err, "Server forced to shutdown")
	}
	app.Close()

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
