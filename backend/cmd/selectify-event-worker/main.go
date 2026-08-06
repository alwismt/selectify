package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	"alwis.dev/selectify/event-worker-pkg/app"
	"alwis.dev/selectify/event-worker-pkg/handlers"
	"alwis.dev/selectify/event-worker-pkg/worker"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/program"
)

func main() {
	ctx := context.Background()
	logger.Init()

	program.AppPrefix = "EVT"

	if err := loadEnv(ctx); err != nil {
		logger.Info(ctx, fmt.Sprintf("env load warning: %v", err))
	}

	app.NewAppEnvironment()
	handlers.InitiateHandlerRegistry()

	runCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info(runCtx, "selectify-event-worker starting")
	worker.NewSQSWorker(app.Repository().EventRepo).Run(runCtx)
	logger.Info(runCtx, "selectify-event-worker exited")
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
