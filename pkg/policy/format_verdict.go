//ff:func feature=policy type=util control=sequence
//ff:what formatVerdict: verdict를 문자열로 변환
package policy

import "fmt"

func formatVerdict(v float64) string {
	return fmt.Sprintf("%.2f", v)
}
