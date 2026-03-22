//ff:func feature=codegen type=codegen control=sequence
//ff:what float64Ptr — helper that returns a pointer to a float64 for testing
package codegen

func float64Ptr(v float64) *float64 { return &v }
