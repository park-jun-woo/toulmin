//ff:type feature=tangl type=parser
//ff:what KoCodec — the Korean surface codec (doc.Codec) over its own parser and the shared scanner
package ko

// KoCodec is the Korean (`ko`) surface codec. It encodes Korean TANGEUL v0.3 markdown
// into a doc.Doc with its own predicate-final parser (조사 무시 + 문미 키워드 디스패치)
// and the shared pkg/tangl/codec/scan for a byte-faithful Block tiling (README
// contract §3). It fully satisfies doc.Codec: Locale and Encode live here, and the
// mode-2 RenderStmt regeneration (with 받침-driven 조사 selection) is in render_stmt.go.
type KoCodec struct{}
