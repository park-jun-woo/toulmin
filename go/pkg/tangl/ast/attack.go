//ff:type feature=tangl type=model
//ff:what Attack — a defeat edge (`don't <target> when <attacker>` / `<attacker> attacks <target> with force F`)
package ast

// Attack is a defeat edge: attacker.Attacks(target[, force]).
// `do not` is equivalent to `don't`. Force is nil for the omitted form
// (`don't X when Y`, force 1, GEUL F_FORCE off) and non-nil for the explicit
// form (`Y attacks X with force F`, F in [0, 1], F_FORCE on — even for 1.0).
type Attack struct {
	Target   string   `json:"target"`
	Attacker string   `json:"attacker"`
	Force    *float64 `json:"force,omitempty"`
	Line     int      `json:"line"`
}
