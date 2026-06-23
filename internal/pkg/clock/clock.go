// Package clock provides shared UTC time helpers.
package clock

import "time"

// Clock returns current time for production and deterministic tests.
type Clock interface {
	Now() time.Time
	NowRFC3339() string
}

// SystemClock uses wall-clock UTC time.
type SystemClock struct{}

// Now returns current UTC time.
func (SystemClock) Now() time.Time { return time.Now().UTC() }

// NowRFC3339 returns current UTC time formatted for persistence.
func (c SystemClock) NowRFC3339() string { return c.Now().Format(time.RFC3339) }

// FixedClock returns a configured time for tests.
type FixedClock struct{ Time time.Time }

// Now returns the configured UTC time.
func (c FixedClock) Now() time.Time { return c.Time.UTC() }

// NowRFC3339 returns the configured UTC time formatted for persistence.
func (c FixedClock) NowRFC3339() string { return c.Now().Format(time.RFC3339) }
