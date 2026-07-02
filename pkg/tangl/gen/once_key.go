//ff:func feature=tangl type=codegen control=sequence
//ff:what onceKey — builds the codegen once-guard key convention
package gen

import "fmt"

// onceKey builds the codegen once-guard key convention:
// "once:<subject>.<case>.<node>#<n>" for the n-th (0-based) "do" on node.
func onceKey(subject, caseName, nodeName string, n int) string {
	return fmt.Sprintf("once:%s.%s.%s#%d", subject, caseName, nodeName, n)
}
