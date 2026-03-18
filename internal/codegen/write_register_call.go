//ff:func feature=codegen type=codegen control=sequence
//ff:what writeRegisterCall — writes a single eng.Register call to buffer
package codegen

import (
	"bytes"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// writeRegisterCall appends one eng.Register(RuleMeta{...}) call to buf.
func writeRegisterCall(buf *bytes.Buffer, meta toulmin.RuleMeta) {
	buf.WriteString("\teng.Register(toulmin.RuleMeta{\n")
	fmt.Fprintf(buf, "\t\tName:      %q,\n", meta.Name)
	fmt.Fprintf(buf, "\t\tQualifier: %g,\n", meta.Qualifier)
	fmt.Fprintf(buf, "\t\tStrength:  %s,\n", strengthString(meta.Strength))
	if len(meta.Defeats) > 0 {
		fmt.Fprintf(buf, "\t\tDefeats:   []string{%s},\n", formatDefeats(meta.Defeats))
	}
	if meta.Backing != "" {
		fmt.Fprintf(buf, "\t\tBacking:   %q,\n", meta.Backing)
	}
	fmt.Fprintf(buf, "\t\tFn:        %s,\n", meta.Name)
	buf.WriteString("\t})\n")
}
