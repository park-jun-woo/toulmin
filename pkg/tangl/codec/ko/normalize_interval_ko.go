//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what normalizeIntervalKo — normalize a KO tick interval to its English canonical form
package ko

import "strings"

// normalizeIntervalKo converts a KO interval token (마다 already stripped) to the
// English canonical form the AST and .tangeul stream store: 초→s, 분→m, 시간→h (so
// "30초"→"30s", "5분"→"5m", "1시간"→"1h"). A clock time such as "7:00" carries no
// unit suffix and passes through verbatim. ok is false for an empty token.
func normalizeIntervalKo(raw string) (string, bool) {
	units := []struct{ ko, en string }{{"시간", "h"}, {"초", "s"}, {"분", "m"}}
	for _, u := range units {
		if !strings.HasSuffix(raw, u.ko) {
			continue
		}
		num := strings.TrimSpace(strings.TrimSuffix(raw, u.ko))
		if num == "" {
			return "", false
		}
		return num + u.en, true
	}
	if raw == "" {
		return "", false
	}
	return raw, true
}
