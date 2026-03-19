# Claim, Not Fact: Applying Toulmin's Argumentation Model to Software Rule Engines

## Abstract

Software rule engines — Rego/OPA, Drools, Semgrep, JSON Schema — share a common design assumption: the data being validated is a "fact." This paper argues that validation targets are not facts but **claims**, and that Toulmin's argumentation model (1958) provides the correct design foundation for software rule engines that has been overlooked for over 60 years.

We present `toulmin`, a Go rule engine that implements Toulmin's six elements — Claim, Ground, Warrant, Backing, Qualifier, Rebuttal — with Nute's strict/defeasible/defeater classification for rule strength and Amgoud's h-Categoriser for verdict computation on a [-1, 1] scale. Rules are Go functions (`func(claim, ground, backing) → (bool, evidence)`) organized into defeats graphs; backing is a first-class runtime value carrying the judgment criteria for each rule. The engine computes verdicts and provides full trace explainability without requiring a custom DSL.

We validate the design by converting filefunc's 22 code structure rules into Toulmin warrants and measuring quantitative effects across three projects. We further demonstrate phase-optional architecture through fullend's compile-time/runtime policy integration.

---

## 1. Introduction

### 1.1 Problem: The Fact Assumption

Modern software rule engines universally treat their inputs as "facts." Drools inserts Java objects into working memory as "facts." Rego evaluates policies against an `input` document implicitly assumed to be true. JSON Schema validates against a document whose structure is taken as given.

This terminology is epistemologically inaccurate. A **fact** is something already established as true — it needs no verification. Yet the entire purpose of a rule engine is to **verify** whether input data satisfies rules. Calling a verification target "already true" is a contradiction.

Verification targets are **claims** — assertions that may be true or false, whose validity is to be determined by rules. This is not merely a terminological correction. The "fact" assumption has contributed to three structural limitations in existing rule engines (§3.1).

### 1.2 The Missing Link

Toulmin analyzed the structure of argumentation into six elements in 1958 [1]: Claim, Ground, Warrant, Backing, Qualifier, and Rebuttal. In subsequent decades, rule engines emerged — CLIPS (1985), Jess (1995), Drools (2001), Rego (2016) — all sharing the fact-based design. Toulmin's model already provided the correct structure, yet it was never applied as the design foundation for software rule engines.

The gap exists because Toulmin's work was published in philosophy and rhetoric, not in software engineering. Rule engine developers had no reason to consult argumentation theory. Meanwhile, each rule engine independently reinvented fragments of Toulmin's model: LegalRuleML added defeasibility mechanisms (Rebuttal), OPA added `# METADATA` annotations (partial Backing), Drools added salience (partial Qualifier). None referenced Toulmin.

### 1.3 Contributions

1. We demonstrate that Toulmin's argumentation model maps naturally onto software rule engine design, and that the "fact" terminology has contributed to structural limitations in existing systems (§3).
2. We implement `toulmin`, a Go rule engine with Toulmin's six elements, Nute's rule strength classification [2], and Amgoud's h-Categoriser [3] for verdict computation (§4).
3. We validate the design through filefunc's 22 rules across three projects with quantitative measurements (§5), and demonstrate phase-optional architecture through fullend's cross-phase policy integration (§6).

---

## 2. Background and Related Work

### 2.1 Toulmin's Argumentation Model

Toulmin [1] analyzed arguments into six elements:

- **Claim**: The assertion being made.
- **Ground (Data)**: The evidence supporting the claim.
- **Warrant**: The rule or principle that connects ground to claim.
- **Backing**: The justification for why the warrant is valid.
- **Qualifier**: The degree of certainty ("certainly," "presumably," "possibly").
- **Rebuttal**: Conditions under which the claim does not hold.

Gabriel et al. [4] implemented five of these elements (excluding Backing) in a BDI agent system using AgentSpeak/Jason, with a qualify function computing confidence as SWW - SWR (sum of warrant weights minus sum of rebuttal weights). This is the only prior software implementation of Toulmin's model we are aware of.

