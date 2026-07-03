//ff:type feature=tangl type=model
//ff:what num32Words — the fixed word count of a NUM32 (2 words, big-endian uint32)
package word

// num32Words is the number of 16-bit words in a NUM32 field: a 32-bit unsigned
// integer stored as high word then low word (GEUL §2.4). Used for epoch stamps
// and the 32-bit frame-length extension.
const num32Words = 2
