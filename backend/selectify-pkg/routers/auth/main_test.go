package auth_test

import (
	"os"
	"testing"

	"alwis.dev/selectify/selectify-pkg/testkit"
)

var ts *testkit.TestSetup

func TestMain(m *testing.M) {
	ts = testkit.NewTestSetup()
	code := m.Run()
	ts.Close()
	os.Exit(code)
}
