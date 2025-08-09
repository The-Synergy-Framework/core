package ids

import (
	"testing"
	"time"
)

func TestUUID(t *testing.T) {
	u, err := NewUUID()
	if err != nil || !IsUUID(u) {
		t.Fatalf("invalid uuid: %v %v", u, err)
	}
	if !IsUUID(MustUUID()) {
		t.Fatalf("MustUUID invalid")
	}
}

func TestULID(t *testing.T) {
	s, err := NewULID(time.Unix(0, 0))
	if err != nil || !IsULID(s) {
		t.Fatalf("invalid ulid: %v %v", s, err)
	}
	if !IsULID(MustULID()) {
		t.Fatalf("MustULID invalid")
	}
}

func TestPrefixed(t *testing.T) {
	s, err := Prefixed("user")
	if err != nil {
		t.Fatalf("prefixed error: %v", err)
	}
	if len(s) < 5 || s[:5] != "user_" {
		t.Fatalf("wrong prefix: %v", s)
	}
}
