//ff:func feature=tangl type=engine control=sequence
//ff:what TestDocResolverValue — tests straight value-resolution delegation to the host
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestDocResolverValue(t *testing.T) {
	host := fakeResolver{values: map[string]func(toulmin.Context, toulmin.Specs) (any, error){
		"\x00v": func(toulmin.Context, toulmin.Specs) (any, error) { return 7, nil },
	}}
	d := docResolver{host: host, ruleSet: map[string]bool{"v": true}}

	fn, ok := d.Value("", "v")
	if !ok {
		t.Fatal("value should delegate to the host even for a rule name")
	}
	if v, err := fn(testCtx(), nil); err != nil || v != 7 {
		t.Fatalf("host value: got %v,%v", v, err)
	}
	if _, ok := d.Value("", "missing"); ok {
		t.Fatal("missing host value should not resolve")
	}
}
