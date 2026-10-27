//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what collectStatement — gather one statement: its leading payload records, then same-span NODE/LINK records
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
)

// collectStatement gathers the records of one statement starting at packets[i]:
// the leading span-less payload run (TEXT and, in Rules, the PREDICATE-DECLARE),
// then the run of span-bearing NODE/LINK records that share a single block
// span. The Writer emits every statement in exactly this shape, so the first
// differing span, a payload after the body, or a control marks the next
// statement. It returns the prefix, body and the index just past them.
func collectStatement(d *dec, packets []rawPacket, i int) (prefix, body []rawPacket, next int, err error) {
	for i < len(packets) && packets[i].Class == classPayload {
		prefix = append(prefix, packets[i])
		i++
	}
	var span toulmin.Span
	spanSet := false
	for ; i < len(packets) && (packets[i].Class == classNode || packets[i].Class == classLink); i++ {
		s, e := spanOfRaw(d, packets[i])
		if e != nil {
			return nil, nil, 0, e
		}
		if spanSet && s != span {
			break
		}
		span, spanSet = s, true
		body = append(body, packets[i])
	}
	return prefix, body, i, nil
}
