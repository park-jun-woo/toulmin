//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what tidKey — the canonical string key of a StmtRef (section code + index path)
package stream

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// tidKey renders ref into a canonical, collision-free string so the tidTable can
// map a StmtRef coordinate to its primary entity TID. The section code and each
// path index are joined with a separator that cannot occur inside a decimal
// integer, so distinct coordinates never alias.
func tidKey(ref doc.StmtRef) string {
	var b strings.Builder
	b.WriteString(strconv.FormatUint(uint64(ref.Sec), 10))
	for _, p := range ref.Path {
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(p))
	}
	return b.String()
}
