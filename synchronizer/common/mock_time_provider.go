package common

import "time"

type MockTimerProvider struct {
	now time.Time
}

func (m *MockTimerProvider) Now() time.Time {
	return m.now
}
