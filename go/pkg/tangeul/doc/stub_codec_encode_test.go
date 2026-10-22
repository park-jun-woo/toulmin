//ff:func feature=tangl type=util control=sequence
//ff:what stubCodec.Encode — unused in Render tests; reports it is a stub
package doc

import "errors"

// Encode is unused by the doc-layer tests (encoding belongs to tangl-en/ko) and
// reports that it is a stub.
func (s stubCodec) Encode(src []byte, path string) (*Doc, error) {
	return nil, errors.New("stubCodec: Encode is a test stub")
}
