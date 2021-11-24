package pool

import (
	"os"
	"testing"

	"github.com/hermeznetwork/hermez-core/log"
)

func TestMain(m *testing.M) {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})

	code := m.Run()
	os.Exit(code)
}
