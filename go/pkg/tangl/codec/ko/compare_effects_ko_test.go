//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what compareEffectsKo — assert effects.Closure agrees over every endpoint of a fixture's KO and EN Docs
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/effects"
)

// compareEffectsKo walks every declared endpoint and asserts effects.Closure over
// the KO-encoded Doc equals the EN-encoded one (error presence alike, then a
// DeepEqual of the entry slices). The endpoint set is read from the EN Doc; the KO
// Doc shares the same AST, so any mismatch means the KO surface leaked a different
// effect graph.
func compareEffectsKo(t *testing.T, name string, dKo, dEn *doc.Doc) {
	for _, ep := range dEn.AST.Provides {
		eKo, errKo := effects.Closure(dKo.AST, ep.Name)
		eEn, errEn := effects.Closure(dEn.AST, ep.Name)
		if (errKo == nil) != (errEn == nil) {
			t.Errorf("%s/%s: effects.Closure error disagree: ko=%v en=%v", name, ep.Name, errKo, errEn)
			continue
		}
		if !reflect.DeepEqual(eKo, eEn) {
			t.Errorf("%s/%s: effects.Closure differs between KO and EN encodings", name, ep.Name)
		}
	}
}
