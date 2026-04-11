//ff:func feature=engine type=engine control=sequence
//ff:what resolveOption — resolves EvalOption from variadic args with defaults
package toulmin

import "fmt"

// resolveOption returns the first option or a zero-value default.
// Duration implies Trace.
// Returns error if an unsupported method is requested.
func resolveOption(opts []EvalOption) (EvalOption, error) {
	if len(opts) == 0 {
		return EvalOption{}, nil
	}
	opt := opts[0]
	if opt.Method == Recursive {
		return EvalOption{}, fmt.Errorf("toulmin: EvalMethod Recursive is not yet implemented")
	}
	if opt.Duration {
		opt.Trace = true
	}
	return opt, nil
}
