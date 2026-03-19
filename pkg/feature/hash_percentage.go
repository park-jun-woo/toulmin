//ff:func feature=feature type=util control=sequence
//ff:what hashPercentage: 사용자 ID를 결정론적 해시로 0.0~1.0 비율에 매핑
package feature

import "hash/fnv"

// hashPercentage maps a user ID to a deterministic float64 in [0.0, 1.0).
// Same ID always returns the same value. Not rand — deterministic.
func hashPercentage(userID string) float64 {
	h := fnv.New32a()
	h.Write([]byte(userID))
	return float64(h.Sum32()) / float64(1<<32)
}
