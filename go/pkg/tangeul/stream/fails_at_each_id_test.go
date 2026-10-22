//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what failsAtEachID — test helper: call must fail when each of its first n issued Edge Ids is the first past the 1-word width
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// failsAtEachID runs call n times on a fresh test encoder over tt whose counter
// is set so that the k-th issued Edge Id (k = 1..n) is 0x10000, one past the
// 1-word width. Every run must fail: this walks each encode-error path of an
// emitter in issue order.
func failsAtEachID(t *testing.T, n int, tt *tidTable, call func(e *enc) error) {
	t.Helper()
	for k := 0; k < n; k++ {
		e := testEnc(tt)
		e.next = 0x10000 - toulmin.TID(k)
		if err := call(e); err == nil {
			t.Fatalf("Edge Id overflow at issue %d: want error", k+1)
		}
	}
}
