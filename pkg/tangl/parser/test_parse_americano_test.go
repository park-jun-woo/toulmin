//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseAmericano — parses the spec's robot control example and checks key fields
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestParseAmericano parses the spec's "로봇 행동 제어" example verbatim and
// checks section counts plus the once/with fields it exercises.
func TestParseAmericano(t *testing.T) {
	doc, err := Parse("testdata/americano.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if doc.Subject != "americano" {
		t.Errorf("Subject = %q, want americano", doc.Subject)
	}
	if len(doc.Sees) != 2 {
		t.Fatalf("len(Sees) = %d, want 2", len(doc.Sees))
	}
	if len(doc.Defs) != 1 || doc.Defs[0].Value != "65°C" {
		t.Fatalf("Defs = %+v, want one ConstDef with value 65°C", doc.Defs)
	}
	if doc.Defs[0].SpecRef == nil || doc.Defs[0].SpecRef.Alias != "sensor" || doc.Defs[0].SpecRef.Name != "Temperature" {
		t.Fatalf("Defs[0].SpecRef = %+v, want sensor.Temperature", doc.Defs[0].SpecRef)
	}
	if len(doc.Cases) != 4 {
		t.Fatalf("len(Cases) = %d, want 4", len(doc.Cases))
	}

	place := doc.Cases[0]
	if place.Name != "can place cup" {
		t.Errorf("Cases[0].Name = %q, want 'can place cup'", place.Name)
	}
	if len(place.Nodes) != 2 || len(place.Attacks) != 1 || len(place.Execs) != 2 {
		t.Fatalf("can place cup: nodes=%d attacks=%d execs=%d, want 2/1/2", len(place.Nodes), len(place.Attacks), len(place.Execs))
	}
	doExec := place.Execs[0]
	if doExec.Kind != ast.DoExec || !doExec.Once || doExec.Node != "order received" {
		t.Errorf("place.Execs[0] = %+v, want DoExec once when 'order received'", doExec)
	}
	runExec := place.Execs[1]
	if runExec.Kind != ast.RunExec || runExec.Case != "can pour water" {
		t.Errorf("place.Execs[1] = %+v, want RunExec to 'can pour water'", runExec)
	}

	serve := doc.Cases[3]
	if serve.Name != "can serve coffee" {
		t.Errorf("Cases[3].Name = %q, want 'can serve coffee'", serve.Name)
	}
	if len(serve.Nodes) != 2 || serve.Nodes[1].Role != ast.ExceptRule {
		t.Fatalf("can serve coffee nodes = %+v", serve.Nodes)
	}
	if len(serve.Nodes[1].With) != 1 || serve.Nodes[1].With[0] != "serve temperature" {
		t.Errorf("low temperature node With = %+v, want ['serve temperature']", serve.Nodes[1].With)
	}
	if len(serve.Execs) != 2 || !serve.Execs[0].Once || serve.Execs[1].Once {
		t.Fatalf("can serve coffee execs = %+v, want [once, not-once]", serve.Execs)
	}

	if len(doc.Provides) != 1 || len(doc.Provides[0].Runs) != 1 || doc.Provides[0].Runs[0] != "can place cup" {
		t.Fatalf("Provides = %+v", doc.Provides)
	}
	if len(doc.Internals) != 1 {
		t.Fatalf("len(Internals) = %d, want 1", len(doc.Internals))
	}
	in := doc.Internals[0]
	if in.Kind != ast.EveryTick || in.Interval != "1s" || in.Until != "can serve coffee" {
		t.Errorf("Internals[0] = %+v, want every 1s until 'can serve coffee'", in)
	}
	if len(in.Runs) != 1 || in.Runs[0] != "can place cup" {
		t.Errorf("Internals[0].Runs = %+v, want ['can place cup']", in.Runs)
	}
}
