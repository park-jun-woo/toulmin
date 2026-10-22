//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what buildLabelMap — group a statement's label edges into per-node ordered label lists
package stream

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// buildLabelMap resolves each label edge's Text and groups the strings by their
// src node TID, preserving emission order (name first, then a field's type). A
// label edge pointing at an unknown Text is a hard error.
func buildLabelMap(d *dec, labelEdges []toulmin.EdgePacket) (map[toulmin.TID][]string, error) {
	labels := map[toulmin.TID][]string{}
	for _, le := range labelEdges {
		txt, ok := d.texts[toulmin.TID(le.DstTID)]
		if !ok {
			return nil, errors.New("tangeul/stream: label edge references unknown Text")
		}
		labels[toulmin.TID(le.Src)] = append(labels[toulmin.TID(le.Src)], txt)
	}
	return labels, nil
}
