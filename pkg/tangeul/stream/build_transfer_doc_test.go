//ff:func feature=tangl type=util control=sequence
//ff:what buildTransferDoc — the §10 account-transfer Doc fixture (programmatic, tiling spans)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildTransferDoc builds the account-transfer Doc of grammar §10 (Subject, See,
// two Cases with requires/rules/attacks/do/undo/run, and a Provides endpoint)
// programmatically. Blocks are laid out in document order with sequential
// non-overlapping spans that tile the Source exactly, so VerifyTiling holds. The
// `run can deposit` edge points at a case declared later, exercising the forward
// reference path.
func buildTransferDoc() *doc.Doc {
	a := &ast.Document{Subject: "transfer"}
	a.Sees = []ast.See{{Alias: "bank", Package: "bank/core"}, {Alias: "log", Package: "tangl/log"}}
	a.Cases = []ast.Case{
		{
			Name:     "can withdraw",
			Requires: []ast.Require{{Field: "from account", Type: "Currency"}, {Field: "amount"}},
			Nodes: []ast.Node{
				{Name: "balance sufficient", Role: ast.GeneralRule, Using: &ast.Ref{Alias: "bank", Name: "hasBalance"}},
				{Name: "account frozen", Role: ast.CounterRule, Using: &ast.Ref{Alias: "bank", Name: "isFrozen"}},
			},
			Attacks: []ast.Attack{{Target: "balance sufficient", Attacker: "account frozen"}},
			Execs: []ast.Exec{
				{Kind: ast.DoExec, Func: &ast.Ref{Alias: "bank", Name: "withdraw"}, Node: "balance sufficient", Once: true},
				{Kind: ast.UndoExec, Func: &ast.Ref{Alias: "bank", Name: "refund"}, Node: "balance sufficient"},
				{Kind: ast.RunExec, Case: "can deposit", Node: "balance sufficient"},
			},
		},
		{
			Name:     "can deposit",
			Requires: []ast.Require{{Field: "to account"}},
			Nodes:    []ast.Node{{Name: "recipient valid", Role: ast.GeneralRule, Using: &ast.Ref{Alias: "bank", Name: "isValidAccount"}}},
			Execs:    []ast.Exec{{Kind: ast.DoExec, Func: &ast.Ref{Alias: "bank", Name: "deposit"}, Node: "recipient valid", Once: true}},
		},
	}
	a.Provides = []ast.Endpoint{{
		Name:     "transfer",
		Requires: []ast.Require{{Field: "from account"}, {Field: "to account"}, {Field: "amount"}},
		Runs:     []string{"can withdraw"},
	}}

	var src []byte
	var blocks []doc.Block
	add := func(kind doc.BlockKind, ref doc.StmtRef) {
		off := uint32(len(src))
		src = append(src, '-', '-', '-', '-')
		blocks = append(blocks, doc.Block{Kind: kind, Span: doc.Span{Off: off, Len: 4}, Ref: ref})
	}
	add(doc.SectionBlock, doc.StmtRef{Sec: codebook.SectionSubject})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionSubject})
	add(doc.SectionBlock, doc.StmtRef{Sec: codebook.SectionSee})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionSee, Path: []int{0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionSee, Path: []int{1}})
	add(doc.SectionBlock, doc.StmtRef{Sec: codebook.SectionCases})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 0, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 0, 1}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 1, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 1, 1}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 2, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 3, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 3, 1}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 3, 2}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{1}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{1, 0, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{1, 1, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionCases, Path: []int{1, 3, 0}})
	add(doc.SectionBlock, doc.StmtRef{Sec: codebook.SectionProvides})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 0, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 0, 1}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 0, 2}})
	add(doc.StmtBlock, doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 1, 0}})

	return &doc.Doc{AST: a, Source: src, Blocks: blocks}
}
