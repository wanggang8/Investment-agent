// Package idgen provides shared business ID generation helpers.
package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

// Generator creates readable business IDs with a caller-provided prefix.
type Generator interface {
	New(prefix string) string
}

// RandomGenerator is the production ID generator.
type RandomGenerator struct{}

// NewGenerator creates the production ID generator.
func NewGenerator() RandomGenerator { return RandomGenerator{} }

// New returns a readable ID in the form prefix_randomhex.
func (RandomGenerator) New(prefix string) string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return prefix + "_fallback"
	}
	return prefix + "_" + hex.EncodeToString(buf[:])
}

// FixedGenerator returns deterministic IDs for tests.
type FixedGenerator struct {
	mu     sync.Mutex
	ids    map[string][]string
	cursor map[string]int
}

// NewFixedGenerator creates a deterministic test generator.
func NewFixedGenerator(ids map[string][]string) *FixedGenerator {
	return &FixedGenerator{ids: ids, cursor: map[string]int{}}
}

// New returns the next configured ID for prefix, or a deterministic fallback.
func (g *FixedGenerator) New(prefix string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	idx := g.cursor[prefix]
	g.cursor[prefix] = idx + 1
	if values := g.ids[prefix]; idx < len(values) {
		return values[idx]
	}
	return prefix + "_fixed"
}
