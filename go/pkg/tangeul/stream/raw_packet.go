//ff:type feature=tangl type=model
//ff:what rawPacket — one BODY record with its TANGEUL record class, before typed decoding
package stream

import "github.com/geul-org/geul/go/pkg/core"

// rawPacket is one BODY record of a TANGEUL document: its record class (from the
// codec, see classifyRecord) and the decoded record. The reader dispatches on
// Class, so record framing is separated from typed decoding.
type rawPacket struct {
	Class recClass
	Rec   core.Record
}
