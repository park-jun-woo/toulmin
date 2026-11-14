//ff:func feature=cli type=command control=sequence
//ff:what splitSubjectEndpoint — splits a "<subject>.<endpoint>" argument into its two parts
package tanglcli

import (
	"fmt"
	"strings"
)

// splitSubjectEndpoint splits s ("transfer.`transfer`" or "transfer.transfer")
// on its first '.' into a subject and an endpoint name, trimming any
// backticks from the endpoint part.
func splitSubjectEndpoint(s string) (string, string, error) {
	idx := strings.Index(s, ".")
	if idx < 0 {
		return "", "", fmt.Errorf("tangl: expected <subject>.<endpoint>, got %q", s)
	}
	return s[:idx], trimBackticks(s[idx+1:]), nil
}
