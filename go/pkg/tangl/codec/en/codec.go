//ff:type feature=tangl type=parser
//ff:what ENCodec — the English surface codec (doc.Codec) over the shared parser and scanner
package en

// ENCodec is the English (`en`) surface codec. It encodes English TANGEUL v0.3
// markdown into a doc.Doc by reusing pkg/tangl/parser for the AST and the shared
// pkg/tangl/codec/scan for a byte-faithful Block tiling (README contract §3). The
// mode-2 RenderStmt regeneration is added by a later phase; Locale and Encode live
// here.
type ENCodec struct{}