### 2.2 Defeasible Reasoning

Nute [2] classified rules into three strengths: **strict** (cannot be defeated), **defeasible** (can be defeated by rebuttals), and **defeaters** (block conclusions without asserting their own). This classification was formalized within ASPIC+ [5] and Dung's abstract argumentation framework [6].

### 2.3 Gradual Semantics

Amgoud and Ben-Naim [3] defined the weighted h-Categoriser for argumentation frameworks:

```
Acc(a) = w(a) / (1 + Σ Acc(attackers))
```

where `w(a)` is the initial weight and `attackers` are nodes attacking `a`. The h-Categoriser is the only gradual semantics satisfying the **Compensation** principle: if an attacker is itself attacked (defended), the original node's acceptability increases. Convergence is guaranteed [3].

### 2.4 Existing Rule Systems

| System | Input model | Defeasibility | Metadata | Phase |
|--------|------------|---------------|----------|-------|
| Rego/OPA | fact (`input`) | `default`/`else` patterns | `# METADATA` (not semantic) | Conftest: multi-phase |
| Drools | fact (working memory) | salience/agenda | None | Runtime only |
| Semgrep | AST nodes | None | None | Compile-time only |
| JSON Schema | JSON document | None | `description` | Compile-time only |
| LegalRuleML | fact | XML elements (bolt-on) | None | N/A |

No existing system treats inputs as claims, provides first-class defeasibility, includes semantic backing/qualifier, and supports phase-optional evaluation simultaneously.

### 2.5 Phase Integration

StaRVOOrS [7] and ppDATE [8] demonstrated that a single specification can feed both static verification (KeY) and runtime monitoring (LARVA). Cockburn's Hexagonal Architecture [9] formalized the adapter pattern separating domain logic from external systems.

---

## 3. From Facts to Claims

### 3.1 How "Fact" Naming Contributes to Design Limitations

The causes of structural limitations in rule engines are multifaceted — technology stack constraints, performance optimization, and use case scope all contribute. However, the "fact" naming acts as a cognitive frame that promotes certain design directions. Terminology shapes design philosophy, and design philosophy constrains implementation.

**Contribution 1: Fact → reduced verification awareness → consistency gaps.**
When Drools names working memory entries "facts," developers are cognitively primed to treat them as already verified. Cross-validation between facts requires explicit rules that developers must remember to write. Java's object model also contributes, but the "fact" naming provides a cognitive frame that makes verification omission feel natural.

**Contribution 2: Fact → data structure binding → rule duplication.**
When Rego treats `input` as a fixed JSON structure, rules reference concrete paths (`input.user.role`). Rego's JSON query syntax is also a contributing factor, but the fact-based design philosophy promotes binding rules to concrete data shapes rather than abstract properties. The same logical rule — "only admins can delete" — must be rewritten for HTTP requests, gRPC metadata, and CLI arguments.

**Contribution 3: Fact → fixed truth values → defeasibility as afterthought.**
Facts are true; true things are not overturned. Systems designed on this premise require external exception mechanisms for defeasibility. LegalRuleML needed `<Override>`, `<Defeater>`, `<Superiority>` XML elements — an entire subsystem added outside the core model. Redefining inputs as claims promotes designs where defeasibility is inherent, though the structural device of Rebuttal (§4) is still required.

**Precedent: JWT claims.**
JWT (RFC 7519) names token fields "claims," not "facts." `sub`, `exp`, `iss` are assertions by the token issuer. Signature verification, expiration checks, and issuer validation are expected precisely because the data is framed as claims. Terminology shaped the direction of design.

### 3.2 Toulmin Mapping

| Toulmin Element | Role in Rule Engine | Cardinality |
|---|---|---|
| Claim | Proposition to be judged. Input to rules | 1 |
| Ground | Evidence data. Supplied by Ground Adapter | 0..N |
| Warrant | Rule (bool function). Node in defeats graph | 1..N |
| Backing | Justification for warrant. Runtime value + annotation | 1..N |
| Qualifier | Initial weight w ∈ [0.0, 1.0] per rule **(see §3.3)** | 1 per rule |
| Rebuttal | Rule (bool function). Attacker in defeats graph | 0..N |

