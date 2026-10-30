//ff:func feature=tangl type=codegen control=sequence
//ff:what TestSnapshotLists — materialized list Specs become plain name → values snapshots
package gen

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestSnapshotLists(t *testing.T) {
	got := snapshotLists(map[string]toulmin.StringListSpec{
		"domains": toulmin.NewStringListSpec("domains", []string{"a,b", "c"}),
		"empty":   toulmin.NewStringListSpec("empty", nil),
	})
	if len(got) != 2 || !reflect.DeepEqual(got["domains"], []string{"a,b", "c"}) || len(got["empty"]) != 0 {
		t.Fatalf("snapshotLists = %#v", got)
	}
	if got := snapshotLists(nil); len(got) != 0 {
		t.Fatalf("nil input = %#v, want empty", got)
	}
}
