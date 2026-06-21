package service_test

import (
	"os"
	"testing"

	"github.com/mystaline/paperid-test/pkg/db"
)

func TestMain(m *testing.M) {
	db.InitIDGenerator()
	os.Exit(m.Run())
}
