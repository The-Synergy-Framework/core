package chrono

import (
	"testing"
	"time"
)

type fakeClock struct{ t time.Time }

func (f fakeClock) Now() time.Time                  { return f.t }
func (f fakeClock) Since(t time.Time) time.Duration { return f.t.Sub(t) }

func TestIsExpired(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	SetDefault(fakeClock{t: base})
	defer SetDefault(nil)
	if IsExpired(base.Add(1 * time.Second)) {
		t.Fatalf("future should not be expired")
	}
	if !IsExpired(base.Add(-1 * time.Second)) {
		t.Fatalf("past should be expired")
	}
}

func TestFormatApprox(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{time.Hour + 2*time.Minute, "1h2m"},
		{3*time.Minute + 4*time.Second, "3m4s"},
		{5 * time.Second, "5s"},
		{250 * time.Millisecond, "250ms"},
	}
	for _, c := range cases {
		if got := FormatApprox(c.d); got != c.want {
			t.Fatalf("FormatApprox(%v)=%q, want %q", c.d, got, c.want)
		}
	}
}
