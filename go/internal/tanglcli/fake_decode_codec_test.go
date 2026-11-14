//ff:type feature=cli type=model
//ff:what fakeDecodeCodec — a test-only doc.Codec used to exercise runDecode's mode-2 (--locale) branches
package tanglcli

// fakeDecodeCodec is a minimal doc.Codec used only to exercise runDecode's
// mode-2 (--locale) branches: a found codec whose RenderStmt either succeeds or
// fails. Encode is unused by runDecode.
type fakeDecodeCodec struct {
	locale string
	fail   bool
}
