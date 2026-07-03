//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what skipParticlesKo — consume leading standalone particle tokens and commas
package ko

import "strings"

// skipParticlesKo repeatedly consumes leading spaces/tabs and then either a
// leading comma or a leading token — delimited by whitespace or a comma — that is
// in koParticles, returning the remaining text. Treating the comma as a token
// delimiter (not just a leading char) lets a particle abutting a separator comma
// (e.g. "의, …", left by strings.Join(mods, ", ")) split into the particle plus the
// comma, so the particle is peeled and the leading-comma branch then drops the ",".
// It stops at the first token that is neither a comma nor an ignored particle (e.g.
// a backtick name or a keyword), so it never reaches past the particle run. KO
// parsers call it after consuming each token.
func skipParticlesKo(s string) string {
	for {
		s = strings.TrimLeft(s, " \t")
		if strings.HasPrefix(s, ",") {
			s = s[1:]
			continue
		}
		tok := s
		if i := strings.IndexAny(s, " \t,"); i >= 0 {
			tok = s[:i]
		}
		if !koParticles[tok] {
			return s
		}
		s = s[len(tok):]
	}
}
