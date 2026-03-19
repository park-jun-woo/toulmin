//ff:func feature=policy type=util control=iteration dimension=1
//ff:what formatTrace: Trace 결과를 사람이 읽을 수 있는 문자열로 변환
package policy

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func formatTrace(traces []toulmin.TraceEntry) string {
	parts := make([]string, len(traces))
	for i, t := range traces {
		parts[i] = fmt.Sprintf("%s=%v", t.Name, t.Activated)
	}
	return strings.Join(parts, ", ")
}