### 3.3 Repositioning the Qualifier: From Claim to Rule

In Toulmin's original model, the Qualifier is attached to the **Claim**. "**Presumably**, this patient should receive penicillin" — the Qualifier is a modal qualifier expressing the degree of certainty the speaker assigns to the claim, after considering Warrant, Ground, and Rebuttal together.

This model repositions the Qualifier from the Claim to **each Rule (warrant, rebuttal, defeater)**. This is an engineering correction necessitated by the transition from argumentation theory to rule engines.

In a rule engine, the claim is simply the validation target. "This file has 3 functions" — this is a factual observation; attaching a degree of certainty to it is meaningless. What determines the quality of the judgment is the **certainty of the rule**:

- "One func per file" — qualifier 1.0 (certain rule)
- "Recommended max 100 lines" — qualifier 0.7 (flexible rule)
- "Allergy possibility" — qualifier 0.95 (strong rebuttal)

This repositioning aligns naturally with Amgoud's h-Categoriser [3]. The h-Categoriser's `w(a)` is the initial weight of each node (rule), and the final output — the verdict — assumes the role that Qualifier played in Toulmin's original model: the degree of certainty of the judgment.

| | Toulmin's Original | This Model |
|--|-------------------|------------|
| Qualifier position | Claim (output side) | Each Rule (input side) |
| Role | Certainty of the final judgment | Initial weight of each rule |
| Certainty computation | Speaker judges holistically | h-Categoriser computes from rule weights |
| Final certainty | Qualifier itself | **Verdict** (h-Categoriser output) |

### 3.4 Epistemic Status of Warrants

Strictly speaking, warrants are also claims — "only admins can delete" is an organizational agreement, not a physical law. Even physical laws are high-reproducibility claims. However, applying this recognition infinitely leads to infinite regress of axioms. As Gödel's incompleteness theorem [10] shows, a sufficiently powerful formal system cannot prove its own consistency. Kelsen's Grundnorm [11] identifies the same structure in legal systems.

Therefore, this model takes a pragmatic position: **warrants are treated as facts within the system.** This is the unavoidable premise shared by all axiomatic systems, bounding the model's scope to "judging claims given warrants."

---

## 4. Engine Design

### 4.1 Rules as Functions with Evidence and Backing

All rules — warrants, rebuttals, defeaters — share the same signature:

```go
type Rule func(claim any, ground any, backing any) (bool, any)
```

A rule returns `(judgment, evidence)`. The bool is the judgment; the second value is domain-specific evidence (e.g., error details, violation context). This enables the engine to not only compute verdicts but also **explain why**.

The third parameter, `backing`, carries the judgment criteria for the rule. Ground and backing serve distinct epistemic roles:

- **Ground** = facts about the judgment target. Ground is per-request data that varies with each evaluation (e.g., a file's AST, a user's HTTP request).
- **Backing** = judgment criteria for the rule. Backing is fixed at declaration time and represents the standard or authority that justifies the warrant (e.g., a threshold value, a regulatory reference, a configuration policy).

This distinction maps directly to Toulmin's original model: ground supports the claim, while backing supports the warrant. By making backing a first-class runtime value, the engine enables the same rule function to be registered multiple times with different backing values, producing different judgment criteria without duplicating logic.

The distinction between warrant, rebuttal, and defeater is not in the function signature but in the **defeats graph**:

| Role | Graph position |
|------|---------------|
| Warrant | Node that can be attacked |
| Rebuttal | Node that attacks a warrant (has own conclusion) |
| Defeater | Node that attacks without asserting its own conclusion |

### 4.2 Graph: Separating Logic from Structure

