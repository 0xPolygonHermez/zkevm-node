package e2e_test

import (
	"testing"

	"github.com/hermeznetwork/hermez-core/db"
)

var cfg = db.Config{
	Database: "polygon-hermez",
	User:     "hermez",
	Password: "polygon",
	Host:     "localhost",
	Port:     "5432",
}

func TestStateTransition(t *testing.T) {

}
