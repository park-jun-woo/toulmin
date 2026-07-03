//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what collectStatement — gather one statement's packets (prefix payloads + same-span Node/Edge)
package stream

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// collectStatement gathers the packets of one statement starting at packets[i]:
// the leading no-span Payload run (label/value Text and, in Rules, the condition
// EXPR), then the run of span-bearing Node/Edge packets that share a single
// block span. The Writer emits every statement in exactly this shape, so the
// first differing span, a Payload after the body, or a control marks the next
// statement. It returns the prefix, body and the index just past them.
func collectStatement(packets []rawPacket, i int) (prefix, body []rawPacket, next int, err error) {
	for i < len(packets) && packets[i].Class == word.ClassPayload {
		prefix = append(prefix, packets[i])
		i++
	}
	var span word.Span
	spanSet := false
	for i < len(packets) && (packets[i].Class == word.ClassNode || packets[i].Class == word.ClassEdge) {
		s, e := spanOfRaw(packets[i])
		if e != nil {
			return nil, nil, 0, e
		}
		if !spanSet {
			span, spanSet = s, true
		} else if s != span {
			break
		}
		body = append(body, packets[i])
		i++
	}
	return prefix, body, i, nil
}