The engine provides two APIs. The Graph API (recommended) separates rule logic (Go functions) from graph structure (defeats, qualifier, strength). The `Warrant`, `Rebuttal`, and `Defeater` methods accept `(fn, backing, qualifier)` and return a `*Rule` reference. The `Defeat` method accepts two `*Rule` references — backing is the second argument to registration methods, and `nil` when no backing is needed:

```go
g := toulmin.NewGraph("file-structure")
w := g.Warrant(CheckOneFileOneFunc, nil, 1.0)
d := g.Defeater(TestFileException, nil, 1.0)
g.Defeat(d, w)
```

Functions are identifiers — no string names needed. The same function can be reused across different graphs with different roles and defeats relationships:

```go
// Same IsAdult function, different graphs, different defeats
votingGraph := toulmin.NewGraph("voting")
w := votingGraph.Warrant(IsAdult, nil, 1.0)
r := votingGraph.Rebuttal(HasCriminalRecord, nil, 1.0)
votingGraph.Defeat(r, w)

contractGraph := toulmin.NewGraph("contract")
w2 := contractGraph.Warrant(IsAdult, nil, 1.0)
r2 := contractGraph.Rebuttal(IsBankrupt, nil, 1.0)
contractGraph.Defeat(r2, w2)
```

Since each registration call returns a distinct `*Rule` reference, the same function can be registered multiple times with different backing values. For example, two instances of a threshold-checking function — one backed by "corporate policy: max 100 lines" and another by "team convention: max 50 lines" — can coexist in the same graph with distinct identities:

```go
g := toulmin.NewGraph("line-limit")
corp := g.Warrant(CheckMaxLines, &LinePolicy{Max: 100, Source: "corporate"}, 1.0)
team := g.Defeater(CheckMaxLines, &LinePolicy{Max: 50, Source: "team"}, 0.8)
g.Defeat(team, corp)
```

This separation means `//tm:backing` is the only annotation on the function itself — it documents why the rule exists. All structural metadata (role, qualifier, strength, defeats) belongs to the graph, not the function. Backing as a runtime value complements the annotation: the annotation explains the rationale in source code, while the runtime value carries the judgment criteria into the engine.

### 4.3 Strength as Graph Constraint

Nute's classification [2] controls attack edges:

| Strength | Meaning |
|----------|---------|
| Strict | No incoming attack edges allowed |
| Defeasible | Incoming attack edges allowed |
| Defeater | Only outgoing attack edges; no verdict of its own |

### 4.4 Verdict Computation: h-Categoriser

We adopt Amgoud's weighted h-Categoriser [3] with a linear transform to [-1, 1]:

```
raw(a) = w(a) / (1 + Σ raw(attackers))     [0, 1]
verdict(a) = 2 × raw(a) - 1                [-1, 1]
```

| Verdict | Meaning |
|---------|---------|
| +1.0 | Fully confirmed violation |
| 0.0 | Undecided |
| -1.0 | Fully rebutted |

Judgment: `verdict > 0` → violation, `verdict == 0` → undecided, `verdict < 0` → rebutted.

### 4.5 Evaluate and EvaluateTrace

The engine provides two evaluation modes:

```go
// Evaluate — verdict + evidence (lightweight)
results, err := g.Evaluate(claim, ground)

// EvaluateTrace — verdict + evidence + full trace (explainability)
results, err := g.EvaluateTrace(claim, ground)
```

EvalResult:

```go
type EvalResult struct {
    Name     string       `json:"name"`
    Verdict  float64      `json:"verdict"`
    Evidence any          `json:"evidence,omitempty"`
    Trace    []TraceEntry `json:"trace"`
}

type TraceEntry struct {
    Name      string  `json:"name"`
    Role      string  `json:"role"`       // "warrant", "rebuttal", "defeater"
    Activated bool    `json:"activated"`
    Qualifier float64 `json:"qualifier"`
    Backing   any     `json:"backing,omitempty"`
    Evidence  any     `json:"evidence,omitempty"`
}
```

