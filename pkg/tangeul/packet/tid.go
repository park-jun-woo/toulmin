//ff:type feature=tangl type=model
//ff:what TID — a segment-local 16-bit target identifier (§4)
package packet

// TID is a segment-local 16-bit target identifier (§4). Every node reference in
// a TANGEUL document is a TID, issued in declaration order; global identity is
// carried by separate binding edges, not by the TID itself. The packet layer
// treats a TID as an opaque word-sized value and never interprets it.
type TID uint16
