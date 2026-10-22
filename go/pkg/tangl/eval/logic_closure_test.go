//ff:func feature=tangl type=engine control=sequence
//ff:what TestLogicClosure — tests and/or operator selection and the unknown-operator error
package eval

import "testing"

func TestLogicClosure(t *testing.T) {
	tf := func(*scope) (bool, any) { return true, "t" }
	ff := func(*scope) (bool, any) { return false, "f" }
	sc := newScope(nil, nil)

	andFn, err := logicClosure("and", []evalFunc{tf, ff})
	if err != nil {
		t.Fatalf("and: %v", err)
	}
	if b, w := andFn(sc); b || w != "f" {
		t.Fatalf("and closure: want false,f got %v,%v", b, w)
	}
	orFn, err := logicClosure("or", []evalFunc{ff, tf})
	if err != nil {
		t.Fatalf("or: %v", err)
	}
	if b, w := orFn(sc); !b || w != "t" {
		t.Fatalf("or closure: want true,t got %v,%v", b, w)
	}
	if _, err := logicClosure("nand", nil); err == nil {
		t.Fatal("want unknown-operator error")
	}
}
