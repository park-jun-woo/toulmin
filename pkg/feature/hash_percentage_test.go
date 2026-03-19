package feature

import (
	"fmt"
	"math"
	"testing"
)

func TestHashPercentageDeterministic(t *testing.T) {
	a := hashPercentage("user-abc")
	b := hashPercentage("user-abc")
	if a != b {
		t.Errorf("expected deterministic: %f != %f", a, b)
	}
}

func TestHashPercentageRange(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h := hashPercentage(fmt.Sprintf("user-%d", i))
		if h < 0 || h >= 1 {
			t.Fatalf("out of range [0,1): %f for user-%d", h, i)
		}
	}
}

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
