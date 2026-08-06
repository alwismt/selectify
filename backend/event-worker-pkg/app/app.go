package app

import (
	"context"
	"io"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
)

var appEnv *applicationEnvironment

func Database() *db.DatabaseConnection {
	return appEnv.dbConn
}

func Repository() *Repo {
	return appEnv.repo
}

func Service() *SVC {
	return appEnv.service
}

type applicationEnvironment struct {
	dbConn  *db.DatabaseConnection
	repo    *Repo
	service *SVC
}

func NewAppEnvironment() {
	appEnv = new(applicationEnvironment)
	appEnv.dbConn = databaseConnection()
	appEnv.repo = NewRepository()
	appEnv.service = NewService()
}

func databaseConnection() *db.DatabaseConnection {
	ctx := context.Background()

	database, err := db.NewDatabaseConnection()
	if err != nil {
		logger.Fatal(ctx, err, "Failed to initialize database connection")
	}

	return database
}

// noopStorage satisfies StorageService without requiring S3 for the event worker.
type noopStorage struct{}

func (noopStorage) UploadFile(_ context.Context, _ io.Reader, _ string, _ string) error {
	return nil
}

func (noopStorage) DeleteFile(_ context.Context, _ string) error {
	return nil
}
