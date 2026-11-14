//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what denormalizeIntervalKo — render an English canonical tick interval back to its KO form
package ko

import "strings"

// denormalizeIntervalKo converts an English canonical interval token (as stored in
// the AST / .geul stream) back to its KO 매 ~마다 form: s→초, m→분, h→시간 (so
// "30s"→"30초", "5m"→"5분", "1h"→"1시간"). A clock time such as "7:00" carries no unit
// suffix (it never ends in h/s/m) and passes through verbatim. It is the inverse of
// normalizeIntervalKo.
func denormalizeIntervalKo(iv string) string {
	units := []struct{ en, ko string }{{"h", "시간"}, {"s", "초"}, {"m", "분"}}
	for _, u := range units {
		if num := strings.TrimSuffix(iv, u.en); num != iv && num != "" {
			return num + u.ko
		}
	}
	return iv
}
