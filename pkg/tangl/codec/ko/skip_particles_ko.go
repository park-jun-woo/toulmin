//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what skipParticlesKo — consume leading standalone particle tokens and commas
package ko

import "strings"

// skipParticlesKo repeatedly consumes leading spaces/tabs and then either a
// leading comma or a leading whitespace-delimited token that is in koParticles,
// returning the remaining text. It stops at the first token that is neither a
// comma nor an ignored particle (e.g. a backtick name or a keyword), so it never
// reaches past the particle run. KO parsers call it after consuming each token.
func skipParticlesKo(s string) string {
	for {
		s = strings.TrimLeft(s, " \t")
		if strings.HasPrefix(s, ",") {
			s = s[1:]
			continue
		}
		tok := s
		if i := strings.IndexAny(s, " \t"); i >= 0 {
			tok = s[:i]
		}
		if !koParticles[tok] {
			return s
		}
		s = s[len(tok):]
	}
}
