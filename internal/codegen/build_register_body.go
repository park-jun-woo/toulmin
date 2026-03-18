//ff:func feature=codegen type=codegen control=iteration dimension=1
//ff:what buildRegisterBody — writes all eng.Register calls to buffer
package codegen

import (
	"bytes"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// buildRegisterBody appends eng.Register calls for all metas to buf.
func buildRegisterBody(buf *bytes.Buffer, metas []toulmin.RuleMeta) {
	for _, m := range metas {
		writeRegisterCall(buf, m)
	}
}
