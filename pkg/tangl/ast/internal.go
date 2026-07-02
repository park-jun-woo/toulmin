//ff:type feature=tangl type=model
//ff:what Internal — a tangl:Internal entry (event/tick driven work)
package ast

// Internal is a single tangl:Internal entry: an event- or tick-driven
// trigger that runs or checks Cases without an external request.
type Internal struct {
	Kind     InternalKind `json:"kind"`
	Event    string       `json:"event,omitempty"`    // `on` clause source (e.g. "sensor.buttonPress", "start")
	Interval string       `json:"interval,omitempty"` // "30s", "1hour", "7:00"
	Until    string       `json:"until,omitempty"`    // termination condition case name
	Runs     []string     `json:"runs,omitempty"`
	Checks   []string     `json:"checks,omitempty"`
	Line     int          `json:"line"`
}
