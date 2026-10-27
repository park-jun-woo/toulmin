# pkg/tangeul

TANGEUL is the **binary single source of truth (SSOT)** for the `tangl`
toolchain. An authored markdown surface (`.md`) is *absorbed* into a compact,
locale-independent `.geul` file — exactly one GEUL DOC (GEUL t002) whose records
follow the TOULMIN codec (`geul/draft/docs/TOULMIN.md`); `check`/`ast`/`effects`/`gen` all run on
that stream, and a decoder renders it back to text. Direct execution of `.md`
is gone — the surface is reachable only through an encoder/decoder codec.

```
 EN .md ──[EN codec]──┐                        ┌──[mode 1]── authored text verbatim (byte-identical)
                      ├──> .geul (SSOT) ───────┤
 KO .md ──[KO codec]──┘          │             └──[mode 2]── EN/KO normalized rendering
                                 ▼
                     tangl loader: Load → doc.Doc{ AST, Source, Blocks }
                                 │
                  ┌──────────────┼──────────────┐
                  ▼              ▼              ▼
               validate       effects         gen      ← unchanged (consume *ast.Document)
```

## Design

`pkg/tangl/ast` is **not modified**: `validate`/`effects`/`gen` keep consuming
`*ast.Document` unchanged. Prose losslessness lives in a separate container.

```go
// pkg/tangeul/doc
type Doc struct {
    AST    *ast.Document // graph (existing toolchain input)
    Source []byte        // Source Text — authored bytes verbatim
    Blocks []Block       // stream-order tiling: prose spans + statement anchor spans
    Locale string        // authoring locale ("en"/"ko"/"")
}

type Block struct {
    Kind BlockKind // ProseBlock | StmtBlock | SectionBlock
    Span Span      // [Off, Off+Len) into Doc.Source
    Ref  StmtRef   // AST coordinate (SectionCode + variable-length index path) for a StmtBlock
}
```

The `Blocks` slice is the single source of truth for element order (since
`ast.Case` stores Requires/Nodes/Attacks/Execs in separate slices). `VerifyTiling`
is the integrity gate: the blocks must tile `Source` exactly, and it runs on every
`Load` and every encoder output.

## Binary format

A `.geul` file is one GEUL t002 document (TOULMIN codec):

```
DOC-1WORD/2WORD/4WORD   total length
  SPEC-VERSION          0x0002
  PROFILE               "tangeul-<codebook major>.<minor>"
  [DEFAULT-LANGUAGE]    authoring locale (en/ko)
BODY
  SOURCE                authored bytes verbatim + SPAN width
  records …             declaration order (NODE / LINK / TEXT / framing)
END
```

- Every record is a GEUL Edge record (TAG WORD `000K KLLL LCCC CCCC`) under a
  direct TOULMIN codec (`0x60`–`0x76`: declarations, ATTACKS, SUPERSEDES, RUN,
  DO, DO-ONCE, UNDO, CHECKING, ENTRY, TRIGGER) or a GEUL ISSUER format number
  (`0xC000`+: SOURCE, TEXT, SECTION, CASE-BEGIN/END, PROSE, NAMESPACE, FUNC,
  FIELD, ALIAS, GROUP and the remaining relations).
- Every record carries an Edge Id; references are Edge Ids. Both are embedded
  LOCAL Reference records (`0000 0LLL L100 0000` + id in minimal words), so no
  id width is declared. The SPAN width is declared in SOURCE. `with`/`grounded
  on` list counts are embedded UNSIGNED-INTEGER Values.
- Verdict/qualifier are SIGNED-DECIMAL Values, the confidence gate an operator
  word plus an UNSIGNED-DECIMAL percent, strings UTF8-TEXT Values.
- The former Free-area `.tangeul` framing (STREAM_HEADER, 16-bit TID, VFIX16,
  GATE16) is retired; the codebook major moved to 1.

## Package layout

| Package | Role |
|---|---|
| `pkg/tangeul` | loader: `Load`/`Save` (`.geul` ↔ `doc.Doc`), `EqualIgnoringLine` AST comparison |
| `github.com/geul-org/gogeul/pkg/core` | GEUL t002 core (imported, not part of this repo): TAG WORD/Length/ISSUER record framing, DOC container (wrap/split), Values (UTF-8, integer, decimal), declared-width ids |
| `github.com/geul-org/gogeul/pkg/toulmin` | TOULMIN codec (imported): code tables, codec/format-number mapping, VFix/Gate/Span, NODE/LINK/control/TEXT/SOURCE records |
| `pkg/tangeul/codebook` | tangl AST ↔ TOULMIN code mappings (role, level, operators, field types, exec kinds) |
| `pkg/tangeul/packet` | EXPR: the PREDICATE-DECLARE condition record (tangl expression AST) |
| `pkg/tangeul/stream` | Doc ↔ GEUL DOC serialization (HEAD, Writer/Reader, Edge Id issue/resolve, order preservation) |
| `pkg/tangeul/doc` | `Doc` container + `VerifyTiling` + the surface `Codec` interface + `Register/LookupCodec` |

Surface codecs live outside this tree: the EN codec is `pkg/tangl/codec/en`
(`ENCodec.Encode`/`.Decode`); the KO codec (`pkg/tangl/codec/ko`) lands with
`plans/tangl-ko`. The CLI wires them in `internal/tanglcli/register_codecs.go`.

## Key API

| Symbol | Signature | Purpose |
|---|---|---|
| `tangeul.Load` | `Load(path string) (*doc.Doc, error)` | read a `.geul` file (extension check → `stream.ReadDoc` → `VerifyTiling`); stamps `AST.Path` |
| `tangeul.Save` | `Save(d *doc.Doc, path string) error` | write a `Doc` to a `.geul` file |
| `tangeul.EqualIgnoringLine` | `EqualIgnoringLine(a, b *ast.Document) bool` | AST equality ignoring `Line` (streams don't store authored line numbers); reused by the codec round-trip checks |
| `doc.RegisterCodec` / `doc.LookupCodec` | — | register / find a surface `Codec` by locale |
| `doc.VerifyTiling` | `VerifyTiling(d *doc.Doc) error` | the tiling integrity gate |

## Decode modes

- **Mode 1 (verbatim)** — replays `Source` bytes exactly; `decode` output is
  byte-identical to the authored `.md`.
- **Mode 2 (regenerated)** — `decode --locale en|ko` re-renders statements via
  the target locale's codec (`StmtBlock`/`SectionBlock` regenerated; `ProseBlock`
  replayed verbatim).

## Verification (Phase 008)

- **Step-1 golden fixtures** (`testdata/{transfer,access}.geul`) are
  hand-assembled with placeholder prose spans — a pre-encoder stand-in enough for
  AST/validate/effects/gen parity, deliberately not for byte identity. The parity
  tests (`Test*Parity`) assert `Load(golden).AST == parser.Parse(.md)` up to `Line`.
- **Step-2 cross-check** (`TestEncoderGoldenXCheck`) confirms the tangl-en encoder
  now produces a stream structurally equivalent to the step-1 golden: same AST (up
  to `Line`) and the same ordered sequence of statement/section blocks (`StmtRef`
  coordinates). Prose spans differ by design (placeholder vs real), so byte equality
  is *not* the criterion; mode-1 byte identity is verified on the encoder's own
  output by the CLI e2e (`TestEncodeDecodeEN` in `internal/tanglcli`).
