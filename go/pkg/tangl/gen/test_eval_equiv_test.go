//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestEvalGenEquiv — eval and gen produce identical (bool, evidence) across the ctx table
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/validate"
)

// equivDocSrc is the cross-test document: every extended construct (field
// path, map key, spec ref, bare host value/predicate call, exists over
// slices with binder paths and a binder-named flat Field, and/or/either/not,
// document-local sibling call, is-in/is-empty literals) with no See imports,
// so the generated code plus the runner stubs compile in a throwaway module.
const equivDocSrc = "## tangl:Subject\n- this document is `equiv`\n\n" +
	"## tangl:Rules\n" +
	"- `credit ok` when `credit score` is at least 650\n" +
	"- `trusted author` when `author`.`trustScore` is at least 0.7\n" +
	"- `is internal` when `attributes`.`internal` equals true\n" +
	"- `within budget` when spec `discount`.`minPrice` is at most `basePrice`\n" +
	"- `high nlp` when result of `nlpScore` is at least 0.7\n" +
	"- `spam mail` when `isSpam` holds\n" +
	"- `has coupon` when any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`\n" +
	"- `active promotion` when any of `promotions` as `p` satisfies\n" +
	"  - `p`.`active` equals true\n" +
	"  - and `p`.`discount` is at least 10\n" +
	"- `ip in list` when any of `blocklist` as `ip` satisfies `ip` equals `clientIp`\n" +
	"- `combo` when\n" +
	"  - any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`\n" +
	"  - or `credit score` is at least 700\n" +
	"- `flagged` when\n" +
	"  - `isSpam` holds\n" +
	"  - or `manualFlag` equals true\n" +
	"- `no notes` when\n" +
	"  - not `notes` is not empty\n" +
	"- `vip member` when\n" +
	"  - either\n" +
	"    - `grade` equals \"VIP\"\n" +
	"    - and `overdue` is empty\n" +
	"  - or `staff` equals true\n" +
	"- `meta ok` when `credit ok` holds\n" +
	"- `grade listed` when `grade` is in \"VIP\", \"GOLD\"\n\n" +
	"## tangl:Cases\n" +
	"- in case of `gate`\n" +
	"  - `credit ok` is a general rule\n"

// equivCtxJSON is the shared ctx table: all-true, all-false/evidence
// passthrough, missing-keys/error/non-slice, and boundary/partial-element
// entries, exercising every rule's true and false paths.
const equivCtxJSON = `[
 {"credit score": 700, "authorTrustScore": 0.9, "attributes": {"internal": true}, "basePrice": 50,
  "nlpScore": 0.9, "spamFlag": false, "coupons": [{"minPrice": 100}, {"minPrice": 40}],
  "promotions": [{"active": false, "discount": 20}, {"active": true, "discount": 15}],
  "blocklist": ["1.1.1.1", "2.2.2.2"], "clientIp": "2.2.2.2", "manualFlag": false,
  "notes": "", "grade": "VIP", "overdue": "", "staff": false},
 {"credit score": 600, "authorTrustScore": 0.5, "attributes": {"public": true}, "basePrice": 20,
  "nlpScore": 0.4, "spamFlag": true, "spamEvidence": {"reason": "links"}, "coupons": [{"minPrice": 100}],
  "promotions": [{"active": true, "discount": 5}], "blocklist": ["1.1.1.1"], "clientIp": "9.9.9.9",
  "manualFlag": true, "notes": "hello", "grade": "GOLD", "overdue": "late", "staff": true},
 {"nlpErr": true, "coupons": "not-a-slice", "promotions": [], "blocklist": null},
 {"credit score": 650, "basePrice": 100, "coupons": [{"minPrice": 100}, {"code": "X"}],
  "spamFlag": true, "grade": "SILVER", "staff": false, "overdue": ""}
]`

// TestEvalGenEquiv is the eval≡gen equivalence contract: the reference
// evaluator (eval.CompileDoc) and the generated Go code must return the same
// (bool, evidence) — witness Fact included, compared in canonical string form
// — for every rule of the cross-test document over every ctx-table entry.
// The eval side is also run twice to observe (bool, evidence) determinism.
func TestEvalGenEquiv(t *testing.T) {
	doc, err := parser.ParseSource(equivDocSrc, "equiv.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	if err := validate.Validate(doc); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	evalLines := equivEvalLines(t, doc, equivCtxJSON)
	again := equivEvalLines(t, doc, equivCtxJSON)
	for i := range evalLines {
		if evalLines[i] != again[i] {
			t.Fatalf("eval nondeterminism at line %d:\n%s\n%s", i, evalLines[i], again[i])
		}
	}
	genLines := equivGenLines(t, doc, equivCtxJSON)
	if len(evalLines) != len(genLines) {
		t.Fatalf("line count differs: eval=%d gen=%d", len(evalLines), len(genLines))
	}
	for i := range evalLines {
		if evalLines[i] != genLines[i] {
			t.Errorf("eval≡gen mismatch:\n eval: %s\n  gen: %s", evalLines[i], genLines[i])
		}
	}
}
