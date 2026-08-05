package app

import "alwis.dev/selectify/internal/repo"

type Repo struct {
	OrderRepo   repo.OrderRepo
	PaymentRepo repo.PaymentRepo
}

func NewRepository() *Repo {
	repository := new(Repo)

	repository.OrderRepo = repo.NewOrderRepo(appEnv.dbConn)
	repository.PaymentRepo = repo.NewPaymentRepo(appEnv.dbConn)

	return repository
}
