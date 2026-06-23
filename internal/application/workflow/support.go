package workflow

import (
	"strings"

	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

var workflowIDGen idgen.Generator = idgen.NewGenerator()
var workflowClock clock.Clock = clock.SystemClock{}

// SetWorkflowIDGenerator lets tests inject deterministic workflow IDs.
func SetWorkflowIDGenerator(gen idgen.Generator) {
	if gen != nil {
		workflowIDGen = gen
	}
}

// SetWorkflowClock lets tests inject deterministic workflow time.
func SetWorkflowClock(c clock.Clock) {
	if c != nil {
		workflowClock = c
	}
}

func workflowNowRFC3339() string {
	return workflowClock.NowRFC3339()
}

func workflowID(prefix string) string {
	return workflowIDGen.New(prefix)
}

func workflowStableID(prefix, hash string) string {
	return prefix + "_" + strings.TrimPrefix(hash, "sha256:")
}
