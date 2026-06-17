//ff:type feature=engine type=model
//ff:what NodeHandler — node event handler signature
package toulmin

// NodeHandler is invoked by Run when a node fires its classified event.
type NodeHandler func(ctx Context, ev NodeEvent, view RunView) error
