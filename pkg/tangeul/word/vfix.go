//ff:type feature=tangl type=model
//ff:what VFix — signed 16-bit fixed-point value (scale 10000) for verdict/qualifier/confidence
package word

// VFix is a signed 16-bit fixed-point number with scale 10000 (VFIX16, GEUL
// §3.1). The stored int16 equals value*10000, so the documented value range
// [-1.0, +1.0] maps to [-10000, +10000]; the int16 headroom (±3.2767) leaves
// room for intermediate computation. Examples: +1.0=0x2710, 0.0=0x0000,
// -1.0=0xD8F0, 0.8=0x1F40.
type VFix int16

// vfixScale is the fixed-point scale (four decimal digits) of a VFix.
const vfixScale = 10000
