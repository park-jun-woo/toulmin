//ff:func feature=tangl type=parser control=selection
//ff:what applySection — parse one section and append its result onto a Document
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// applySection dispatches sec by name to its section parser and appends the
// parsed result onto doc.
func applySection(doc *ast.Document, sec section, path string) error {
	switch sec.Name {
	case "Subject":
		name, err := parseSubjectSection(sec, path)
		if err != nil {
			return err
		}
		doc.Subject = name
	case "See":
		sees, err := parseSeeSection(sec, path)
		if err != nil {
			return err
		}
		doc.Sees = append(doc.Sees, sees...)
	case "Definitions":
		defs, err := parseDefinitionsSection(sec, path)
		if err != nil {
			return err
		}
		doc.Defs = append(doc.Defs, defs...)
	case "Rules":
		rules, err := parseRulesSection(sec, path)
		if err != nil {
			return err
		}
		doc.Rules = append(doc.Rules, rules...)
	case "Cases":
		cases, err := parseCasesSection(sec, path)
		if err != nil {
			return err
		}
		doc.Cases = append(doc.Cases, cases...)
	case "Provides":
		eps, err := parseProvidesSection(sec, path)
		if err != nil {
			return err
		}
		doc.Provides = append(doc.Provides, eps...)
	case "Internal":
		ins, err := parseInternalSection(sec, path)
		if err != nil {
			return err
		}
		doc.Internals = append(doc.Internals, ins...)
	}
	return nil
}
