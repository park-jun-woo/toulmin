//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestJosaFor — josaFor covers 받침/무받침 Hangul, the ㄹ 로/으로 특칙, digit 낭독, and ASCII
package ko

import "testing"

// TestJosaFor checks 받침-driven allomorph selection across all input classes: a
// 받침 Hangul syllable takes the batchim form, an open syllable the open form, a ㄹ
// jongseong takes 로 only in the 으로/로 pair (batchim elsewhere), each digit follows
// its 낭독 받침, and ASCII / empty are treated as 받침 없음.
func TestJosaFor(t *testing.T) {
	cases := []struct {
		name, batchim, open, want string
	}{
		{"금액", "은", "는", "은"},      // 액: ㄱ jongseong → 받침
		{"직원", "이면", "면", "이면"},    // 원: ㄴ jongseong → 받침
		{"센서", "은", "는", "는"},      // 서: no jongseong → 무받침
		{"신용점수", "이", "가", "가"},    // 수: 무받침
		{"연체 이력", "이", "가", "이"},   // 력: ㄱ jongseong → 받침
		{"팔", "은", "는", "은"},       // ㄹ jongseong, 은/는 pair → 받침
		{"팔", "으로", "로", "로"},      // ㄹ jongseong, 으로/로 pair → 로 특칙
		{"은행", "으로", "로", "으로"},    // 행: ㅇ jongseong 받침 (not ㄹ) → 으로
		{"650", "을", "를", "을"},     // 0(영) → 받침
		{"5", "을", "를", "를"},       // 5(오) → 무받침
		{"7", "으로", "로", "으로"},     // 7(칠) → 받침; ㄹ 특칙 is Hangul-only, digit table → 으로
		{"refund", "으로", "로", "로"}, // ASCII → 무받침
		{`"VIP"`, "과", "와", "와"},   // trailing quote → 기타 → 무받침
		{"", "은", "는", "는"},        // empty → 무받침
	}
	for _, c := range cases {
		if got := josaFor(c.name, c.batchim, c.open); got != c.want {
			t.Errorf("josaFor(%q, %q, %q) = %q, want %q", c.name, c.batchim, c.open, got, c.want)
		}
	}
}
