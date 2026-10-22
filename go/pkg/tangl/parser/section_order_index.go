//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what sectionOrderIndex — look up a section name's required document order
package parser

// sectionNames lists the seven recognized tangl: sections in required order.
var sectionNames = []string{
	"Subject", "See", "Definitions", "Rules", "Cases", "Provides", "Internal",
}

// sectionOrderIndex returns the required document-order index of a section
// name, or ok=false if the name is not one of the seven recognized sections.
func sectionOrderIndex(name string) (int, bool) {
	for i, n := range sectionNames {
		if n == name {
			return i, true
		}
	}
	return -1, false
}