EvaluateTrace provides full explainability: which rules activated, in what role, with what qualifier, under what backing, producing what evidence. The `Backing` field in TraceEntry records the judgment criteria that were in effect for each rule execution, enabling auditors to understand not only what the rule decided but on what authority. This is a key advantage of symbolic reasoning — the verdict's derivation is transparent and auditable.

### 4.6 Evaluation Flow

```
0. Cycle detection: DFS on defeat edges at graph construction time
   → error returned if cycle found (before any func execution)
1. Start from each warrant node
2. Run warrant func(claim, ground, backing) → false? skip
   backing is the value registered with the node at graph construction time
3. If true, traverse attackers (rebuttal/defeater) recursively
4. Each attacker: run func(claim, ground, attacker's backing) → false? contributes 0 → true? recurse deeper
5. h-Categoriser at each node: raw(a) = w(a) / (1 + Σ raw(attackers))
   verdict(a) = 2 * raw(a) - 1
6. Func results cached — each (func, backing) pair runs at most once per evaluation
7. Only rules reachable from the warrant's attack chain are executed
```

### 4.7 YAML Graph Definition and Code Generation

Graphs can be defined in YAML and compiled to Go code:

```yaml
graph: file-structure
rules:
  - name: CheckOneFileOneFunc
    role: warrant
    qualifier: 1.0
  - name: TestFileException
    role: defeater
    qualifier: 1.0
defeats:
  - from: TestFileException
    to: CheckOneFileOneFunc
```

```bash
toulmin graph file-structure.yaml   # generates graph_gen.go
```

### 4.8 Metadata via Annotations

The `//tm:backing` annotation documents **why** the rule exists — the rationale in human-readable form. Backing is also a first-class runtime value passed to the rule function via the Graph API (§4.2). The two mechanisms are complementary:

- **`//tm:backing` annotation**: Static documentation visible in source code. Explains the authority or rationale behind the warrant (e.g., a theorem, a standard, a convention).
- **Backing parameter**: Runtime value injected by the engine. Carries judgment criteria that the rule function can use in its logic (e.g., threshold values, configuration policies).

When a rule's backing is purely documentary (no runtime data needed), the function receives `nil` as the backing parameter and the annotation alone suffices. When backing carries operational data, the annotation still documents the rationale while the runtime value provides the criteria:

```go
//tm:backing "Böhm-Jacopini theorem — all control flow reducible to sequence, selection, iteration"
func CheckOneFileOneFunc(claim any, ground any, backing any) (bool, any) {
    gf := ground.(*FileGround)
    if len(gf.Funcs) > 1 {
        return true, &Evidence{Got: len(gf.Funcs), Expected: 1}
    }
    return false, nil
}

//tm:backing "test files conventionally group multiple test funcs"
func TestFileException(claim any, ground any, backing any) (bool, any) {
    return strings.HasSuffix(claim.(string), "_test.go"), nil
}
```

When backing carries runtime data:

```go
//tm:backing "corporate coding standard — maximum lines per file"
func CheckMaxLines(claim any, ground any, backing any) (bool, any) {
    policy := backing.(*LinePolicy)
    gf := ground.(*FileGround)
    if gf.Lines > policy.Max {
        return true, &Evidence{Got: gf.Lines, Max: policy.Max, Source: policy.Source}
    }
    return false, nil
}
```

---

## 5. Case Study 1: filefunc — Toulmin Conversion of 22 Rules

### 5.1 Background

filefunc is a code structure convention and CLI tool for LLM-native Go development. It enforces 22 rules (F1-F6, Q1-Q3, A1-A16, C1-C4) via Go functions operating on AST at compile time. All rules were converted to Toulmin warrants.

### 5.2 Strength Classification

| Strength | Count | Ratio | Examples |
|----------|-------|-------|----------|
| Strict | 15 | 68% | F1, F2, F3, F4, A1-A3, A6-A16 |
| Defeasible | 4 | 18% | Q1, Q2, Q3, C4 |
| Defeater | 3 | 14% | F5, F6, test file exceptions |

Most rules are strict — code structure conventions inherently minimize exceptions. Defeasible rules are those with context-dependent thresholds (depth limits Q1-Q2, size limits Q3).

