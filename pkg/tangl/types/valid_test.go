//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestValid — table-driven check of Valid dispatching to each TANGL basic type validator
package types

import "testing"

func TestValid(t *testing.T) {
	cases := []struct {
		name     string
		typeName string
		value    any
		want     bool
	}{
		{"text ok", Text, "hello", true},
		{"text wrong kind", Text, 42, false},

		{"integer ok int", Integer, 42, true},
		{"integer ok int64", Integer, int64(42), true},
		{"integer wrong kind float", Integer, 42.5, false},

		{"number ok float64", Number, 3.14, true},
		{"number ok float32", Number, float32(3.14), true},
		{"number wrong kind int", Number, 3, false},

		{"boolean ok true", Boolean, true, true},
		{"boolean ok false", Boolean, false, true},
		{"boolean wrong kind", Boolean, "true", false},

		{"email ok", Email, "user@example.com", true},
		{"email missing at", Email, "userexample.com", false},
		{"email missing domain dot", Email, "user@examplecom", false},
		{"email wrong kind", Email, 123, false},

		{"date ok", Date, "2026-07-01", true},
		{"date bad format", Date, "07/01/2026", false},
		{"date wrong kind", Date, 20260701, false},

		{"time ok hms", Time, "13:45:00", true},
		{"time ok hm", Time, "13:45", true},
		{"time bad format", Time, "1pm", false},

		{"url ok", URL, "https://example.com/path", true},
		{"url no scheme", URL, "example.com/path", false},
		{"url wrong kind", URL, 123, false},

		{"currency ok int", Currency, 1000, true},
		{"currency ok float", Currency, 19.99, true},
		{"currency wrong kind", Currency, "19.99", false},

		{"quantity ok int", Quantity, 5, true},
		{"quantity ok float", Quantity, 5.5, true},
		{"quantity wrong kind", Quantity, "5", false},

		{"unknown type name", "Frobnicate", "anything", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Valid(c.typeName, c.value)
			if got != c.want {
				t.Errorf("Valid(%q, %#v) = %v, want %v", c.typeName, c.value, got, c.want)
			}
		})
	}
}
