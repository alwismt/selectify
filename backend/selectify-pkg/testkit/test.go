//go:build test

package testkit

import (
	"net/http/httptest"

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
	if t.DB != nil {
		_ = t.DB.Close()
	}
	if t.S != nil {
		t.S.Close()
	}
}

func CreateServer() *httptest.Server {
	app.NewAppEnvironment()
	return httptest.NewServer(routers.CreateHandler())
}
