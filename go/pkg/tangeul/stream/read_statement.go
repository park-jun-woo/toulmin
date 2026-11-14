//ff:func feature=tangl type=util control=sequence
//ff:what readStatement — collect and interpret one statement, appending its Block; return next index
package stream

// readStatement collects one statement's packets starting at packets[i],
// interprets them into an AST element and its Block (appending the Block), and
// returns the index just past the statement for the reader loop to continue from.
func readStatement(d *dec, packets []rawPacket, i int) (int, error) {
	prefix, body, next, err := collectStatement(d, packets, i)
	if err != nil {
		return 0, err
	}
	blk, err := interpretStmt(d, prefix, body)
	if err != nil {
		return 0, err
	}
	d.blocks = append(d.blocks, blk)
	return next, nil
}
