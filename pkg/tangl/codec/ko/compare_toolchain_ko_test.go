//ff:func feature=tangl type=parser control=sequence
//ff:what compareToolchainKo — assert validate/Lint/gen agree on one fixture's KO and EN Docs, then defer to effects
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/gen"
	"github.com/park-jun-woo/toulmin/pkg/tangl/validate"
)

// compareToolchainKo encodes fixture name's KO and EN sources (encodePairKo) and
// asserts the locale-blind toolchain agrees: Validate accepts or rejects both
// alike, Lint yields the same warning count, and gen.Generate emits byte-identical
// Go. The per-endpoint effects comparison is delegated to compareEffectsKo.
func compareToolchainKo(t *testing.T, name string) {
	dKo, dEn := encodePairKo(t, name)
	if (validate.Validate(dKo.AST) == nil) != (validate.Validate(dEn.AST) == nil) {
		t.Errorf("%s: Validate disagree between KO and EN encodings", name)
	}
	if len(validate.Lint(dKo.AST)) != len(validate.Lint(dEn.AST)) {
		t.Errorf("%s: Lint warning count differs between KO and EN encodings", name)
	}
	genKo, errKo := gen.Generate(dKo.AST)
	genEn, errEn := gen.Generate(dEn.AST)
	if (errKo == nil) != (errEn == nil) {
		t.Errorf("%s: gen.Generate error disagree between KO and EN encodings", name)
	}
	if errKo == nil && errEn == nil && genKo != genEn {
		t.Errorf("%s: gen.Generate output differs between KO and EN encodings", name)
	}
	compareEffectsKo(t, name, dKo, dEn)
}
