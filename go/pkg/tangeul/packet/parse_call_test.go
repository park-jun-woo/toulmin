//ff:func feature=tangl type=util control=sequence
//ff:what TestParseCall — parseCall rebuilds an ast.Call from [alias][name] and rejects truncated/out-of-range records
package packet

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestParseCall exercises parseCall's successful rebuild — an alias and a name
// label index resolved through labels, advancing *pos with node count 1 — plus
// its two guards: a truncated record and an out-of-range alias or name label
// index.
func TestParseCall(t *testing.T) {
	labels := []string{"", "nlp", "score"}

	pos := 0
	recs := []core.Word{1, 2}
	got, count, err := parseCall(recs, &pos, labels)
	if err != nil || count != 1 || pos != 2 {
		t.Fatalf("Call: got %#v, count %d, pos %d, err %v", got, count, pos, err)
	}
	if want := (ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}); !reflect.DeepEqual(got, want) {
		t.Fatalf("Call = %#v, want %#v", got, want)
	}

	pos = 0
	if _, _, err := parseCall([]core.Word{1}, &pos, labels); err == nil {
		t.Error("truncated Call record: expected error")
	}

	pos = 0
	if _, _, err := parseCall([]core.Word{9, 0}, &pos, labels); err == nil {
		t.Error("out-of-range alias label: expected error")
	}

	pos = 0
	if _, _, err := parseCall([]core.Word{0, 9}, &pos, labels); err == nil {
		t.Error("out-of-range name label: expected error")
	}
}
