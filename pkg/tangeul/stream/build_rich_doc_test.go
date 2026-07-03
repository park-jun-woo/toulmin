//ff:func feature=tangl type=util control=sequence
//ff:what buildRichDoc — a Doc fixture exercising every section, prose, gates, with/qualified/checking
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildRichDoc covers the paths the transfer fixture omits: a prose block, a
// const Definition with an `as` spec reference plus a struct Definition, a Rules
// inline rule with a nested condition tree, an in-case node carrying `with`,
// another carrying `checking` (a forward case) and a qualifier, a do edge with a
// confidence gate, endpoint checks, and both Internal trigger kinds.
func buildRichDoc() *doc.Doc {
	qual := 0.8
	a := &ast.Document{Subject: "sys"}
	a.Sees = []ast.See{{Alias: "sensor", Package: "robot/sensor"}}
	a.Defs = []ast.Definition{
		{Name: "serve temperature", Kind: ast.ConstDef, Value: "65°C", SpecRef: &ast.Ref{Alias: "sensor", Name: "Temperature"}},
		{Name: "config", Kind: ast.StructDef, Fields: []ast.Field{{Name: "host", Type: "Text"}, {Name: "port", Type: "Integer"}}},
	}
	a.Rules = []ast.InlineRule{{Name: "ready", Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
		ast.Compare{Field: "level", Op: "is at least", Value: "5"},
		ast.Not{Term: ast.Compare{Field: "name", Op: "is empty"}},
	}}}}
	a.Cases = []ast.Case{
		{
			Name:     "can serve",
			Requires: []ast.Require{{Field: "user", Type: "Text"}},
			Nodes: []ast.Node{
				{Name: "brewing complete", Role: ast.GeneralRule, Using: &ast.Ref{Alias: "sensor", Name: "isComplete"}},
				{Name: "low temperature", Role: ast.ExceptRule, Using: &ast.Ref{Alias: "sensor", Name: "temperatureLow"}, With: []string{"serve temperature"}},
				{Name: "authenticate", Role: ast.GeneralRule},
				{Name: "composed", Role: ast.GeneralRule, Checking: "can check", Qualified: &qual},
			},
			Attacks: []ast.Attack{{Target: "brewing complete", Attacker: "low temperature"}},
			Execs: []ast.Exec{
				{Kind: ast.DoExec, Func: &ast.Ref{Alias: "sensor", Name: "serve"}, Node: "brewing complete", Once: true, Certainty: &ast.Certainty{Op: "at least", Percent: 75}},
				{Kind: ast.DoExec, Func: &ast.Ref{Name: "report"}, Node: "authenticate"},
			},
		},
		{
			Name:  "can check",
			Nodes: []ast.Node{{Name: "ok", Role: ast.GeneralRule, Using: &ast.Ref{Alias: "sensor", Name: "ping"}}},
			Execs: []ast.Exec{{Kind: ast.DoExec, Func: &ast.Ref{Alias: "sensor", Name: "ack"}, Node: "ok", Once: true}},
		},
	}
	a.Provides = []ast.Endpoint{{Name: "serve", Requires: []ast.Require{{Field: "user"}}, Runs: []string{"can serve"}, Checks: []string{"can check"}}}
	a.Internals = []ast.Internal{
		{Kind: ast.EveryTick, Interval: "1s", Until: "can check", Runs: []string{"can serve"}},
		{Kind: ast.OnEvent, Event: "sensor.buttonPress", Runs: []string{"can serve"}, Checks: []string{"can check"}},
	}

	var src []byte
	var blocks []doc.Block
	add := func(kind doc.BlockKind, ref doc.StmtRef) {
		off := uint32(len(src))
		src = append(src, '-', '-', '-', '-')
		blocks = append(blocks, doc.Block{Kind: kind, Span: doc.Span{Off: off, Len: 4}, Ref: ref})
	}
	sr := func(sec codebook.SectionCode, p ...int) doc.StmtRef { return doc.StmtRef{Sec: sec, Path: p} }

	add(doc.SectionBlock, sr(codebook.SectionSubject))
	add(doc.StmtBlock, sr(codebook.SectionSubject))
	add(doc.ProseBlock, doc.StmtRef{})
	add(doc.SectionBlock, sr(codebook.SectionSee))
	add(doc.StmtBlock, sr(codebook.SectionSee, 0))
	add(doc.SectionBlock, sr(codebook.SectionDefinitions))
	add(doc.StmtBlock, sr(codebook.SectionDefinitions, 0))
	add(doc.StmtBlock, sr(codebook.SectionDefinitions, 1))
	add(doc.StmtBlock, sr(codebook.SectionDefinitions, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionDefinitions, 1, 1))
	add(doc.SectionBlock, sr(codebook.SectionRules))
	add(doc.StmtBlock, sr(codebook.SectionRules, 0))
	add(doc.SectionBlock, sr(codebook.SectionCases))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 0, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 1, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 1, 2))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 1, 3))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 2, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 3, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 3, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1, 3, 0))
	add(doc.SectionBlock, sr(codebook.SectionProvides))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0, 0, 0))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0, 2, 0))
	add(doc.SectionBlock, sr(codebook.SectionInternal))
	add(doc.StmtBlock, sr(codebook.SectionInternal, 0))
	add(doc.StmtBlock, sr(codebook.SectionInternal, 0, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionInternal, 1))
	add(doc.StmtBlock, sr(codebook.SectionInternal, 1, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionInternal, 1, 2, 0))

	return &doc.Doc{AST: a, Source: src, Blocks: blocks}
}
