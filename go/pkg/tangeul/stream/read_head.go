//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what readHead — validate a TANGEUL HEAD (spec, profile, language) and extract its settings
package stream

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/core"
)

// readHead applies each HEAD Meta record via readHeadRecord and requires both
// SPEC-VERSION and PROFILE, so a document of another spec or profile is
// rejected before its BODY is read.
func readHead(head []core.Record) (headInfo, error) {
	var h headInfo
	for _, r := range head {
		if err := readHeadRecord(&h, r); err != nil {
			return headInfo{}, err
		}
	}
	if !h.spec || !h.profile {
		return headInfo{}, errors.New("tangeul/stream: HEAD lacks SPEC-VERSION or PROFILE")
	}
	return h, nil
}
