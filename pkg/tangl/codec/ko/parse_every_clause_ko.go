//ff:func feature=tangl type=parser control=sequence
//ff:what parseEveryClauseKo — parse '<간격>마다[, `C` 까지]' into a canonical interval and until case
package ko

import "strings"

// parseEveryClauseKo parses the every-tick clause after 매: the interval (its 마다
// suffix stripped, then normalized to the English canonical form by
// normalizeIntervalKo) and an optional ", `C` 까지" until clause. The comma splits
// the interval from the until clause; 까지 is a keyword (not a particle).
func parseEveryClauseKo(s string, line int, path string) (interval string, until string, err error) {
	parts := strings.SplitN(s, ",", 2)
	raw := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(parts[0]), "마다"))
	iv, ok := normalizeIntervalKo(raw)
	if !ok {
		return "", "", errAtKo(path, line, "unrecognized interval %q", raw)
	}
	if len(parts) < 2 {
		return iv, "", nil
	}
	untilBody, ok := takeSuffixKo(strings.TrimSpace(parts[1]), "까지")
	if !ok {
		return "", "", errAtKo(path, line, "expected '까지' in until clause, got %q", parts[1])
	}
	name, rest, ok := takeBacktick(untilBody)
	if !ok {
		return "", "", errAtKo(path, line, "expected backtick-quoted case name before '까지'")
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return "", "", errAtKo(path, line, "unexpected trailing text: %q", rest)
	}
	return iv, name, nil
}
