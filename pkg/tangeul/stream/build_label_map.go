//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what buildLabelMap — group a statement's label edges into per-node ordered label lists
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// buildLabelMap resolves each label edge's Text and groups the strings by their
// src node TID, preserving emission order (name first, then a field's type). A
// label edge pointing at an unknown Text is a hard error.
func buildLabelMap(d *dec, labelEdges []packet.EdgePacket) (map[uint16][]string, error) {
	labels := map[uint16][]string{}
	for _, le := range labelEdges {
		txt, ok := d.texts[uint16(le.DstTID)]
		if !ok {
			return nil, errors.New("tangeul/stream: label edge references unknown Text")
		}
		labels[uint16(le.Src)] = append(labels[uint16(le.Src)], txt)
	}
	return labels, nil
}
