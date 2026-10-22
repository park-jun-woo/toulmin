//ff:type feature=tangl type=model
//ff:what recClass — the reader's grouping class of a TANGEUL record (node, link, payload, control, source)
package stream

// recClass groups TANGEUL records the way the reader consumes them: span-bearing
// NODE and LINK records form a statement body, TEXT and PREDICATE-DECLARE
// records are the statement prefix, framing controls stand alone, and SOURCE
// opens the BODY.
type recClass uint8

const (
	classNode recClass = iota
	classLink
	classPayload
	classControl
	classSource
)
