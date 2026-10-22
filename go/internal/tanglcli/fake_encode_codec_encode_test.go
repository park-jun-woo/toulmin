//go:build fakecodec

//ff:func feature=cli type=command control=selection
//ff:what fakeEncodeCodec.Encode — a forced error, or a tiling-valid/tiling-broken transfer fixture Doc
package tanglcli

import "github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"

// Encode returns f.encodeErr when set (driving runEncode's Encode-error branch);
// otherwise it returns the transfer fixture Doc, with its Blocks emptied when
// f.badTiling is set so doc.VerifyTiling's trailing-gap check fails (driving
// runEncode's VerifyTiling-error branch).
func (f *fakeEncodeCodec) Encode(src []byte, path string) (*doc.Doc, error) {
	switch {
	case f.encodeErr != nil:
		return nil, f.encodeErr
	case f.badTiling:
		d := buildTransferDoc()
		d.Blocks = nil
		return d, nil
	default:
		return buildTransferDoc(), nil
	}
}
