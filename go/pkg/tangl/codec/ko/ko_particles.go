//ff:type feature=tangl type=model
//ff:what koParticles — the set of ignored Korean particle tokens (문법 v0.3 §조사 정책)
package ko

// koParticles is the ignore set of standalone Korean particle tokens the lexer
// drops when they appear as whitespace-delimited tokens (문법 v0.3 §조사 정책).
// Both members of every 받침-driven allomorph pair are present so writing either
// form never changes parsing: 은/는, 이/가, 을/를, 와/과, 로/으로, 면/이면.
//
// Tokens that only look like particles are deliberately excluded (§주의):
// 까지(until), 보다/중(복합 연산자), 이다/필수다/무효다/규칙이다(계사 결합 키워드),
// 마다(단위 복합). The comma is handled directly by skipParticlesKo, not here.
var koParticles = map[string]bool{
	"은":  true,
	"는":  true,
	"이":  true,
	"가":  true,
	"을":  true,
	"를":  true,
	"의":  true,
	"로":  true,
	"으로": true,
	"에":  true,
	"에서": true,
	"와":  true,
	"과":  true,
	"도":  true,
	"면":  true,
	"이면": true,
}
