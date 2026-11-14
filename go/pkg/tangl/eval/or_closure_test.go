//ff:func feature=tangl type=engine control=sequence
//ff:what TestOrClosure — tests short-circuit or over empty, first-true, and all-false terms
package eval

import "testing"

func TestOrClosure(t *testing.T) {
	tf := func(w any) evalFunc { return func(*scope) (bool, any) { return true, w } }
	ff := func(w any) evalFunc { return func(*scope) (bool, any) { return false, w } }
	sc := newScope(nil, nil)

	if b, w := orClosure(nil)(sc); b || w != nil {
		t.Fatalf("empty or: want false,nil got %v,%v", b, w)
	}
	if b, w := orClosure([]evalFunc{ff("a"), tf("b"), ff("c")})(sc); !b || w != "b" {
		t.Fatalf("first-true or: want true,b got %v,%v", b, w)
	}
	if b, w := orClosure([]evalFunc{ff("a"), ff("b")})(sc); b || w != "b" {
		t.Fatalf("all-false or: want false with last witness b, got %v,%v", b, w)
	}
}
