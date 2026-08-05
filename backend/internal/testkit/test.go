//go:build test

package testkit

import (
	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/program"
	"context"
	"errors"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
)

type TestSetup struct {
	C  context.Context
	DB *db.DatabaseConnection
}

func NewTestSetup() *TestSetup {
	program.AppPrefix = "api"
	logger.Init()

	filePath, err := getTestEnvironmentFile()
	if err != nil {
		_ = logger.Error(context.Background(), err, "failed to get env files")
		os.Exit(1)
	}

	err = godotenv.Load(filePath)
	if err != nil {
		_ = logger.Error(context.Background(), err, "failed to load env")
		os.Exit(1)
	}

	return &TestSetup{
		C: context.Background(),
	}
}

func (ts *TestSetup) ConnectDatabase() {
	if ts.DB == nil {
		var err error
		ts.DB, err = db.NewDatabaseConnection()
		if err != nil {
			_ = logger.Error(ts.C, err, "failed to connect to database")
			os.Exit(1)
		}
	}
}

func getTestEnvironmentFile() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to get caller")
	}

	basePath := filepath.Dir(currentFilePath)
	filePath := filepath.Join(basePath, "test.env")
	return filePath, nil
}