### 5.3 Quantitative Results

| Project | Files (before→after) | Avg LOC/file (before→after) | SRP violations resolved | Depth violations resolved |
|---------|---------------------|---------------------------|------------------------|--------------------------|
| filefunc | — (born compliant) | 25.1 | 0 | 0 |
| fullend | 87→1,260 | 244→25.4 | 66→0 | 148→0 |
| whyso | 12→99 | 147.8→24.4 | 12→0 | 23→0 |

### 5.4 Benefits of Toulmin Conversion

1. **Backing as first-class value**: Why each rule exists ("Böhm-Jacopini theorem," "AI agent's read unit is a file") is both documented via `//tm:backing` annotations and passed as a runtime value to the rule function. When backing carries operational data (thresholds, policies), the same function can be registered multiple times with different backing values, each returning a distinct `*Rule` reference.
2. **Rebuttal/Defeater structured**: Exception conditions (F5 test files, F6 grouped consts) are explicit `defeats` relations, not if-branches buried in code.
3. **Strength classification**: "This rule has no exceptions" (strict) vs "this rule may be defeated" (defeasible) is declaratively expressed.

---

## 6. Case Study 2: fullend — Phase-Optional Architecture

### 6.1 Background

fullend validates consistency across 9 SSOTs (OpenAPI, SSaC, Rego, etc.) and generates code. 18 cross-check rules are implemented as Go functions; runtime policies are separate Rego files. Where filefunc (§5) operates at a single phase (compile-time), fullend demonstrates the need for phase-optional architecture.

### 6.2 Problem: Duplicate Rules

The same logical rule is written twice:

**Compile-time** (Go crosscheck):
```go
func CheckRegoClaimsPresence(policies, openapi) []Error { ... }
```

**Runtime** (Rego policy):
```rego
allow if {
  input.action == "ExecuteWorkflow"
  input.claims.user_id > 0
}
```

### 6.3 Toulmin Solution

```go
//tm:backing "OWASP API Security Top 10 A2:2023"
func AuthEndpointRequiresClaims(claim any, ground any, backing any) (bool, any) {
    g := ground.(*EndpointGround)
    if g.Security.Contains("bearerAuth") && !g.PolicyRule.References("claims") {
        return true, &Evidence{Endpoint: g.Path, Missing: "claims reference"}
    }
    return false, nil
}

//tm:backing "Public API endpoints do not require authentication by design"
func PublicEndpointException(claim any, ground any, backing any) (bool, any) {
    g := ground.(*EndpointGround)
    return g.Annotation.Contains("x-public"), nil
}

// Graph definition
g := toulmin.NewGraph("endpoint-auth")
w := g.Warrant(AuthEndpointRequiresClaims, nil, 1.0)
d := g.Defeater(PublicEndpointException, nil, 1.0)
g.Defeat(d, w)
```

**Compile-time Ground Adapter**: `EndpointGround.Security` ← OpenAPI spec parse, `EndpointGround.PolicyRule` ← Rego AST parse.

**Runtime Ground Adapter**: `EndpointGround.Security` ← request route middleware, `EndpointGround.PolicyRule` ← OPA evaluation.

One warrant function serves both compile-time cross-validation and runtime policy enforcement. The graph structure is identical; only the Ground Adapter changes. The function itself is reusable across graphs — the same `AuthEndpointRequiresClaims` can participate in different defeats relationships in different contexts.

---

## 7. Discussion

### 7.1 What This Paper Does Not Claim

This paper does not claim theoretical originality. The structural foundation is Toulmin (1958), the strength classification is Nute (1994), the verdict computation is Amgoud (2017), the adapter pattern is Cockburn (2005) and StaRVOOrS (2015). Our contribution is **identifying that these independently developed pieces connect** — that Toulmin's argumentation model is the missing design foundation for software rule engines — and **demonstrating this through implementation**.

### 7.2 Limitations

