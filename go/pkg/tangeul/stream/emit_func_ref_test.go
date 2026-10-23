//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitFuncRef — local (no alias), declared alias, and undeclared alias hard error
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitFuncRef covers emitFuncRef's alias branch: an empty Alias (local
// reference) takes aliasTID 0 without a map lookup, a non-empty Alias present in
// aliasByName resolves, and one absent from aliasByName is a hard error. The
// EncodeNode/labelEdge `if err != nil` guards are not exercised: NodeFlagFunc is
// always set with a non-nil Func and the label Edge always carries a resolved
// DstText TID with no gate flag, so neither encode call can fail for the fixed
// shape this function builds.
func TestEmitFuncRef(t *testing.T) {
	e := testEnc(&tidTable{aliasByName: map[string]toulmin.TID{"bank": 2}})
	span := toulmin.Span{Off: 0, Len: 1}

	t.Run("local", func(t *testing.T) {
		var prefix, body []core.Word
		tid, err := emitFuncRef(e, &prefix, &body, ast.Ref{Name: "withdraw"}, span)
		if err != nil {
			t.Fatalf("emitFuncRef: %v", err)
		}
		if len(body) == 0 || len(prefix) == 0 {
			t.Error("expected non-empty prefix and body")
		}
		_ = tid
	})

	t.Run("declared alias", func(t *testing.T) {
		var prefix, body []core.Word
		tid, err := emitFuncRef(e, &prefix, &body, ast.Ref{Alias: "bank", Name: "withdraw"}, span)
		if err != nil {
			t.Fatalf("emitFuncRef: %v", err)
		}
		if len(body) == 0 || len(prefix) == 0 {
			t.Error("expected non-empty prefix and body")
		}
		_ = tid
	})

	t.Run("undeclared alias", func(t *testing.T) {
		var prefix, body []core.Word
		if _, err := emitFuncRef(e, &prefix, &body, ast.Ref{Alias: "missing", Name: "withdraw"}, span); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
