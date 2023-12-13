package common

import "time"

// MockTimerProvider is a mock implementation of the TimerProvider interface that return the internal variable
type MockTimerProvider struct {
	now time.Time
}

// Now in the implementation of TimeProvider.Now()
func (m *MockTimerProvider) Now() time.Time {
	return m.now
}
