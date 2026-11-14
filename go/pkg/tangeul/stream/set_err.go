//ff:func feature=tangl type=util control=sequence
//ff:what setErr — record the first record-encoding error on the Writer state
package stream

// setErr keeps err on e unless an earlier error is already recorded, so the
// error emitBody reports is the first record that failed to encode. A nil err
// changes nothing, so the record helpers pass every encoder result through.
func setErr(e *enc, err error) {
	if e.err == nil {
		e.err = err
	}
}
