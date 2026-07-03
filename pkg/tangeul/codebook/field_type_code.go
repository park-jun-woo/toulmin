//ff:type feature=tangl type=model
//ff:what FieldTypeCode — the ten TANGEUL basic field-type codes (ast.Field.Type)
package codebook

// FieldTypeCode is the code for a TANGEUL basic type name written in "as <Type>"
// clauses. It maps to the ten canonical ast.Field.Type strings (the same set as
// pkg/tangl/types.Name). Free-area L1 codebook value, FROZEN append-only.
type FieldTypeCode uint16

// The ten TANGEUL basic type codes, assigned 0x0201..0x020A.
const (
	// FieldText is "Text" (0x0201).
	FieldText FieldTypeCode = 0x0201
	// FieldInteger is "Integer" (0x0202).
	FieldInteger FieldTypeCode = 0x0202
	// FieldNumber is "Number" (0x0203).
	FieldNumber FieldTypeCode = 0x0203
	// FieldBoolean is "Boolean" (0x0204).
	FieldBoolean FieldTypeCode = 0x0204
	// FieldEmail is "Email" (0x0205).
	FieldEmail FieldTypeCode = 0x0205
	// FieldDate is "Date" (0x0206).
	FieldDate FieldTypeCode = 0x0206
	// FieldTime is "Time" (0x0207).
	FieldTime FieldTypeCode = 0x0207
	// FieldURL is "URL" (0x0208).
	FieldURL FieldTypeCode = 0x0208
	// FieldCurrency is "Currency" (0x0209).
	FieldCurrency FieldTypeCode = 0x0209
	// FieldQuantity is "Quantity" (0x020A).
	FieldQuantity FieldTypeCode = 0x020A
)
