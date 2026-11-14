//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunComposeDepthGuard — a long acyclic Run chain trips the runMaxDepth backstop
package toulmin

import (
	"strings"
	"testing"
)

func TestRunComposeDepthGuard(t *testing.T) {
	active := func(ctx Context, specs Specs) (bool, any) { return true, nil }

	// A long chain of DISTINCT graphs (no cycle) longer than runMaxDepth.
	n := runMaxDepth + 5
	graphs := make([]*Graph, n)
	for i := range graphs {
		graphs[i] = NewGraph("chain")
	}
	for i := 0; i < n; i++ {
		r := graphs[i].Rule(active)
		if i+1 < n {
			r.Run(graphs[i+1])
		}
	}

	_, _, err := graphs[0].Run(NewContext())
	if err == nil {
		t.Fatal("expected run depth exceeded error")
	}
	if !strings.Contains(err.Error(), "depth exceeded") {
		t.Errorf("error should mention depth exceeded: %v", err)
	}
}
