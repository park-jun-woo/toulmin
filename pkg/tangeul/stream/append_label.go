//ff:func feature=tangl type=util control=sequence
//ff:what appendLabel — issue a Text payload for s and append a label Edge binding it to a node
package stream

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// appendLabel issues a Text payload for s onto the statement prefix and appends a
// label Edge (src -> that Text) onto body, the standard way a node's spelling is
// bound losslessly (README §6). The Text rides in the prefix so grouping stays
// prefix-then-body.
func appendLabel(e *enc, pre, body *[]word.Word, src uint16, s string, span word.Span) error {
	w, err := labelEdge(src, addText(e, pre, s), span)
	if err != nil {
		return err
	}
	*body = append(*body, w...)
	return nil
}
