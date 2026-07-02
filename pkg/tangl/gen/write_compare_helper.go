//ff:func feature=tangl type=codegen control=sequence
//ff:what writeCompareHelper — writes the shared tanglCompare function
package gen

import "strings"

// writeCompareHelper writes the shared tanglCompare function that every
// generated Compare leaf calls: it fetches the field from ctx and applies
// one of the eight tangl comparison operators against value.
func writeCompareHelper(w *strings.Builder) {
	w.WriteString(`func tanglCompare(ctx toulmin.Context, field, op, value string) bool {
	v, ok := ctx.Get(field)
	switch op {
	case "is empty":
		return !ok || v == nil || fmt.Sprintf("%v", v) == ""
	case "is not empty":
		return ok && v != nil && fmt.Sprintf("%v", v) != ""
	case "equals":
		return ok && fmt.Sprintf("%v", v) == value
	case "is in":
		if !ok {
			return false
		}
		for _, item := range strings.Split(value, ",") {
			if strings.TrimSpace(item) == fmt.Sprintf("%v", v) {
				return true
			}
		}
		return false
	case "is greater than", "is less than", "is at most", "is at least":
		if !ok {
			return false
		}
		fv, ferr := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
		tv, terr := strconv.ParseFloat(value, 64)
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
