//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestResolveFuncRef — unknown Func TID errors; zero alias stays local; non-zero alias resolves via nodeName
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestResolveFuncRef drives resolveFuncRef directly (white-box) over its three
// branches: an unregistered Func node TID (the funcRef lookup misses — hard
// error), a registered Func node with alias 0 (a local reference — the
// nodeName lookup is skipped, Alias stays ""), and a registered Func node with
// a non-zero alias TID (Alias resolves through d.nodeName).
func TestResolveFuncRef(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(d *dec)
		fTID    toulmin.TID
		wantErr bool
		want    ast.Ref
	}{
		{
			name:    "unknown Func node TID is a hard error",
			setup:   func(d *dec) {},
			fTID:    99,
			wantErr: true,
		},
		{
			name: "zero alias stays a local reference",
			setup: func(d *dec) {
				d.funcRef[5] = funcInfo{alias: 0, name: "withdraw"}
			},
			fTID: 5,
			want: ast.Ref{Alias: "", Name: "withdraw"},
		},
		{
			name: "non-zero alias resolves through nodeName",
			setup: func(d *dec) {
				d.funcRef[7] = funcInfo{alias: 3, name: "log"}
				d.nodeName[3] = "bank"
			},
			fTID: 7,
			want: ast.Ref{Alias: "bank", Name: "log"},
		},
	}

	check := func(t *testing.T, setup func(d *dec), fTID toulmin.TID, wantErr bool, want ast.Ref) {
		d := newDec(testWire)
		setup(d)
		got, err := resolveFuncRef(d, fTID)
		if wantErr && err == nil {
			t.Fatalf("expected error, got nil (ref=%+v)", got)
		}
		if wantErr {
			return
		}
		if err != nil {
			t.Fatalf("resolveFuncRef: %v", err)
		}
		if got != want {
			t.Errorf("ref = %+v, want %+v", got, want)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { check(t, tt.setup, tt.fTID, tt.wantErr, tt.want) })
	}
}
