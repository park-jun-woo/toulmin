//ff:func feature=tangl type=codegen control=sequence
//ff:what normalizeInterval — parses a tangl interval into a time.Duration
package gen

import (
	"strings"
	"time"
)

// normalizeInterval parses a tangl interval like "30s" or "1hour" into a
// time.Duration, expanding the handful of word-suffixes tangl allows
// (hour/hours/minute/minutes/second/seconds) that time.ParseDuration
// itself doesn't accept. A clock-of-day schedule like "7:00" isn't a
// duration at all and is reported as unsupported (ok=false) — the caller
// falls back to a fixed interval and flags it with a comment.
func normalizeInterval(s string) (time.Duration, bool) {
	r := strings.NewReplacer("hours", "h", "hour", "h", "minutes", "m", "minute", "m", "seconds", "s", "second", "s")
	norm := r.Replace(s)
	d, err := time.ParseDuration(norm)
	if err != nil {
		return 0, false
	}
	return d, true
}
