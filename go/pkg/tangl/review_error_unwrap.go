//ff:func feature=tangl type=model control=sequence
//ff:what ReviewError.Unwrap — exposes both wrapped errors to errors.Is/errors.As
package tangl

// Unwrap returns both the original cause and the compensation error, so
// errors.Is and errors.As can traverse into either one.
func (r *ReviewError) Unwrap() []error {
	return []error{r.Cause, r.Compensate}
}
