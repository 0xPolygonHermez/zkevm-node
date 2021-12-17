package pool

import (
	"os"
	"testing"

	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/log"
)

var cfg = db.Config{
	User:     "test_user",
	Password: "test_password",
	Name:     "test_db",
	Host:     "localhost",
	Port:     "5432",
}

func TestMain(m *testing.M) {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})

	code := m.Run()
	os.Exit(code)
}
