//ff:func feature=engine type=model control=sequence
//ff:what Trace.Ctx — this Run's context
package toulmin

// Ctx returns this Run's context.
func (t Trace) Ctx() Context { return t.ctx }
