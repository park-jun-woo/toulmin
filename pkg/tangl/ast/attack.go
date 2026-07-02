//ff:type feature=tangl type=model
//ff:what Attack — a defeat edge (`don't <target> when <attacker>`)
package ast

// Attack is a defeat edge: attacker.Attacks(target).
// `do not` is equivalent to `don't`.
type Attack struct {
	Target   string `json:"target"`
	Attacker string `json:"attacker"`
	Line     int    `json:"line"`
}
