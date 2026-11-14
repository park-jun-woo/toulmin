//ff:func feature=cli type=command control=sequence
//ff:what fakeDecodeCodec.Encode — unused stub; runDecode only exercises RenderStmt/Locale
package tanglcli

import "github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"

// Encode is unused by runDecode's mode-2 path and always reports a nil Doc and
// nil error.
func (f fakeDecodeCodec) Encode(src []byte, path string) (*doc.Doc, error) { return nil, nil }
