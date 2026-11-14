//ff:type feature=tangl type=codegen
//ff:what equivAuthor — cross-test struct walked by the author.trustScore path
package gen

// equivAuthor is the cross-test struct the `author`.`trustScore` data path
// walks (case-insensitive exported-field match); the generated runner defines
// an identical copy.
type equivAuthor struct {
	TrustScore float64
}
