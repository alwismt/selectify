package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"alwis.dev/selectify/internal/config"
	"alwis.dev/selectify/internal/logger"
)

type DatabaseConnection struct {
	RoDb *sqlx.DB
	RwDb *sqlx.DB
}

func NewDatabaseConnection() (*DatabaseConnection, error) {
	// Load database configuration
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		logger.Fatal(context.Background(), err, "Failed to load database configuration")
	}

	roDb, err := openReadOnlyDb(dbConfig)
	if err != nil {
		return nil, err
	}

	rwDb, err := openReadWriteDb(dbConfig)
	if err != nil {
		_ = roDb.Close()
		return nil, err
	}

	return &DatabaseConnection{
		RwDb: rwDb,
		RoDb: roDb,
	}, err
}

func openReadWriteDb(dbConfig *config.DatabaseConfiguration) (*sqlx.DB, error) {
	return openDb(dbConfig.ReadWriteDbUrl(), dbConfig.DatabaseConnectionSettings)
}

func openReadOnlyDb(dbConfig *config.DatabaseConfiguration) (*sqlx.DB, error) {
	return openDb(dbConfig.ReadOnlyDbUrl(), dbConfig.DatabaseConnectionSettings)
}

func openDb(url string, dbConnSettings config.DatabaseConnectionSettings) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(dbConnSettings.MaxIdle)
	db.SetMaxOpenConns(dbConnSettings.MaxOpen)
	db.SetConnMaxLifetime(dbConnSettings.MaxLifetime)

	return db, nil
}

func (db *DatabaseConnection) Close() (err error) {
	if err = db.RoDb.Close(); err != nil {
		_ = logger.Errorf(context.TODO(), err, "Failed to close RoDb, %v", err)
	}

	if err = db.RwDb.Close(); err != nil {
		_ = logger.Errorf(context.TODO(), err, "Failed to close RwDb, %v", err)
	}
	return
}
