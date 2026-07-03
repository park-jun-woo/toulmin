//ff:type feature=tangl type=model
//ff:what TextFormat — the format code of a Text payload (§7.4); only UTF-8 in v0.1
package packet

// TextFormat is the format code written as word 2 of a Text payload (§7.4). The
// L1 format codebook is not yet formalized (§11), so v0.1 emits and accepts only
// UTF-8; the decoder rejects any other code rather than guessing an encoding.
type TextFormat uint16

// TextFormatUTF8 is the UTF-8 payload format code.
const TextFormatUTF8 TextFormat = 0x0001
