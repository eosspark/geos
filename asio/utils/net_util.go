package utils

import "errors"

const (
	IPv4len = 4
)

type FD int
type IP []byte

// Decimal to integer starting at &s[i0].
// Returns number, new offset, success.
func dtoi(s string, i0 int) (n int, i int, ok bool) {
	n = 0
	for i = i0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= 0xFFFFFF {
			return 0, i, false
		}
	}
	if i == i0 {
		return 0, i, false
	}
	return n, i, true
}

func ParseIPv4(s string) ([IPv4len]byte, error) {
	var p [IPv4len]byte
	i := 0
	for j := 0; j < IPv4len; j++ {
		if i >= len(s) {
			// Missing octets.
			return p, errors.New("parseIPv4 failed1")
		}
		if j > 0 {
			if s[i] != '.' {
				return p, errors.New("parseIPv4 failed2")
			}
			i++
		}
		var (
			n  int
			ok bool
		)
		n, i, ok = dtoi(s, i)
		if !ok || n > 0xFF {
			return p, errors.New("parseIPv4 failed3")
		}
		p[j] = byte(n)
	}
	if i != len(s) {
		return p, errors.New("parseIPv4 failed4")
	}
	return p, nil
}

func ParsePort(port string) (int, error) {
	p, i, ok := dtoi(port, 0)
	if !ok || i != len(port) {
		return 0, errors.New("invalid port 1")
	}
	if p < 0 || p > 0xFFFF {
		return 0, errors.New("invalid port 2")
	}
	return p, nil
}

