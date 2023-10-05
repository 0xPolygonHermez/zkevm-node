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
	concensus := os.Getenv(TestConcensusENV)
	return concensus == Rollup || concensus == Validium
}

// IsRollup returns true if the test is supposed to run in Rollup mode
func IsRollup() bool {
	concensus := os.Getenv(TestConcensusENV)
	return concensus == Rollup
}

// IsValidium returns true if the test is supposed to run in Validium mode
func IsValidium() bool {
	concensus := os.Getenv(TestConcensusENV)
	return concensus == Validium
}
