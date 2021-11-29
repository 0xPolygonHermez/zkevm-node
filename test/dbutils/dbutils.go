package dbutils

import (
	"io"
	"os/exec"
	"time"

	"github.com/hermeznetwork/hermez-core/log"
)

const (
	dockerInstanceName = "test-instance"
	// DBHost is the host for the PostgreSQL instance
	DBHost = "localhost"
	// DBPort is the port for the PostgreSQL instance
	DBPort = "5432"
)

// StartPostgreSQL starts a docker PostgreSQL server with the database initialized with optional sqlFile
func StartPostgreSQL(dbName string, dbUser, dbPassword, sqlFile string) error {
	// Make sure we have the image
	cmd := exec.Command("/usr/bin/docker", "pull", "postgres")
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Start the container
	cmd1 := exec.Command("/usr/bin/docker", "run", "--rm", "--name", dockerInstanceName, "-p", DBPort+":5432", "-e",
		"POSTGRES_PASSWORD="+dbPassword, "-e", "POSTGRES_USER="+dbUser, "-e", "POSTGRES_DB="+dbName, "postgres") // #nosec
	err = cmd1.Start()
	if err != nil {
		return err
	}

	const safeTimeDelay = 5
	time.Sleep(safeTimeDelay * time.Second)

	// Check if we have to run a SQL Script
	if sqlFile != "" {
		cmd2 := exec.Command("/usr/bin/cat", sqlFile) // #nosec
		stdout, err := cmd2.CombinedOutput()
		if err != nil {
			return err
		}

		cmd3 := exec.Command("docker", "exec", "-i", dockerInstanceName, "psql", "-d", dbName, "-U", dbUser) // #nosec
		stdin, err := cmd3.StdinPipe()
		if err != nil {
			return err
		}

		go func() {
			defer func(wc io.WriteCloser) {
				_ = wc.Close()
			}(stdin)
			if _, err = io.WriteString(stdin, string(stdout)); err != nil {
				log.Error(err)
			}
		}()

		err = cmd3.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// StopPostgreSQL stops the docker PostgreSQL
func StopPostgreSQL() error {
	cmd := exec.Command("/usr/bin/docker", "stop", dockerInstanceName)
	err := cmd.Start()

	if err != nil {
		return err
	}

	return nil
}
