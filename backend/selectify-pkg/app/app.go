package app

import (
	"context"

	"google.golang.org/grpc"

	"alwis.dev/selectify/internal/db"
	grpcclient "alwis.dev/selectify/internal/grpc-client"
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

func GrpcClient() *ProtoClient {
	return appEnv.grpcCli
}

func Close() {
	if appEnv == nil || appEnv.grpcPayClient == nil {
		return
	}
	if err := appEnv.grpcPayClient.Close(); err != nil {
		logger.Warnf(context.Background(), "failed to close grpc payment connection, %v", err)
	}
}

type applicationEnvironment struct {
	dbConn        *db.DatabaseConnection
	grpcPayClient *grpc.ClientConn
	repo          *Repo
	service       *SVC
	grpcCli       *ProtoClient
}

func NewAppEnvironment() {
	appEnv = new(applicationEnvironment)
	appEnv.dbConn = databaseConnection()
	appEnv.grpcPayClient = grpcPaymentClientConnection()
	appEnv.repo = NewRepository()
	appEnv.grpcCli = NewGrpcClient()
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

func grpcPaymentClientConnection() *grpc.ClientConn {
	grpcClient := grpcclient.NewPaymentClientConnection()
	return grpcClient
}
