//ff:type feature=tangl type=model
//ff:what StmtRef — an AST element coordinate (section code + variable-length index path)
package doc

import "github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"

// StmtRef is the AST coordinate of a statement, linking a packet/Block to an AST
// element without mutating the AST. Sec is a codebook SectionCode (1 Subject ..
// 7 Internal) and Path is a variable-length index path whose shape is section
// specific; the normative encoding for all seven sections is the table in
// README contract §2, implemented by ResolveStmt.
//
//	Subject     nil               subject declaration (Document.Subject, plain string)
//	See         [i]               i-th `see` declaration
//	Definitions [defIdx]          definition declaration itself
//	Definitions [defIdx, fldIdx]  struct `has` field
//	Rules       [ruleIdx]         inline rule subtree in full (never deeper)
//	Cases       [caseIdx]         case declaration itself
//	Cases       [caseIdx, sk, si] case body: sk 0 Require 1 Node 2 Attack 3 Exec
//	Provides    [epIdx]           endpoint declaration itself
//	Provides    [epIdx, sk, si]   endpoint body: sk 0 Require 1 Run 2 Check
//	Internal    [i]               internal trigger declaration itself
//	Internal    [i, sk, si]       internal body: sk 1 Run 2 Check (0 unused)
type StmtRef struct {
	Sec  codebook.SectionCode
	Path []int
}
