//ff:func feature=tangl type=parser control=selection
//ff:what koSectionNames — map a Korean `## tangl:<이름>` heading to its canonical English section name
package ko

// koSectionNames maps a Korean section heading name (the verbatim label after
// `## tangl:`, e.g. "주제") to the canonical English section name the shared parser
// and scanner recognize ("Subject"), so KO documents reuse the locale-neutral
// dispatch. ok is false for any other name (a scanner/parser drift the caller
// reports). Order matches the required document order (문법 v0.3 §일곱 섹션).
func koSectionNames(name string) (string, bool) {
	switch name {
	case "주제":
		return "Subject", true
	case "참조":
		return "See", true
	case "정의":
		return "Definitions", true
	case "규칙":
		return "Rules", true
	case "사안":
		return "Cases", true
	case "업무":
		return "Provides", true
	case "내부":
		return "Internal", true
	default:
		return "", false
	}
}
