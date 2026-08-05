package app

import (
	"context"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
)

var appEnv *applicationEnvironment

func Database() *db.DatabaseConnection {
	return appEnv.dbConn
}

//func Repository() *Repo {
//	return appEnv.repo
//}

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

	// Initialize database connection
	database, err := db.NewDatabaseConnection()
	if err != nil {
		logger.Fatal(ctx, err, "Failed to initialize database connection")
	}

	return database
}
