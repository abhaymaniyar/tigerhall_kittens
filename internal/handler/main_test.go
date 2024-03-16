package handler

import (
	"context"
	"tigerhall_kittens/test_helpers"

	"os"
	"testing"
)

func TestMain(m *testing.M) {
	test_helpers.LoadEnvForTest()
	test_helpers.InitializeLogger()
	test_helpers.SetupPostgresConnection(context.Background(), test_helpers.EnvConfig)
	os.Exit(m.Run())
}
