//ff:func feature=tangl type=parser control=sequence
//ff:what parseRunCheckKo — parse an endpoint/internal body '`C` 진행한다' or '`C` 판정한다'
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseRunCheckKo classifies a Provides/Internal body item: 진행한다 (run, slotKind
// 1) or 판정한다 (check, slotKind 2), extracting the backtick-quoted case name. It
// returns ok=false (no error) for anything else (e.g. a 필수다 require handled by the
// caller). Unlike the case-scoped run edge, this body form has no 면 condition node.
func parseRunCheckKo(it scan.Item, path string) (name string, kind int, ok bool, err error) {
	body, isRun := takeSuffixKo(it.Text, "진행한다")
	if isRun {
		kind = 1
	} else if b, isCheck := takeSuffixKo(it.Text, "판정한다"); isCheck {
		body, kind = b, 2
	} else {
		return "", 0, false, nil
	}
	name, rest, found := takeBacktick(body)
	if !found {
		return "", 0, true, errAtKo(path, it.No, "expected backtick-quoted case name")
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return "", 0, true, errAtKo(path, it.No, "unexpected trailing text: %q", rest)
	}
	return name, kind, true, nil
}
