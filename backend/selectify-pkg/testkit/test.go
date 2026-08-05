//go:build test

package testkit

import (
	"net/http/httptest"
	"os"

	"alwis.dev/selectify/internal/testkit"
	"alwis.dev/selectify/selectify-pkg/app"
	"alwis.dev/selectify/selectify-pkg/routers"
)

func NewTestSetup() *TestSetup {
	ts := new(TestSetup)
	ts.TestSetup = testkit.NewTestSetup()
	ts.ConnectDatabase()
	ts.S = CreateServer()

	return ts
}

type TestSetup struct {
	*testkit.TestSetup
	S *httptest.Server
}

func (t TestSetup) Close() {
	_ = t.DB.Close()
	t.S.Close()
	os.Exit(0)
}

func CreateServer() *httptest.Server {
	app.NewAppEnvironment()
	return httptest.NewServer(routers.CreateHandler())
}
