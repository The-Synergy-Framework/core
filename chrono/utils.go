package chrono

import "time"

const (
	FiveMinutes    = 5 * time.Minute
	FifteenMinutes = 15 * time.Minute
	Day            = 24 * time.Hour
	Week           = 7 * Day
)

// Now returns the current time using the Default clock.
func Now() time.Time { return Default.Now() }

// Since returns the duration since t using the Default clock.
func Since(t time.Time) time.Duration { return Default.Since(t) }

// IsExpired reports whether t is strictly before now.
func IsExpired(t time.Time) bool { return t.Before(Now()) }

// FormatApprox formats a duration into a short human-ish form (e.g., 1h2m, 3m4s, 5s, 250ms).
func FormatApprox(d time.Duration) string {
	if d < 0 {
		return "-" + FormatApprox(-d)
	}
	switch {
	case d >= time.Hour:
		h := d / time.Hour
		m := (d % time.Hour) / time.Minute
		if m == 0 {
			return itoa(int(h)) + "h"
		}
		return itoa(int(h)) + "h" + itoa(int(m)) + "m"
	case d >= time.Minute:
		m := d / time.Minute
		s := (d % time.Minute) / time.Second
		if s == 0 {
			return itoa(int(m)) + "m"
		}
		return itoa(int(m)) + "m" + itoa(int(s)) + "s"
	case d >= time.Second:
		s := d / time.Second
		return itoa(int(s)) + "s"
	default:
		ms := d / time.Millisecond
		return itoa(int(ms)) + "ms"
	}
}

// itoa is a tiny, allocation-free integer to ASCII converter for small non-negative ints.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + (n % 10))
		n /= 10
	}
	return string(buf[i:])
}
