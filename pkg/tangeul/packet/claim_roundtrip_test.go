//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestClaimRoundTrip — each Claim object kind round-trips and the §10 claim matches its words
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestClaimRoundTrip checks each Claim object kind (TID / NUM / SPAN) encodes to
// the expected word count and decodes back identically, and that the §10 balance
// sufficient claim lands on exactly the documented words.
func TestClaimRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		claim Claim
		words int
	}{
		{"tid", Claim{Subject: 0x0011, Predicate: 0x7400, Kind: ObjTID, ObjTID: 0x0012, Mode: 0}, 4},
		{"num", Claim{Subject: 0x0005, Predicate: 0x0107, Kind: ObjNum, ObjNum: 0x0001_2345, Mode: 2}, 5},
		{"span", Claim{Kind: ObjSpan, ObjSpan: word.Span{TID: 1, Off: 9, Len: 4}, Mode: 0}, 6},
	}
	for _, c := range cases {
		ws, err := encodeClaim(c.claim)
		if err != nil {
			t.Fatalf("%s: encodeClaim: %v", c.name, err)
		}
		if len(ws) != c.words {
			t.Fatalf("%s: encodeClaim len = %d, want %d", c.name, len(ws), c.words)
		}
		got, n, err := decodeClaim(ws, c.claim.Kind)
		if err != nil || n != c.words || !reflect.DeepEqual(got, c.claim) {
			t.Fatalf("%s: decodeClaim = %v, %d, %v; want %v", c.name, got, n, err, c.claim)
		}
	}
	golden, _ := encodeClaim(cases[0].claim)
	want := []word.Word{0x0011, 0x7400, 0x0012, 0x0000}
	if !reflect.DeepEqual(golden, want) {
		t.Fatalf("§10 claim words = %v, want %v", golden, want)
	}
}
