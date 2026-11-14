//ff:func feature=cli type=command control=iteration dimension=1
//ff:what detectLocale — infer the authoring locale from the first `## tangl:` section header
package tanglcli

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// detectLocale infers the authoring locale of a text surface from the keyword of
// its first `## tangl:` section header (README §5): the seven section keywords are
// locale specific (EN `Subject`/`Cases`/… <-> KO `주제`/`사안`/…), so the first one
// pins en or ko. A source with no `## tangl:` header, or a header keyword in
// neither locale's set, is an error — encode cannot pick a codec without it.
func detectLocale(src []byte) (string, error) {
	en := map[string]bool{"Subject": true, "See": true, "Definitions": true, "Rules": true, "Cases": true, "Provides": true, "Internal": true}
	ko := map[string]bool{"주제": true, "참조": true, "정의": true, "규칙": true, "사안": true, "업무": true, "내부": true}
	sc := bufio.NewScanner(bytes.NewReader(src))
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		rest, ok := strings.CutPrefix(line, "## tangl:")
		if !ok {
			continue
		}
		keyword := strings.TrimSpace(rest)
		switch {
		case en[keyword]:
			return "en", nil
		case ko[keyword]:
			return "ko", nil
		default:
			return "", fmt.Errorf("cannot determine locale: unknown section keyword %q in first `## tangl:` header", keyword)
		}
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("cannot determine locale: no `## tangl:` section header found")
}
