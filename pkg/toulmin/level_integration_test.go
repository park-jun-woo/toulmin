//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestLevelIntegration — golden: role default Levels, explicit override/opt-out, Fact roundtrip, determinism
package toulmin

import "testing"

func TestLevelIntegration(t *testing.T) {
	g := buildLevelGraph()
	_, tr, err := g.Run(NewContext())
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	wantLevel := map[string]Level{
		"lvlWarrant":  Unset,  // role rule → default Unset
		"lvlCounter":  Fail,   // role counter → default Fail
		"lvlExcept":   Review, // role except → default Review
		"lvlOverride": Review, // explicit .Level(Review) overrides counter's Fail default
		"lvlOptOut":   Unset,  // explicit .Level(Unset) opts out, does not fall back to Fail
	}
	for name, want := range wantLevel {
		entry, ok := tr.Get(name)
		if !ok {
			t.Fatalf("trace missing entry %q", name)
		}
		if entry.Level != want {
			t.Errorf("%s Level = %v, want %v", name, entry.Level, want)
		}
	}

	// Fact roundtrips through Evidence unchanged.
	factEntry, ok := tr.Get("lvlFact")
	if !ok {
		t.Fatal("trace missing lvlFact entry")
	}
	fact, ok := factEntry.Evidence.(Fact)
	if !ok {
		t.Fatalf("lvlFact Evidence = %T, want Fact", factEntry.Evidence)
	}
	if fact.Where != "field.x" || fact.Expected != 10 || fact.Actual != 3 {
		t.Errorf("Fact roundtrip mismatch: %+v", fact)
	}

	// Determinism: a second Run yields identical Level + Evidence per node.
	_, tr2, err := g.Run(NewContext())
	if err != nil {
		t.Fatalf("second Run error: %v", err)
	}
	for _, e1 := range tr.All() {
		e2, ok := tr2.Get(e1.Name)
		if !ok {
			t.Fatalf("second trace missing %q", e1.Name)
		}
		if e1.Level != e2.Level {
			t.Errorf("%s Level not deterministic: %v vs %v", e1.Name, e1.Level, e2.Level)
		}
		if e1.Evidence != e2.Evidence {
			t.Errorf("%s Evidence not deterministic: %v vs %v", e1.Name, e1.Evidence, e2.Evidence)
		}
	}
}
