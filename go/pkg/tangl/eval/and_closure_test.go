//ff:func feature=tangl type=engine control=sequence
//ff:what TestAndClosure — tests short-circuit and over empty, first-false, and all-true terms
package eval

import "testing"

func TestAndClosure(t *testing.T) {
	tf := func(w any) evalFunc { return func(*scope) (bool, any) { return true, w } }
	ff := func(w any) evalFunc { return func(*scope) (bool, any) { return false, w } }
	sc := newScope(nil, nil)

	if b, w := andClosure(nil)(sc); !b || w != nil {
		t.Fatalf("empty and: want true,nil got %v,%v", b, w)
	}
	if b, w := andClosure([]evalFunc{tf("a"), ff("b"), tf("c")})(sc); b || w != "b" {
		t.Fatalf("first-false and: want false,b got %v,%v", b, w)
	}
	if b, w := andClosure([]evalFunc{tf("a"), tf("b")})(sc); !b || w != "b" {
		t.Fatalf("all-true and: want true with last witness b, got %v,%v", b, w)
	}
}
