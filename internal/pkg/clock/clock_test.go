package clock

import (
	"testing"
	"time"
)

func TestSystemClockReturnsUTC(t *testing.T) {
	got := SystemClock{}.Now()
	if got.Location() != time.UTC {
		t.Fatalf("expected UTC location, got %v", got.Location())
	}
}

func TestFixedClockReturnsDeterministicRFC3339(t *testing.T) {
	fixed := time.Date(2026, 5, 29, 4, 0, 0, 0, time.UTC)
	c := FixedClock{Time: fixed}
	if got := c.NowRFC3339(); got != "2026-05-29T04:00:00Z" {
		t.Fatalf("NowRFC3339 = %q", got)
	}
}
