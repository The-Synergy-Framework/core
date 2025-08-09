package ids

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"
)

// NewUUID generates a random UUID v4 as a lowercase string with hyphens.
func NewUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("uuid: rand: %w", err)
	}
	// Set version (4) and variant (10xx)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

// MustUUID generates a UUID v4 or panics.
func MustUUID() string {
	s, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return s
}

// IsUUID reports whether s looks like a valid UUID v4 string.
func IsUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !isHex(byte(c)) {
				return false
			}
		}
	}
	// version 4 at pos 14 (0-based)
	if s[14] != '4' {
		return false
	}
	// variant at pos 19 must be one of 8,9,a,b
	s19 := s[19]
	return s19 == '8' || s19 == '9' || s19 == 'a' || s19 == 'b'
}

func isHex(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}

// Crockford base32 alphabet (no I, L, O, U)
const crockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// NewULID generates a ULID string for time t using crypto-random entropy.
// ULID is a 26-char, lexicographically sortable identifier.
func NewULID(t time.Time) (string, error) {
	// 48-bit timestamp (ms)
	ts := uint64(t.UnixNano() / 1e6)
	var entropy [10]byte
	if _, err := rand.Read(entropy[:]); err != nil {
		return "", fmt.Errorf("ulid: rand: %w", err)
	}
	var buf [26]byte
	encodeTime(ts, buf[0:10])
	encodeEntropy(entropy[:], buf[10:26])
	return string(buf[:]), nil
}

// MustULID generates a ULID for now or panics.
func MustULID() string {
	s, err := NewULID(time.Now())
	if err != nil {
		panic(err)
	}
	return s
}

// IsULID reports whether s looks like a valid ULID (26 Crockford chars).
func IsULID(s string) bool {
	if len(s) != 26 {
		return false
	}
	s = strings.ToUpper(s)
	for i := 0; i < 26; i++ {
		if indexCrock(s[i]) < 0 {
			return false
		}
	}
	return true
}

func indexCrock(b byte) int {
	switch b {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return int(b - '0')
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H':
		return int(b-'A') + 10
	case 'J', 'K':
		return int(b-'J') + 18
	case 'M', 'N':
		return int(b-'M') + 20
	case 'P', 'Q', 'R', 'S', 'T', 'V', 'W', 'X', 'Y', 'Z':
		return int(b-'P') + 22
	default:
		return -1
	}
}

func encodeTime(ts uint64, dst []byte) {
	// produce 10 chars, msb first
	for i := 9; i >= 0; i-- {
		dst[i] = crockford[ts&31]
		ts >>= 5
	}
}

func encodeEntropy(src []byte, dst []byte) {
	// src: 10 bytes (80 bits) -> dst: 16 chars (80 bits)
	var acc uint32
	var bits uint8
	pos := 0
	for i := 0; i < 10; i++ {
		acc = (acc << 8) | uint32(src[i])
		bits += 8
		for bits >= 5 {
			bits -= 5
			idx := (acc >> bits) & 31
			dst[pos] = crockford[idx]
			pos++
		}
	}
}

// Prefixed returns prefix + '_' + uuid v4.
func Prefixed(prefix string) (string, error) {
	u, err := NewUUID()
	if err != nil {
		return "", err
	}
	if prefix == "" {
		return u, nil
	}
	return prefix + "_" + u, nil
}

// ParseOrErr validates an id with a validator; returns error on invalid.
func ParseOrErr(id string, valid func(string) bool) error {
	if valid == nil {
		return errors.New("ids: nil validator")
	}
	if !valid(id) {
		return fmt.Errorf("ids: invalid id: %s", id)
	}
	return nil
}
