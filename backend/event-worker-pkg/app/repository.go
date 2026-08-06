package app

import "alwis.dev/selectify/internal/repo"

type Repo struct {
	EventRepo    repo.EventRepo
	UserRepo     repo.UserRepo
	UserFileRepo repo.UserFileRepo
	TxRepo       repo.TransactionRepo
}

func NewRepository() *Repo {
	repository := new(Repo)
	repository.EventRepo = repo.NewEventRepo(appEnv.dbConn)
	repository.UserRepo = repo.NewUserRepo(appEnv.dbConn)
	repository.UserFileRepo = repo.NewUserFileRepo(appEnv.dbConn)
	repository.TxRepo = repo.NewTransactionRepository(appEnv.dbConn)
	return repository
}
