//ff:func feature=tangl type=codegen control=sequence
//ff:what writeCompareValHelper — writes the shared tanglCompareVal function
package gen

import "strings"

// writeCompareValHelper writes the shared tanglCompareVal function: the
// value-vs-value generalization of the eight tangl comparison operators,
// mirroring pkg/tangl/eval.compareVals — both operands are compared in their
// fmt "%v" string form, and a no-value lhs (ok=false) satisfies only
// "is empty".
func writeCompareValHelper(w *strings.Builder) {
	w.WriteString(`func tanglCompareVal(lhs any, ok bool, op string, rhs any) bool {
	switch op {
	case "is empty":
		return !ok || lhs == nil || fmt.Sprintf("%v", lhs) == ""
	case "is not empty":
		return ok && lhs != nil && fmt.Sprintf("%v", lhs) != ""
	case "equals":
		return ok && fmt.Sprintf("%v", lhs) == fmt.Sprintf("%v", rhs)
	case "is in":
		if !ok {
			return false
		}
		for _, item := range strings.Split(fmt.Sprintf("%v", rhs), ",") {
			if strings.TrimSpace(item) == fmt.Sprintf("%v", lhs) {
				return true
			}
		}
		return false
	case "is greater than", "is less than", "is at most", "is at least":
		if !ok {
			return false
		}
		fv, ferr := strconv.ParseFloat(fmt.Sprintf("%v", lhs), 64)
		tv, terr := strconv.ParseFloat(fmt.Sprintf("%v", rhs), 64)
		if ferr != nil || terr != nil {
			return false
		}
		switch op {
		case "is greater than":
			return fv > tv
		case "is less than":
			return fv < tv
		case "is at most":
			return fv <= tv
		default:
			return fv >= tv
		}
	default:
		return false
	}
}

`)
}
