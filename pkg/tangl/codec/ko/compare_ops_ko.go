//ff:type feature=tangl type=model
//ff:what compareOpsKo — KO comparison-operator suffix table mapping to canonical English Op
package ko

// compareOpKo is one bidirectional row of the KO 규칙 operator table: it maps a KO
// trailing operator phrase — the sentence-final 서술형 (Ko) and its 관형형 attributive
// variant (KoAdn, used inside a one-line exists body) — to its canonical English
// Compare.Op and whether the comparison is negated (Negated, wrapped in ast.Not),
// and — for the reverse render direction — records whether the operator carries a
// value (HasValue) and whether the value takes a 와/과 object particle before the
// predicate (ObjParticle). Only equals uses the object particle (`f` 가 <v> 와 같으면);
// the value-less 비어 있음 forms carry neither.
type compareOpKo struct {
	Ko          string
	KoAdn       string
	Op          string
	Negated     bool
	HasValue    bool
	ObjParticle bool
}

// compareOpsKo is the KO 규칙 operator table (문법 v0.3 §연산자·§논리), ordered
// longest phrase first so a compound negative (보다 크지 않으면) is matched before its
// shorter positive (보다 크면). The same row order also satisfies the longest-first
// property for the 관형형 (KoAdn) column (보다 크지 않은 before 보다 큰,
// 중 하나가 아닌 before 중 하나인). 비어 있지 않으면 maps directly to the base operator
// "is not empty" — not a Not wrapper — because that operator is one of the eight
// canonical operators. Op strings match pkg/tangl/parser.compareOps exactly. The
// single table is shared by parseCompareKo/parseCompareAdnKo (suffix match) and
// compareKoPhrase/compareKoAdnPhrase (reverse render lookup) so the directions never
// drift.
var compareOpsKo = []compareOpKo{
	{"비어 있지 않으면", "비어 있지 않은", "is not empty", false, false, false},
	{"보다 크지 않으면", "보다 크지 않은", "is greater than", true, true, false},
	{"보다 작지 않으면", "보다 작지 않은", "is less than", true, true, false},
	{"중 하나가 아니면", "중 하나가 아닌", "is in", true, true, false},
	{"이상이 아니면", "이상이 아닌", "is at least", true, true, false},
	{"이하가 아니면", "이하가 아닌", "is at most", true, true, false},
	{"비어 있으면", "비어 있는", "is empty", false, false, false},
	{"같지 않으면", "같지 않은", "equals", true, true, true},
	{"보다 작으면", "보다 작은", "is less than", false, true, false},
	{"중 하나면", "중 하나인", "is in", false, true, false},
	{"보다 크면", "보다 큰", "is greater than", false, true, false},
	{"이상이면", "이상인", "is at least", false, true, false},
	{"이하면", "이하인", "is at most", false, true, false},
	{"같으면", "같은", "equals", false, true, true},
}
