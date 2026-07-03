//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what trimTrailingParticleKo — strip trailing ignored particle tokens from a value
package ko

import "strings"

// trimTrailingParticleKo removes trailing whitespace-delimited particle tokens
// (koParticles members, e.g. 와/과) from a comparison value. skipParticlesKo only
// peels leading particles, but the KO 정규형 writes the object marker after the
// value ("VIP" 와 같으면), so the value's trailing 와/과 must be dropped here.
func trimTrailingParticleKo(s string) string {
	s = strings.TrimSpace(s)
	for {
		fields := strings.Fields(s)
		if len(fields) == 0 {
			return s
		}
		last := fields[len(fields)-1]
		if !koParticles[last] {
			return s
		}
		s = strings.TrimSpace(s[:len(s)-len(last)])
	}
}
