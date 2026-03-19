//ff:type feature=state type=model
//ff:what transition: 등록된 전이의 메타데이터
package state

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// transition holds a registered transition's metadata.
type transition struct {
	from  string
	event string
	to    string
	graph *toulmin.Graph
}
