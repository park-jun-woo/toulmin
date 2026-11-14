//ff:type feature=engine type=model
//ff:what defeatEdge — a defeat relationship between two rules in the builder
package toulmin

type defeatEdge struct {
	from  string
	to    string
	force float64 // attack force F in [0, 1] (defeat edges only; unused by supersede edges)
}
