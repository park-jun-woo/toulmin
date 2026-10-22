//ff:type feature=approve type=model
//ff:what Budget: 예산 정보
package approve

// Budget represents budget status.
type Budget struct {
	Remaining float64
	Frozen    bool
}
