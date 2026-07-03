//ff:type feature=tangl type=model
//ff:what stubCodec — a test-only Codec that returns canned statement text for Render tests
package doc

// stubCodec is a minimal test Codec: Locale is configurable so independent tests
// register distinct locales, and RenderStmt returns canned single-line and
// multi-line statements exercising the core renderer's skeleton regeneration.
type stubCodec struct {
	locale string
}
