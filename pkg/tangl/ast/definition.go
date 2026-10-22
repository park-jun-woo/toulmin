//ff:type feature=tangl type=model
//ff:what Definition — a tangl:Definitions entry (constant or struct term)
package ast

// Definition is a single entry of the tangl:Definitions section.
type Definition struct {
	Name     string   `json:"name"`
	Kind     DefKind  `json:"kind"`
	Value    string   `json:"value,omitempty"`    // ConstDef literal source ("650", "65°C")
	SpecRef  *Ref     `json:"specRef,omitempty"`  // `as` clause (`credit`.`Threshold`)
	Fields   []Field  `json:"fields,omitempty"`   // StructDef `has` clauses
	Strings  []string `json:"strings,omitempty"`  // StringListDef decoded items
	LoadType string   `json:"loadType,omitempty"` // FileSpecDef registry name
	FilePath string   `json:"filePath,omitempty"` // FileSpecDef source path
	Line     int      `json:"line"`
}
