//ff:func feature=tangl type=engine control=sequence
//ff:what TestStepInto — tests map-key, struct-field, pointer, and miss step resolution
package eval

import "testing"

func TestStepInto(t *testing.T) {
	if _, ok := stepInto(nil, "k"); ok {
		t.Fatal("nil value should not resolve")
	}
	if v, ok := stepInto(map[string]any{"k": 1}, "k"); !ok || v != 1 {
		t.Fatalf("map key: got %v,%v", v, ok)
	}
	if _, ok := stepInto(map[string]any{}, "k"); ok {
		t.Fatal("missing map key should not resolve")
	}
	if _, ok := stepInto((*testAuthor)(nil), "trustScore"); ok {
		t.Fatal("nil pointer should not resolve")
	}
	if v, ok := stepInto(&testAuthor{TrustScore: 0.8}, "trustScore"); !ok || v != 0.8 {
		t.Fatalf("pointer struct field: got %v,%v", v, ok)
	}
	if v, ok := stepInto(testAuthor{TrustScore: 0.5}, "TrustScore"); !ok || v != 0.5 {
		t.Fatalf("struct field: got %v,%v", v, ok)
	}
	if _, ok := stepInto("scalar", "k"); ok {
		t.Fatal("non-struct value should not resolve")
	}
	if _, ok := stepInto(testAuthor{}, "ghost"); ok {
		t.Fatal("missing field should not resolve")
	}
}