- Rule bodies are Go functions (opaque to the engine). Partial evaluation, serialization, and cross-language portability are not supported in this implementation.
- The h-Categoriser's [-1, 1] transform is a linear mapping; more sophisticated verdict aggregation may be needed for complex rule interactions.
- Claim modeling requires further work — relational assertions involving two subjects exceed the expressiveness of simple (subject, property, value) triples.
- Quantitative validation is limited to the author's own projects. External validation on popular open-source libraries is planned.

### 7.3 Relationship to Existing Systems

This engine does not replace Rego, Drools, or Semgrep. It provides a **superordinate abstraction**. Rego can serve as a runtime Ground Adapter backend. Semgrep rules can be reused in compile-time Ground Adapters. Existing systems become Ground Adapter implementations for specific phases.

---

## 8. Conclusion

For 60 years, software rule engines have treated validation targets as "facts" and built designs on this assumption. Toulmin's argumentation model (1958) already provided the correct structure — Claim (not Fact), Ground, Warrant, Backing, Qualifier, Rebuttal — but the disciplinary gap between philosophy and software engineering prevented its application.

We implemented `toulmin`, a Go rule engine that bridges this gap. Rules are Go functions accepting `(claim, ground, backing)` and returning `(bool, evidence)` — judgment and its basis. Backing is a first-class runtime value carrying the judgment criteria for each rule, complementing the `//tm:backing` annotation that documents the rationale. The engine organizes rules into defeats graphs via the Graph API, evaluates them, and computes verdicts via Amgoud's h-Categoriser on a [-1, 1] scale. Nute's strict/defeasible/defeater classification controls attack edges in the graph. EvaluateTrace provides full explainability — which rules activated, in what role, under what backing, with what evidence.

The design requires no custom DSL — Go's type system handles claim-rule matching, and the engine core is under a few hundred lines. The same rule function can be reused across different graphs with different defeats relationships, separating rule logic from graph structure. Validation on filefunc's 22 rules across three projects demonstrates that the Toulmin model naturally accommodates rule strength classification, structured defeasibility, and explicit backing. fullend's cross-phase case demonstrates that Ground Adapters enable phase-optional rule reuse.

The `toulmin` library is available under MIT License.

---

## References

[1] Toulmin, S. *The Uses of Argument*. Cambridge University Press, 1958.

[2] Nute, D. "Defeasible Reasoning." In *Handbook of Logic in Artificial Intelligence and Logic Programming*, Vol. 3, Oxford University Press, 1994.

[3] Amgoud, L., and Ben-Naim, J. "Ranking-based Semantics for Argumentation Frameworks." *SUM 2013*, LNCS 8078, 2013.

[4] Gabriel, V.O., Panisson, A.R., Bordini, R.H., Adamatti, D.F., Billa, C.Z. "Reasoning in BDI agents using Toulmin's argumentation model." *Theoretical Computer Science*, 805:76–91, 2020.

[5] Modgil, S., and Prakken, H. "The ASPIC+ Framework for Structured Argumentation: A Tutorial." *Argument and Computation*, 5(1):31–62, 2014.

[6] Dung, P.M. "On the Acceptability of Arguments and its Fundamental Role in Nonmonotonic Reasoning, Logic Programming and n-person Games." *Artificial Intelligence*, 77(2):321–357, 1995.

[7] Ahrendt, W., et al. "StaRVOOrS: A Tool for Combined Static and Runtime Verification of Java." *Runtime Verification*, LNCS 9333, 2015.

[8] Ahrendt, W., et al. "A Specification Language for Static and Runtime Verification of Data and Control Properties." *FM 2015*, LNCS 9109, 2015.

[9] Cockburn, A. "Hexagonal Architecture (Ports and Adapters)." 2005.

[10] Gödel, K. "Über formal unentscheidbare Sätze der Principia Mathematica und verwandter Systeme I." *Monatshefte für Mathematik und Physik*, 38(1):173–198, 1931.

[11] Kelsen, H. *Reine Rechtslehre (Pure Theory of Law)*. Franz Deuticke, Wien, 1934.

