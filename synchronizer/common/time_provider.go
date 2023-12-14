package common

import (
	"time"
)

// TimeProvider is a interface for classes that needs time and we want to be able to unittest it
type TimeProvider interface {
	// Now returns current time
	Now() time.Time
}

// DefaultTimeProvider is the default implementation of TimeProvider
type DefaultTimeProvider struct{}

// Now returns current time
func (d DefaultTimeProvider) Now() time.Time {
	return time.Now()
}
