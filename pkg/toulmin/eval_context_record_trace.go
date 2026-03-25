//ff:func feature=engine type=engine control=sequence
//ff:what recordTrace — executes rule function and appends trace entry on first visit
package toulmin

import (
	"fmt"
	"time"
)

// recordTrace runs the rule function for id and records a TraceEntry.
// Called once per node when calcTrace visits it for the first time.
// When duration is true, measures execution time of the rule function.
// Sets ec.err if the rule function panics.
func (ec *evalContext) recordTrace(id string, ctx Context, duration bool) {
	ec.ran[id] = true
	var dur time.Duration
	if duration {
		start := time.Now()
		active, evidence, err := safeCall(ec.fnMap[id], ctx, ec.specsMap[id])
		dur = time.Since(start)
		if err != nil {
			ec.err = fmt.Errorf("rule %q: %w", id, err)
			return
		}
		ec.active[id] = active
		ec.evidence[id] = evidence
	} else {
		active, evidence, err := safeCall(ec.fnMap[id], ctx, ec.specsMap[id])
		if err != nil {
			ec.err = fmt.Errorf("rule %q: %w", id, err)
			return
		}
		ec.active[id] = active
		ec.evidence[id] = evidence
	}
	role := ec.roleMap[id]
	if role == "" {
		role = inferRole(ec.strMap, ec.attackerSet, id)
	}
	ec.trace = append(ec.trace, TraceEntry{
		Name:      id,
		Role:      role,
		Activated: ec.active[id],
		Qualifier: ec.qualMap[id],
		Evidence:  ec.evidence[id],
		Specs:     ec.specsMap[id],
		Duration:  dur,
	})
}
