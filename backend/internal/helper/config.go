package helper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"alwis.dev/selectify/internal/logger"
)

func LoadEnv(ctx context.Context) error {
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
