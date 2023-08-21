package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidIP(t *testing.T) {
	var tests = []struct {
		name     string
		ip       string
		expected bool
	}{
		{"Valid IPv4", "127.0.0.1", true},
		{"Valid IPv6", "2001:db8:0:1:1:1:1:1", true},
		{"Invalid IP", "300.0.0.1", false},
		{"Empty IP", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidIP(tt.ip)
			assert.Equal(t, tt.expected, result)
		})
	}
}
