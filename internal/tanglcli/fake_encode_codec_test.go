//go:build fakecodec

//ff:type feature=cli type=model
//ff:what fakeEncodeCodec — a test-only, mutable doc.Codec used to exercise runEncode's post-lookup branches
package tanglcli

// fakeEncodeCodec is a minimal, mutable doc.Codec used only to exercise
// runEncode's post-lookup branches (Encode success/error, VerifyTiling
// success/error) once it is registered for locale "ko" (the real en codec owns
// "en"). Its fields are mutated between subtests — doc.RegisterCodec panics on a
// duplicate locale, so it is registered exactly once and reused. RenderStmt is
// unused by runEncode.
type fakeEncodeCodec struct {
	encodeErr error
	badTiling bool
}
