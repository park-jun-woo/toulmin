//ff:func feature=tangl type=util control=selection
//ff:what classifyRecord — the reader class of a BODY record from its TOULMIN codec or GEUL format number
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// classifyRecord returns the reader class of r. A record outside the TANGEUL
// profile (another kind, another Issuer, an unassigned codec) is an error: the
// profile defines every BODY record.
func classifyRecord(r core.Record) (recClass, error) {
	ref, err := toulmin.RecordRef(r)
	if err != nil {
		return 0, err
	}
	_, _, _, isNode := toulmin.NodeFromCodec(ref)
	_, _, isLink := toulmin.PropFromCodec(ref)
	_, isControl := toulmin.ControlFromFormat(ref.Format)
	switch {
	case isNode:
		return classNode, nil
	case isLink:
		return classLink, nil
	case ref.Codec == toulmin.CodecPredicateDeclare || ref.Format == toulmin.FormatText:
		return classPayload, nil
	case isControl:
		return classControl, nil
	case ref.Format == toulmin.FormatSource:
		return classSource, nil
	default:
		return 0, fmt.Errorf("tangeul/stream: record 0x%02X/0x%04X is not in the TANGEUL profile", ref.Codec, ref.Format)
	}
}
