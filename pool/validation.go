package pool

import "net"

// IsValidIP returns true if the given string is a valid IP address
func IsValidIP(ip string) bool {
	return ip != "" && net.ParseIP(ip) != nil
}
