package operations

import "os"

const (
	Rollup             = "rollup"
	Validium           = "validium"
	Irrelevant         = "irrelevant"
	DockerConcensusENV = "ZKEVM_CONCENSUS"
	TestConcensusENV   = "ZKEVM_TEST_CONCENSUS"
)

// IsConcensusRelevant returns true if the test is supposed to run in Validium / Rollup mode, false if it's irrelevant
func IsConcensusRelevant() bool {
	consensus := os.Getenv(TestConcensusENV)
	return consensus == Rollup || consensus == Validium
}

// IsRollup returns true if the test is supposed to run in Rollup mode
func IsRollup() bool {
	consensus := os.Getenv(TestConcensusENV)
	return consensus == Rollup
}

// IsValidium returns true if the test is supposed to run in Validium mode
func IsValidium() bool {
	consensus := os.Getenv(TestConcensusENV)
	return consensus == Validium
}
