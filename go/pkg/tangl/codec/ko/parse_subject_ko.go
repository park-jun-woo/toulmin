//ff:func feature=tangl type=parser control=sequence
//ff:what parseSubjectKo — parse '이 문서는 `name` 이다' into the subject namespace
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseSubjectKo parses the tangl:주제 statement "이 문서는 `name` 이다", returning the
// subject namespace. The sentence-final 이다 is consumed first so the shared
// dispatch never confuses it with the 이 particle.
func parseSubjectKo(it scan.Item, path string) (string, error) {
	body, ok := takeSuffixKo(it.Text, "이다")
	if !ok {
		return "", errAtKo(path, it.No, "expected '이 문서는 `name` 이다', got %q", it.Text)
	}
	rest, ok := takeKeywordKo(body, "이 문서는")
	if !ok {
		return "", errAtKo(path, it.No, "expected '이 문서는' prefix, got %q", it.Text)
	}
	name, rest, ok := takeBacktick(rest)
	if !ok {
		return "", errAtKo(path, it.No, "expected backtick-quoted subject name")
	}
	if strings.TrimSpace(rest) != "" {
		return "", errAtKo(path, it.No, "unexpected trailing text: %q", rest)
	}
	return name, nil
}
