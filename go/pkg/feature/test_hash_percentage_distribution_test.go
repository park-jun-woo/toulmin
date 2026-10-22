//ff:func feature=feature type=util control=iteration dimension=1
//ff:what TestHashPercentageDistribution — tests hash percentage has uniform distribution
package feature

import (
	"fmt"
	"math"
	"testing"
)

func TestHashPercentageDistribution(t *testing.T) {
	n := 10000
	buckets := 10
	counts := make([]int, buckets)
	for i := 0; i < n; i++ {
		h := hashPercentage(fmt.Sprintf("user-%d", i))
		bucket := int(h * float64(buckets))
		if bucket >= buckets {
			bucket = buckets - 1
		}
		counts[bucket]++
	}
	expected := float64(n) / float64(buckets)
	for i, c := range counts {
		deviation := math.Abs(float64(c)-expected) / expected
		if deviation > 0.25 {
			t.Errorf("bucket %d: %d (%.1f%% deviation from expected %.0f)", i, c, deviation*100, expected)
		}
	}
}
