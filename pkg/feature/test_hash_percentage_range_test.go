//ff:func feature=feature type=util control=iteration dimension=1
//ff:what TestHashPercentageRange — tests hash percentage values are in [0,1)
package feature

import (
	"fmt"
	"testing"
)

func TestHashPercentageRange(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h := hashPercentage(fmt.Sprintf("user-%d", i))
		if h < 0 || h >= 1 {
			t.Fatalf("out of range [0,1): %f for user-%d", h, i)
		}
	}
}
