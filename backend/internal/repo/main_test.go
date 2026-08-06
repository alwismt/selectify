package repo_test

import (
	"context"
	"os"
	"testing"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/internal/testkit"
	"alwis.dev/selectify/internal/types"

	"github.com/jmoiron/sqlx"
)

var (
	ctx    context.Context
	ts     *testkit.TestSetup
	dbConn *sqlx.DB

	testProduct  = &model.Product{ID: 1, SKU: "TSET-SKU-5352206841292483884", Name: "Test Product"}
	testProduct2 = &model.Product{ID: 2, SKU: "TSET-SKU-5398977550428877284", Name: "Test Product"}
	testUser     = &model.User{ID: 1, Email: "travis@alwis.dev", FirstName: "Travis", LastName: "Alwis", Phone: "+37012345678", Status: types.UserStatusActive}

	// Repos
	productRepo         repo.ProductRepo
	productFileRepo     repo.ProductFileRepo
	productVariantsRepo repo.ProductVariantsRepo
	transactionRepo     repo.TransactionRepo
	userRepo            repo.UserRepo
	userFileRepo        repo.UserFileRepo
	userRoleRepo        repo.UserRoleRepo
	userSessionRepo     repo.UserSessionRepo
)

func TestMain(m *testing.M) {
	ts = testkit.NewTestSetup()
	ts.ConnectDatabase()
	setupRepo()
	dbConn = ts.DB.RwDb
	ctx = ts.C

	code := m.Run()

	if dbConn != nil {
		_ = dbConn.Close()
	}
	os.Exit(code)
}

func setupRepo() {
	productRepo = repo.NewProductRepo(ts.DB)
	productFileRepo = repo.NewProductFileRepo(ts.DB)
	productVariantsRepo = repo.NewProductVariantRepo(ts.DB)
	transactionRepo = repo.NewTransactionRepository(ts.DB)
	userRepo = repo.NewUserRepo(ts.DB)
	userFileRepo = repo.NewUserFileRepo(ts.DB)
	userRoleRepo = repo.NewUserRoleRepo(ts.DB)
	userSessionRepo = repo.NewUserSessionRepo(ts.DB)
}
