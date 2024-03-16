package handler

import (
	"context"
	"os"
	"testing"

	"tigerhall_kittens/test_helpers"
)

func TestMain(m *testing.M) {
	test_helpers.LoadEnvForTest()
	test_helpers.InitializeLogger()
	test_helpers.SetupPostgresConnection(context.Background(), test_helpers.EnvConfig)
	os.Exit(m.Run())
}
