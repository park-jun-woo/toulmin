//ff:func feature=tangl type=parser control=sequence
//ff:what Doc — convert a scanner Span into the equivalent doc.Span
package scan

import "github.com/park-jun-woo/toulmin/pkg/tangeul/doc"

// Doc converts a scanner Span into the equivalent doc.Span, letting a consuming
// codec hand scanner byte ranges straight to the tiling layer without copying.
func (s Span) Doc() doc.Span {
	return doc.Span{Off: s.Off, Len: s.Len}
}
