//ff:func feature=tangl type=parser control=sequence
//ff:what TestSpanDoc — tests Span.Doc converts to an equal-valued doc.Span
package scan

import "testing"

func TestSpanDoc(t *testing.T) {
	s := Span{Off: 12, Len: 34}
	d := s.Doc()
	if d.Off != s.Off || d.Len != s.Len {
		t.Errorf("Doc() = %+v, want Off=%d Len=%d", d, s.Off, s.Len)
	}
}
