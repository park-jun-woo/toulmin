//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what readPackets — classify every BODY record into a rawPacket
package stream

import "github.com/geul-org/gogeul/pkg/core"

// readPackets classifies each BODY record (already framed by
// core.SplitDocument) into a rawPacket. An unclassifiable record is an error.
func readPackets(body []core.Record) ([]rawPacket, error) {
	out := make([]rawPacket, 0, len(body))
	for _, r := range body {
		c, err := classifyRecord(r)
		if err != nil {
			return nil, err
		}
		out = append(out, rawPacket{Class: c, Rec: r})
	}
	return out, nil
}
