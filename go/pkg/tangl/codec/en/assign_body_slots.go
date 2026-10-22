//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what assignBodySlots — assign run/check StmtRefs to a body's items by per-kind counters
package en

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// assignBodySlots walks an endpoint/internal body in source order, classifying each
// item by parseRunCheck (slotKind 1 Run, 2 Check) and assigning the next per-kind
// counter as slotIdx. The backtick name is checked against the AST slice value; a
// mismatch, overflow, or malformed line is a scanner/parser drift reported as error.
func assignBodySlots(children []scan.Item, code toulmin.SectionCode, idx int, runs, checks []string, refByLine map[int]doc.StmtRef) error {
	runN, checkN := 0, 0
	for _, ch := range children {
		if ch.Struck {
			continue
		}
		slotKind, name, isBody, err := parseRunCheck(ch.Text)
		if err != nil {
			return fmt.Errorf("tangl/codec/en: line %d: %w", ch.No, err)
		}
		if !isBody {
			continue
		}
		slice, slotIdx := checks, checkN
		if slotKind == 1 {
			slice, slotIdx = runs, runN
		}
		if err := verifyName(slice, slotIdx, name); err != nil {
			return fmt.Errorf("tangl/codec/en: line %d: %w", ch.No, err)
		}
		refByLine[ch.No] = doc.StmtRef{Sec: code, Path: []int{idx, slotKind, slotIdx}}
		runN, checkN = bumpCounter(slotKind, runN, checkN)
	}
	return nil
}
