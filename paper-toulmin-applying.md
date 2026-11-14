# Claim, Not Fact: Applying Toulmin's Argumentation Model to Software Rule Engines

© 2026 PARK JUN WOO(mail@parkjunwoo.com). This paper is licensed under the Creative Commons Attribution 4.0 International License (CC BY 4.0), https://creativecommons.org/licenses/by/4.0/. The software in this repository is licensed separately under the MIT License.

## Abstract

Software rule engines (Rego/OPA, Drools, Semgrep, JSON Schema, and others) typically treat the data being validated as "facts." This paper argues that validation targets are not facts but **claims**, and that Toulmin's argumentation model (1958), although rarely used in rule engine design, is a suitable design foundation for it.

We present `toulmin`, a rule engine that implements this view (this edition describes v0.9.0). It maps Toulmin's six elements (Claim, Ground, Warrant, Backing, Qualifier, Rebuttal) onto engine components, uses Nute's strict/defeasible/defeater classification for rule strength, and computes verdicts with a weighted h-Categoriser rescaled to [-1, 1]. A rule is a function that takes the judgment context and its judgment criteria (spec) and returns a rule verdict v ∈ [-1, 1] together with evidence; rules form a graph connected by attack edges. The v0.9.0 verdict computation extends the h-Categoriser in four directions: an applicability determined by the rule verdict; a distinction between rebuttals, which enter the denominator, and exceptions, which multiplicatively reduce applicability; an attack force F ∈ [0, 1] on each attack edge; and crisp exclusion (Supersedes), which marks a node to be dropped from downstream tallies without changing any verdict. The same computation exists in a Go reference implementation and in TypeScript, Python, and C ports, and the markdown policy language TANGEUL provides an optional declarative surface over the core.

We examine the design through two case studies. We convert filefunc's 22 code structure rules into Toulmin warrants and report figures from three projects to which the rule set was applied, and we show a phase-optional design that unifies fullend's compile-time and runtime policies in a single rule function.

---

## 1. Introduction

### 1.1 Problem: The Fact Assumption

Modern software rule engines commonly treat their inputs as "facts." Drools inserts Java objects into working memory as facts. Rego evaluates policies against an `input` document taken as given. JSON Schema checks whether a document (an instance) conforms to a schema, taking the document as given.

This terminology is epistemologically inaccurate. A **fact** is something already established as true; it needs no verification. Yet the reason a rule engine exists is to **verify** whether input data satisfies rules. Calling a verification target "already true" is a contradiction.

Verification targets are **claims**: assertions that may be true or false and whose validity must be judged against rules. We hold that this is more than a terminological correction. In §3.1 we put forward the hypothesis that the "fact" assumption is entangled with three structural limitations of existing rule engines.

### 1.2 The Missing Link

Toulmin analyzed the structure of argumentation into six elements in 1958 [1]: Claim, Ground, Warrant, Backing, Qualifier, and Rebuttal. Many rule engines appeared afterwards (CLIPS 1985, Jess 1995, Drools 2001, OPA/Rego 2016). They share a fact-based design, and within the scope of our survey we found no rule engine that takes Toulmin's model as its design foundation.

One possible reason is that Toulmin's work has been read in philosophy and rhetoric and did not become a reference for rule engine development. Meanwhile, individual rule engines acquired mechanisms that can be matched with parts of Toulmin's model: LegalRuleML's defeasibility elements (corresponding to Rebuttal), OPA's `# METADATA` annotations (corresponding to part of Backing), and Drools' salience (corresponding to part of Qualifier).

### 1.3 Contributions

1. **Concept**: We present the view of validation targets as claims and map Toulmin's six elements onto rule engine components (§3).
2. **Formal definition**: We define a verdict computation that extends the weighted h-Categoriser [3] with rule verdicts and applicability, a distinction between rebuttals and exceptions, attack force, and reporting-only exclusion (§4).
3. **Implementation**: We describe `toulmin`, which implements this computation in a Go reference implementation and TypeScript, Python, and C ports, together with the declarative surface TANGEUL (§5).
4. **Case studies**: We examine the design through the conversion of filefunc's 22 rules (§6) and fullend's phase-integration design (§7).

---

## 2. Background and Related Work

### 2.1 Toulmin's Argumentation Model

Toulmin [1] analyzed arguments into six elements.

- **Claim**: The assertion being made.
- **Ground (Data)**: The material offered in support of the claim.
- **Warrant**: The rule or principle that licenses the step from ground to claim.
- **Backing**: The justification for why the warrant is valid.
- **Qualifier**: The degree of certainty ("certainly," "presumably," "possibly").
- **Rebuttal**: Exceptional conditions under which the claim does not hold.

Verheij [19] formalized Toulmin's scheme for the evaluation of arguments. Gabriel et al. [4] implemented five of these elements (excluding Backing) in a BDI agent system using AgentSpeak/Jason, with a qualify function computing SWW − SWR (sum of warrant weights minus sum of rebuttal weights). Within the scope of our survey we found no prior work that uses Toulmin's model as the design foundation of a general-purpose software rule engine.

### 2.2 Defeasible Reasoning

Nute [2] classified rules into three strengths: **strict** (cannot be defeated), **defeasible** (can be defeated by contrary evidence), and **defeaters** (block the conclusions of other rules without asserting their own). Defeasible logic in Nute's line also includes a superiority relation among rules [18].

Pollock [12] distinguished two kinds of defeaters. A **rebutting defeater** supports the negation of the conclusion; an **undercutting defeater** attacks the connection between the grounds and the conclusion, stating that it does not hold in the present situation. The structured argumentation framework ASPIC+ [5] distinguishes strict from defeasible rules and handles both rebutting and undercutting attacks, with its semantics given on top of Dung's abstract argumentation framework [6].

### 2.3 Abstract Argumentation and Gradual Semantics

Dung's abstract argumentation framework [6] treats arguments as nodes and attacks as edges and determines sets of acceptable arguments (extensions). Labelling semantics [14] assigns each argument an IN, OUT, or UNDEC label. In these semantics, judgments are crisp.

Gradual semantics assigns each argument a continuous acceptability degree. The h-Categoriser was introduced by Besnard and Hunter [15], and Amgoud et al. [3] defined the weighted h-Categoriser with an initial weight per node:

```
Acc(a) = w(a) / (1 + Σ_{b ∈ Att(a)} Acc(b))
```

where `w(a) ∈ [0, 1]` is the initial weight and `Att(a)` is the set of nodes attacking `a`. As the formula shows directly, when an attacker's acceptability falls (for example, because the attacker is itself attacked), the acceptability of `a` rises. For the definition and properties on general graphs with cycles, see [3]. Amgoud and Ben-Naim [16] set out principles for ranking-based semantics. The idea of attaching weights to attack edges appears in the weighted argument systems of Dunne et al. [17], where the weights determine which attacks may be disregarded within an inconsistency budget.

### 2.4 Existing Rule Systems

| System | Input model | Defeasibility | Metadata | Phase |
|--------|------------|---------------|----------|-------|
| Rego/OPA | fact (`input`) | `default`/`else` patterns | `# METADATA` (not used in the decision) | runtime policy evaluation; configuration files via Conftest |
| Drools | fact (working memory) | salience/agenda (firing order) | rule attributes | runtime |
| Semgrep | source code patterns | none | rule `metadata` (not used in the decision) | static analysis |
| JSON Schema | instance (JSON document) | none | annotation keywords such as `description` | whenever the validator is called (static or runtime) |
| LegalRuleML | legal norm representation (XML) | elements added on top of the core model | legal metadata | representation language (no execution phase of its own) |

Among the systems we examined, none combines treating inputs as claims with first-class defeasibility, spec and qualifier that take part in the judgment, and phase-optional evaluation.

### 2.5 Phase Integration

StaRVOOrS [7] and ppDATE [8] showed that a single specification can feed both static verification (KeY) and runtime monitoring (LARVA). Cockburn's Hexagonal Architecture [9] set out the adapter pattern that separates domain logic from external systems.

---

## 3. From Facts to Claims

### 3.1 "Fact" Naming and Design Limitations: Three Hypotheses

Structural limitations in rule engines have many causes: technology stack constraints, performance optimization, and use case scope all play a part. This section adds the hypothesis that the name "fact" acts as a cognitive frame that makes certain design directions feel natural. The three pathways below are conceptual arguments and have not been tested empirically (§8.2).

**Pathway 1: Fact → weakened verification awareness → consistency gaps.**
When Drools calls working memory entries "facts," developers tend to regard them as already verified. Cross-validation between facts happens only when developers explicitly add rules for it. Java's object model also plays a part, but the name "fact" provides a frame in which omitting verification feels natural.

**Pathway 2: Fact → binding to data structure → rule duplication.**
Rego receives `input` as a JSON document, and rules reference concrete paths (`input.user.role`). Rego's JSON query syntax is one factor, but fact-based design encourages binding rules to concrete data shapes rather than abstract properties. As a result, the same logical rule ("only admins can delete") must be rewritten for HTTP requests, gRPC metadata, and CLI arguments.

**Pathway 3: Fact → fixed truth values → defeasibility as an afterthought.**
Facts are true, and true things are not overturned. Systems designed on this premise must attach external exception mechanisms for defeasibility. LegalRuleML adding XML elements such as `<Override>`, `<Defeater>`, and `<Superiority>` outside its core model is an example. Redefining inputs as claims makes designs with built-in defeasibility easier, but the structural devices of rebuttals and exceptions (§4) are still required.

**Precedent: JWT claims.**
JWT (RFC 7519 [20]) calls token fields "claims," not "facts." `sub`, `exp`, and `iss` are assertions by the token issuer. Signature verification, expiration checks, and issuer validation fit well with the premise that this data consists of assertions. It is an example of a name and a verification obligation going together.

### 3.2 Toulmin Mapping

| Toulmin element | Counterpart in the engine | Cardinality |
|---|---|---|
| Claim | The proposition being judged. One per warrant node, identified by the node name (`TraceEntry.Name`) | 1 per node |
| Ground | Evidence data. Passed as the context `ctx` and filled by a Ground Adapter (`TraceEntry.Ground`) | 0..N |
| Warrant | A rule function and its node. Can be attacked | 1..N |
| Backing | Spec: judgment criteria fixed at declaration time (`TraceEntry.Specs`) | 0..N per rule |
| Qualifier | A weight q ∈ [0, 1] per rule **(§3.3)** | 1 per rule |
| Rebuttal | Attacking nodes: rebuttals (Counter) and exceptions (Except) **(§4.4)** | 0..N |

Input data enters the context as ground, but the engine does not treat the proposition being judged as established. How far the proposition holds comes out only as a verdict, after rule verdicts, attacks, and qualifiers have been taken into account.

### 3.3 Repositioning the Qualifier: From Claim to Rule

In Toulmin's original model, the Qualifier is attached to the **Claim**. In "**Presumably**, this patient should receive penicillin," the Qualifier is a modal qualifier expressing the degree of certainty the speaker assigns to the claim after weighing Ground, Warrant, and Rebuttal together.

This model moves the Qualifier from the Claim to **each rule node (general rule, rebuttal, exception)**. It is an engineering adjustment that arises in applying argumentation theory to rule engines.

In a rule engine, the observation "this file contains 3 functions" is ground, and the claim to be judged is a proposition such as "this file violates the one-function-per-file rule." There is rarely reason to attach a degree of certainty to the observation itself. What determines the quality of the judgment is the **certainty of the rule** that leads from the observation to the claim.

- "One func per file": qualifier 1.0 (certain rule)
- "Recommended max 100 lines": qualifier 0.7 (flexible rule)
- "Allergy possibility": qualifier 0.95 (strong rebuttal)

This repositioning aligns naturally with the weighted h-Categoriser [3]. The rule's qualifier q takes the place of the node's initial weight `w(a)`, and the final output, the verdict, assumes the role that the Qualifier played in Toulmin's model: the degree of certainty of the judgment. In v0.9.0 the rule function also yields a degree (the rule verdict v, §4.2). q expresses how trustworthy the rule is, v how far the rule applies to the case at hand, and the two combine as a product (§4.3).

| | Toulmin's original | This model |
|--|-------------------|------------|
| Qualifier position | Claim (output side) | Each rule node (input side) |
| Role | Certainty of the final judgment | The rule's weight q |
| Certainty computation | Speaker judges holistically | Engine computes from q, rule verdict v, attacks and their force |
| Final certainty | Qualifier itself | **Verdict** |

### 3.4 Epistemic Status of Warrants

Strictly speaking, warrants are also claims. "Only admins can delete" is an organizational agreement, not a physical law, and even physical laws are claims with very high reproducibility. Applying this recognition without end, however, leads to an infinite regress of justification. Even in formal systems, axioms are presupposed rather than proved within the system. Gödel's second incompleteness theorem [10] shows that a consistent formal system containing arithmetic cannot prove even its own consistency within the system. Kelsen's basic norm (Grundnorm) [11] exhibits a similar structure: the chain of validity in a legal system stops at a presupposed norm.

This model therefore takes a pragmatic decision: **warrants are treated as facts within the system.** This is a premise that axiomatic systems commonly carry, and it bounds the model's scope to "judging the validity of claims given warrants."

---

## 4. The Judgment Model

This section defines the verdict computation of `toulmin` v0.9.0. The definition follows the Go reference implementation (`go/pkg/toulmin`) and its tests, and the four ports implement the same computation (§5.5).

### 4.1 Judgment Graph

A judgment graph consists of the following.

- A node set N. Each node x has a rule function `f_x`, a list of specs `σ_x`, a qualifier `q_x ∈ [0, 1]` (default 1), a role `role(x) ∈ {rule, counter, except}`, and a strength `str(x) ∈ {defeasible, strict, defeater}`. Nodes registered with `Rule` and `Counter` are defeasible; nodes registered with `Except` are defeaters.
- An attack edge set A ⊆ N × N. An edge (y, x) means "y attacks x" and has a force `F_yx ∈ [0, 1]` (default 1).
- A supersede edge set P ⊆ N × N. An edge (y, x) means "y supersedes x."

A and P must each be acyclic. At the start of evaluation the engine runs a depth-first search over each edge set separately and, if it finds a cycle, returns an error without executing any rule function.

Let `Att(x) = {y | (y, x) ∈ A}` be the attackers of node x, split into two parts:

- exceptions `E(x) = {y ∈ Att(x) | str(y) = defeater or role(y) = except}`
- rebuttals `R(x) = Att(x) \ E(x)`

Whether an attack is computed as a rebuttal or as an exception is thus decided by **the attacker's role**, not by the edge. A node that is neither an attacker nor a defeater is a **warrant**; evaluation produces one result per warrant.

### 4.2 Rule Verdict and Applicability

In a context c (`ctx`), a rule function returns a rule verdict `v_x = f_x(c, σ_x) ∈ [-1, 1]` and evidence. +1 means "fully applies," 0 "half applies," and −1 "does not apply." If the rule verdict is NaN or outside [-1, 1], or if the rule function panics (throws), the whole evaluation ends with an error.

The applicability of node x is `a(x) = (v_x + 1) / 2 ∈ [0, 1]`. If `v_x > -1`, we say x **applies** (`TraceEntry.Activated`).

Legacy bool rule functions are still accepted: true becomes `v = +1`, false becomes `v = -1`. It matters that "does not apply" is −1, not 0. A rule that returns 0 stays in the computation with applicability 0.5, and since 0 is a legal value the engine cannot detect this mistake.

### 4.3 Verdict Computation

The verdict `verdict(x)` is defined recursively.

- If `v_x = -1`, then `verdict(x) = -1`. The attackers of x are not evaluated (lazy evaluation).
- Otherwise

```
S(x)       = Σ_{y ∈ R(x)} F_yx · (verdict(y) + 1) / 2         rebuttals: denominator term
G(x)       = Π_{y ∈ E(x)} (1 − F_yx · (verdict(y) + 1) / 2)    exceptions: applicability gate
a_eff(x)   = a(x) · G(x)                                       effective applicability
raw(x)     = q_x · a_eff(x) / (1 + S(x))                       ∈ [0, 1]
verdict(x) = 2 · raw(x) − 1                                    ∈ [-1, 1]
```

- If `str(x) = strict`, all attacks are ignored: `S(x) = 0`, `G(x) = 1`.

An attacker y that does not apply has `verdict(y) = -1`, so its contribution `(verdict(y) + 1) / 2` is 0. Since `q_x` and `a_eff(x)` lie in [0, 1] and `S(x) ≥ 0`, we have `raw(x) ∈ [0, 1]` and `verdict(x) ∈ [-1, 1]`. The effective applicability `a_eff` of each warrant is reported as `EvalResult.Applicability`.

**Relation to the h-Categoriser.** If every rule is bool (v ∈ {−1, +1}), there are no exception attacks, and every force is 1, then for an applying node x

```
raw(x) = q_x / (1 + Σ_{y ∈ Att(x), v_y = +1} raw(y))
```

that is, the weighted h-Categoriser [3] on the subgraph of applying nodes, with qualifiers as weights. A repository test (`graph_evaluate_bool_compat_test.go`) confirms that in this case both a bool rule graph and a graded rule graph returning ±1 produce values bit-for-bit equal to the pre-v0.9.0 formula `raw = q / (1 + Σ raw(attackers))`.

### 4.4 Rebuttals, Exceptions, and Attack Force

Rebuttals and exceptions are declared with the same attack edge but computed differently.

- A **rebuttal (Counter)** is a claim against the conclusion. It is added to the denominator `1 + S(x)` and lowers the verdict. A single full rebuttal brings an unattacked full claim (q = 1, v = +1) to 0 (balance) without overturning it, and rebuttals alone never bring the verdict to −1 as long as `q_x · a_eff(x) > 0`.
- An **exception (Except)** says "the rule does not apply in this case." It multiplies the target's applicability instead of adding to the denominator. A single full exception brings the target's effective applicability to 0 and its verdict to −1.

This distinction can be seen as corresponding to Pollock's rebutting/undercutting distinction [12]. A rebuttal contests the conclusion; an exception states that the connection from grounds to conclusion does not apply to this case. An exception to an exception is likewise one edge: an exception node can itself be attacked, and as its verdict falls, the gate closes less.

The force `F` is a property of the edge that multiplies the attacker's contribution; it is a different axis from a node's strength or qualifier. The table below gives the verdict of an unattacked full claim (q = 1, v = +1) attacked with force F by a single fully applying attacker (v = +1). In closed form, the rebuttal gives `(1 − F)/(1 + F)` and the exception gives `1 − 2F`.

| F | 0 | 0.25 | 0.5 | 0.75 | 1 |
|---|---|---|---|---|---|
| Rebuttal | 1 | 0.6 | 0.333 | 0.143 | 0 |
| Exception | 1 | 0.5 | 0 | −0.5 | −1 |

The values were obtained by building these graphs with the v0.9.0 Go core and calling `Evaluate` (rounded to three decimals); the rebuttal values for F = 0, 0.5, 1 and the exception values for F = 0, 0.25, 0.5 match those checked by a repository test (`graph_evaluate_attack_force_test.go`). Multiple exceptions compound multiplicatively: three half-applying exceptions (v = 0.5) leave only `(1 − 0.75)³ ≈ 0.016` of the applicability (`graph_evaluate_except_gate_test.go`), whereas three rebuttals bring the same claim down only to −0.5.

A force that is NaN or outside [0, 1] is rejected at declaration time; the qualifier follows the same rule.

The idea of weighting attack edges also appears in weighted argument systems [17], but here the force is not a budget that decides whether to disregard an attack; it is a coefficient that multiplies the sum term and the gate of the h-Categoriser.

### 4.5 Strength as Graph Constraint

Nute's classification [2] works in the engine as follows.

| Strength | Meaning in the engine |
|----------|---------|
| Strict | All declared attacks are ignored in the computation (`S = 0`, `G = 1`) |
| Defeasible | Attacks are taken into account |
| Defeater | Has no conclusion of its own and never becomes a warrant; as an attacker it is computed as an exception gate, and it can itself be attacked |

The public builder API of v0.9.0 has no method that marks a node as strict (§8.2). A defeasible rule that receives no attack edge obtains the same verdict as a strict rule, so the strict classification in the case study (§6) appears in the engine as "a rule with no incoming attack edges."

### 4.6 Supersedes: A Reporting-Only Crisp Priority

In addition to attack edges, the engine provides **supersede (Supersedes)** edges. For the supersede edge set P,

```
Excluded(x)  ⇔  ∃ y : (y, x) ∈ P  and  a_eff(y) > 0
```

The superseder's effective applicability reflects exception gates but not rebuttals, and a strict superseder has `a_eff = a`. A superseder that does not apply has `a_eff = 0`. Even a partial effective applicability such as 0.3 excludes the whole target.

Exclusion does not change any verdict. An excluded node is only marked `Excluded` in the trace: its verdict stays the same, and the attacks and exclusions it exerts on other nodes remain in force. Supersede edges never enter the h-Categoriser traversal. Dropping a marked node from a tally is the job of downstream aggregation (for example, choosing a single discount).

Consider an example from the repository tests (`run_supersedes_except_gate_test.go`). A VIP discount supersedes the general discount. If a "suspended VIP" exception fully gates the VIP discount (`a_eff = 0`), the exclusion is lifted and the general discount stays in the tally. An exception with force 0.5 (`a_eff = 0.5`) or a rebuttal against the VIP discount (`a_eff = 1`) keeps the exclusion. A cycle of supersede edges is treated as an authoring error (mutual exclusion) and rejected before evaluation.

Theoretically, exclusion is a device at a different level from attack. Raz's exclusionary reasons [13] are not first-order reasons that outweigh others on the scale but second-order reasons not to put certain reasons on the scale at all, and supersede edges follow this structure of exclusion from the tally. The superiority relation of defeasible logic [18] likewise expresses crisp priority among rules. Unlike labelling semantics [14], however, where an argument labelled OUT loses its power to attack, an excluded node keeps both its attacks and its power to exclude. Exclusion is a reporting-only crisp layer placed on top of the gradual computation, which it leaves untouched.

### 4.7 Interpreting Verdicts

| Verdict | Meaning |
|---------|---------|
| +1 | Fully applies (`q · a_eff = 1`) with no effective rebuttal |
| 0 | Balance; e.g., one full claim against one full rebuttal, or `q · a_eff = 1/2` with no attacks |
| −1 | Effective applicability 0 (a full exception) or q = 0 |

The verdict is **the degree to which the claim holds**. The positive side is affirmative (the claim is supported), the negative side negative (the claim is rebutted). A warrant that does not apply (`v = -1`) is left out of the results instead of producing a verdict of −1, so "does not apply" is distinguished from "rebutted." The engine only computes values; domain-specific interpretation such as thresholds is up to the user (e.g., access control allows when `verdict > 0`). For intuitive results we recommend phrasing rule claims affirmatively ("access allowed," where +1 means allowed). Negatively phrased rules ("violation detected") are valid, but +1 then means a violation, which reads counter-intuitively.

---

## 5. Implementation

### 5.1 Rules as Functions Returning Verdicts and Evidence

All rules (general rules, rebuttals, exceptions) share the same signature. In Go:

```go
func(ctx toulmin.Context, specs toulmin.Specs) (float64, any) // rule verdict v ∈ [-1, 1], evidence
func(ctx toulmin.Context, specs toulmin.Specs) (bool, any)    // legacy signature: true → +1, false → -1
```

The first parameter, `ctx`, is a `Context` interface with `Get(key string) (any, bool)` and `Set(key string, value any)`; the second, `specs`, carries the rule's judgment criteria as runtime values. The second return value is domain-specific evidence (e.g., the location of a violation and the expected value); the engine does not interpret it and passes it to results and trace. An optional convention type `Fact{Where, Expected, Actual}` holds location, expected value, and actual value.

| | Context | Spec |
|--|---------|------|
| Toulmin element | Ground | Backing |
| Definition | Facts about the judgment target (Get/Set access) | Judgment criteria of the rule |
| Timing | Varies per request | Fixed at declaration time |
| Example | Number of functions in a file, user's role | Threshold for allowed functions, role name |

```go
// Binary rule: written with the bool signature.
func CheckOneFileOneFunc(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
    gf, _ := ctx.Get("file")
    f := gf.(*FileGround)
    if len(f.Funcs) > 1 {
        return true, &Evidence{Got: len(f.Funcs), Expected: 1}
    }
    return false, nil
}

// Rule that uses its spec as the judgment criterion
func CheckLineCount(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
    gf, _ := ctx.Get("file")
    f := gf.(*FileGround)
    limit := specs[0].(*LineLimit)
    if f.Lines > limit.Max {
        return true, &Evidence{Got: f.Lines, Limit: limit.Max}
    }
    return false, nil
}

// Graded rule: maps a classifier score p ∈ [0, 1] to the rule verdict v = 2p − 1.
func LooksLikeSpam(ctx toulmin.Context, specs toulmin.Specs) (float64, any) {
    p, _ := ctx.Get("spam_score")
    score := p.(float64)
    return 2*score - 1, score
}
```

The distinction between warrant, rebuttal, and exception lies not in the function signature but in the node's position in the graph and its registered role.

| Role | Graph position |
|------|---------------|
| Warrant (`Rule`) | A node that can be attacked and attacks no other node. Produces a result |
| Rebuttal (`Counter`) | A node that attacks another node (has its own conclusion). Computed as a denominator term |
| Exception (`Except`) | A node that attacks without a conclusion of its own. Computed as an applicability gate |

### 5.2 Graph API: Separating Logic from Structure

The Graph API separates rule logic (functions) from graph structure (roles, qualifiers, attacks, exclusions).

```go
g := toulmin.NewGraph("api:access")
auth    := g.Rule(isAuthenticated)
blocked := g.Counter(isIPBlocked)
exempt  := g.Except(isInternalIP)
blocked.Attacks(auth)        // rebuttal, force 1
exempt.Attacks(blocked)      // exception to a rebuttal: reduces the blocking rule's applicability
hint := g.Counter(isSuspicious)
hint.Attacks(auth, 0.3)      // weak rebuttal with force 0.3
```

`g.Rule(fn)`, `g.Counter(fn)`, and `g.Except(fn)` return a `*Rule` reference on which the builder methods `.With(spec)` (adds a spec) and `.Qualifier(q)` (default 1.0) can be chained. A spec must implement the `Spec` interface (`SpecName() string`, `Validate() error`), and specs with function fields are rejected. `Attacks(target, force...)` declares an attack edge, and `Supersedes(target)` a supersede edge.

Registering the same function with different specs yields distinct nodes. The node identifier is the function identifier followed by the names and serialized values of its specs.

```go
g := toulmin.NewGraph("access")
admin  := g.Rule(isInRole).With(&RoleSpec{Role: "admin"})
editor := g.Rule(isInRole).With(&RoleSpec{Role: "editor"}).Qualifier(0.8)
```

Functions are identifiers, so no string names are needed. The same function can be reused across graphs with different roles and attack relations:

```go
// Same IsAdult function, different graphs, different attack relations
votingGraph := toulmin.NewGraph("voting")
vw := votingGraph.Rule(IsAdult)
vr := votingGraph.Counter(HasCriminalRecord)
vr.Attacks(vw)

contractGraph := toulmin.NewGraph("contract")
cw := contractGraph.Rule(IsAdult)
cr := contractGraph.Counter(IsBankrupt)
cr.Attacks(cw)
```

Exclusion is declared as follows:

```go
general := g.Rule(generalDiscount)
vip     := g.Rule(vipDiscount)
vip.Supersedes(general)               // marks general Excluded while vip has a_eff > 0
g.Except(suspendedVIP).Attacks(vip)   // fully gating vip lifts the exclusion
```

All structural metadata (role, qualifier, attacks, force, exclusion) belongs to the graph, not to the function; functions contain only judgment logic.

### 5.3 Evaluation Procedure

`Evaluate` judges in the following order:

```
0. Cycle detection: DFS over attack edges and supersede edges separately.
   A cycle is an error (before any rule function runs)
1. Start from each warrant
2. Run the rule function f(ctx, specs) to obtain the rule verdict v
   (panic, NaN, or out of range is an evaluation error)
3. If v = -1, the verdict is -1 and the attackers are not evaluated.
   A warrant in this state is left out of the results
4. If v > -1 and the node is not strict, evaluate the attackers recursively;
   exceptions go into the gate G, rebuttals into the sum S
5. a_eff = a·G, raw = q·a_eff/(1+S), verdict = 2·raw − 1. Cached per node
6. With trace off, each node's function runs at most once within one Evaluate
   and results are shared across warrants. With trace on, state is reset and
   re-evaluated per warrant, so a function may run more than once
7. With trace on, an exclusion pass runs per warrant: each superseder is evaluated
   and, if its a_eff > 0, the target's TraceEntry is marked Excluded (verdict unchanged)
```

Because of lazy evaluation, `Evaluate` with trace off executes only the rules reachable from a warrant along attack chains.

### 5.4 Results, Trace, and Execution

```go
results, err := g.Evaluate(ctx)                                    // verdict + evidence
results, err = g.Evaluate(ctx, toulmin.EvalOption{Trace: true})    // + per-node trace
results, err = g.Evaluate(ctx, toulmin.EvalOption{Duration: true}) // + execution time (trace auto-enabled)
```

```go
type EvalResult struct {
    Name          string       `json:"name"`
    Verdict       float64      `json:"verdict"`
    Evidence      any          `json:"evidence,omitempty"`
    Applicability float64      `json:"applicability"` // a_eff
    Trace         []TraceEntry `json:"trace"`
}

type TraceEntry struct {
    Name        string        `json:"name"`         // = Claim
    Role        string        `json:"role"`
    Activated   bool          `json:"activated"`    // v > -1
    RuleVerdict float64       `json:"rule_verdict"` // rule verdict v
    Qualifier   float64       `json:"qualifier"`
    Verdict     float64       `json:"verdict"`
    Evidence    any           `json:"evidence,omitempty"`
    Ground      any           `json:"ground,omitempty"` // = Ground (ctx)
    Specs       Specs         `json:"specs,omitempty"`  // = Backing
    Duration    time.Duration `json:"duration,omitempty"`
    Excluded    bool          `json:"excluded,omitempty"`
    Level       Level         `json:"level"`
}
```

For each warrant, the trace of `Evaluate` records the role, rule verdict, applies flag, qualifier, spec, evidence, and exclusion mark of every evaluated node. Because spec is part of the trace, nodes that register the same function with different specs can be told apart by their judgment criteria. Being able to audit how a verdict was derived in this way is an advantage of symbolic reasoning. The trace entries of the `Evaluate` path, however, do not carry each node's final verdict (`Verdict`); a whole-graph trace with per-node verdicts and exclusion marks is returned by `Run`. The results of `Evaluate` with trace off carry no exclusion marks. The default judgment method `Matrix` is a lazy recursive h-Categoriser; `Recursive` is not yet implemented and returns an error if selected.

`Evaluate` has no side effects. `Run` evaluates the whole graph, then runs, in registration order, the handlers (`RunOn`) registered on nodes that apply and have a positive verdict, and for a node with a declared sub-graph it `Run`s that graph with the same context (execution composition, cycles rejected, depth capped at 64). The Go port can also attach a reporting level `Level` (`Unset`/`Fail`/`Review`), separate from the verdict, to each node. To record why a rule exists, use ordinary comments; there is no special annotation syntax.

### 5.5 Four Ports

The same verdict computation is implemented in four directories of the repository: `go/` (reference implementation, module `github.com/park-jun-woo/toulmin/go`), `ts/`, `py/`, and `c/`. The Go module depends on the GEUL library `github.com/geul-org/gogeul`, which is used by the TANGEUL toolchain.

| | Go | TypeScript | Python | C |
|---|---|---|---|---|
| Graded rule | `(float64, any)` | `[number, unknown]` | `(float, Any)` | `rc_verdict_func` + `rc_graph_*_verdict` |
| Telling bool rules apart | by function type at registration | by return value type on each call | by return value type on each call | separate registration functions (`rc_graph_rule` etc.) |
| Attack force | `Attacks(t, f)` | `attacks(t, f)` | `attacks(t, force=f)` | `rc_rule_attacks_with_force(r, t, f)` |
| Exclusion | `Supersedes(t)` | `supersedes(t)` | `supersedes(t)` | `rc_rule_supersedes(r, t)` |
| Effective applicability | `EvalResult.Applicability` | `applicability` | `applicability` | `rc_eval_result.applicability` |
| Invalid force or qualifier | panic at declaration | throws | `ValueError` | deferred builder error (`RC_ERR_BUILDER` at evaluation); force also returns `RC_ERR_INVALID_ARG` from the call |
| Invalid rule verdict | `Evaluate` error | throws | `RuntimeError` | `RC_ERR_RULE` |

The verdict computation (§4.3) and the exclusion test (§4.6) are the same in all four ports. The differences are these. In TypeScript and Python, a return value that is neither a number nor a boolean is an error, not "false." A C rule function returns a status code, and a failure becomes `RC_ERR_RULE`. With trace on in `Evaluate`, the role shown in an entry is the registered role in Go but a role inferred from graph position in the other three ports (the label differs when a node registered as a general rule attacks another node; verdicts are not affected). The reporting level `Level`, the TANGEUL toolchain, and the domain framework packages exist only in the Go port.

### 5.6 Declarative Surface: TANGEUL

The core engine requires no separate DSL; rules are functions in each language. On top of it, optionally, sits TANGEUL (Toulmin Amgoud Nute Graph - General Encoding Unified Language), a markdown policy language that non-developers can read and audit (Go port `go/pkg/tangl`, `go/pkg/tangeul`, CLI `tangl`). A TANGEUL document compiles to a `toulmin.Graph` without a separate runtime.

- **Rules and cases are separate.** The `tangl:Rules` section owns the implementation of a judgment (an expression rule, or a `uses` rule that calls a host function) and its inputs (`with` specs, `grounded on` grounds). The `tangl:Cases` section cites the same name and adds only the role (general, counter, except), the qualifier, and the reporting level.
- **Definitions.** A string-list definition (`` `domains` means ["a.com", "b.com"] ``) becomes a built-in spec with exact membership checking, and a file definition (`` `domains` reads `JSON list` from "./domains.json" ``) is read once and frozen before evaluation or code generation.
- **Attacks and force.** `` don't `X` when `Y` `` is an attack of force 1, and `` `Y` attacks `X` with force 0.5 `` compiles to `Y.Attacks(X, 0.5)`. The force lies in [0, 1] with at most four decimal places. The same attack statement is computed as a gate if the attacker is an except rule and as a denominator term otherwise (§4.4).
- **Exclusion.** `` `X` supersedes `Y` `` is a supersede edge, and `` supersede in order: `A` > `B` > `C` `` expands to every earlier→later pair (A→B, A→C, B→C).
- **The source of truth is `.geul`.** `tangl encode` absorbs an authored `.md` surface into a locale-independent `.geul` file, that is, one GEUL document following the TOULMIN codec [21]. The file holds both the AST and the authored source bytes, and the attack force is carried in the `F_FORCE` slot of the ATTACKS record. Checking, effect analysis, and code generation run on this stream, and `tangl decode` either restores the source byte for byte (mode 1) or regenerates an English or Korean surface from the AST (mode 2). The English and Korean surfaces parse to the same AST.

TANGEUL expression rules and case composition (`checking`) yield only true or false (+1/−1). Graded judgments arise only when a `(float64, any)` host function is used through `uses`, and the reference evaluator `pkg/tangl/eval` is still bool-only. A `checking` that refers to a case in another document is accepted by the parsers and the code generator but cannot yet be stored in `.geul`.

---

## 6. Case Study 1: Toulmin Conversion of filefunc's Rules

### 6.1 Background

filefunc is a code structure convention and CLI tool for LLM-native Go development. At the time of conversion, filefunc implemented 22 rules (F, Q, A, and C families) as Go functions evaluated at compile time on the AST, and all of them were converted into Toulmin warrants. The classification and figures below refer to the rule set at the time of conversion.

### 6.2 Strength Classification

| Strength | Count | Ratio | Examples |
|----------|-------|-------|----------|
| Strict | 15 | 68% | F1–F4, A1–A3, among others |
| Defeasible | 4 | 18% | Q1, Q2, Q3, C4 |
| Defeater | 3 | 14% | F5, F6, test file exceptions |

Most rules are strict, because code structure conventions inherently try to minimize exceptions. Defeasible rules are those with context-dependent thresholds (depth limits Q1–Q2, size limit Q3). As noted in §4.5, strict rules appear in the engine as rules that receive no attack edges, and defeaters as `Except` nodes.

### 6.3 Results of Applying the Rule Set

| Project | Files (before→after) | Avg LOC/file (before→after) | SRP violations resolved | Depth violations resolved |
|---------|---------------------|---------------------------|------------------------|--------------------------|
| filefunc | — (born compliant) | 25.1 | 0 | 0 |
| fullend | 87→1,260 | 244→25.4 | 66→0 | 148→0 |
| whyso | 12→99 | 147.8→24.4 | 12→0 | 23→0 |

These figures were measured by the author before and after bringing each project into compliance with the filefunc rules. They show that a rule set expressed as a Toulmin graph could be applied to real codebases and drive violations to zero, but they do not isolate the effect of the engine design (the Toulmin conversion) from an if-else implementation (§8.2).

### 6.4 Benefits of Toulmin Conversion

1. **Spec as a first-class value**: Judgment criteria (thresholds, policies) are passed to rule functions as runtime spec values. The same function registered with different specs becomes distinct rules.
2. **Structured rebuttals and exceptions**: Exception conditions (F5 test files, F6 grouped constants) are declared as explicit attack edges rather than buried in if-branches.
3. **Strength classification**: "This rule has no exceptions" (strict) and "this rule may be defeated" (defeasible) are visible in the graph structure.

---

## 7. Case Study 2: Phase-Optional Design in fullend

### 7.1 Background

fullend validates consistency across 9 SSOTs (OpenAPI, SSaC, Rego, etc.) and generates code. 18 cross-check rules are implemented as Go functions, and runtime policies are separate Rego files. Whereas filefunc (§6) operates in a single phase (compile time), fullend is a case that calls for a phase-optional design.

### 7.2 Problem: Duplicate Rules

The same logical rule is written twice (sketch):

**Compile time** (Go cross-check):
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

### 7.3 Toulmin Design

These are binary judgments, so they use the bool signature (§5.1).

```go
// OWASP API Security Top 10 (2023) API2: Broken Authentication
func AuthEndpointRequiresClaims(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
    ep, _ := ctx.Get("endpoint")
    g := ep.(*EndpointGround)
    if g.Security.Contains("bearerAuth") && !g.PolicyRule.References("claims") {
        return true, &Evidence{Endpoint: g.Path, Missing: "claims reference"}
    }
    return false, nil
}

// Public API endpoints do not require authentication by design
func PublicEndpointException(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
    ep, _ := ctx.Get("endpoint")
    g := ep.(*EndpointGround)
    return g.Annotation.Contains("x-public"), nil
}

// Graph definition
g := toulmin.NewGraph("endpoint-auth")
w := g.Rule(AuthEndpointRequiresClaims)
d := g.Except(PublicEndpointException)
d.Attacks(w)   // for a public endpoint: effective applicability 0 → verdict -1
```

**Compile-time Ground Adapter**: `EndpointGround.Security` ← OpenAPI spec parse, `EndpointGround.PolicyRule` ← Rego AST parse.

**Runtime Ground Adapter**: `EndpointGround.Security` ← request route middleware, `EndpointGround.PolicyRule` ← OPA evaluation.

One warrant function serves both compile-time cross-validation and runtime policy enforcement. The graph structure is the same; only the Ground Adapter changes. The function itself is reusable across graphs, so the same `AuthEndpointRequiresClaims` can take part in different attack relations in different contexts. This is a design-level case study and contains no quantitative evaluation of whether the two Ground Adapters yield the same judgment on the same input.

---

## 8. Discussion

### 8.1 What This Paper Does Not Claim

This paper does not claim theoretical originality for its components. The structural foundation comes from Toulmin (1958) [1], the strength classification from Nute (1994) [2], the backbone of the verdict computation from the h-Categoriser (Besnard and Hunter 2001 [15]; weighted version, Amgoud et al. 2017 [3]), the rebutting/undercutting distinction from Pollock (1987) [12], and the adapter pattern from Cockburn (2005) [9] and StaRVOOrS (2015) [7]. Rule verdicts and applicability, the exception gate, attack force, and reporting-only exclusion are engineering extensions added by this engine, and their formal properties are not analyzed in this paper. Our contribution is **showing that these independently developed pieces connect**, that is, presenting the view that Toulmin's argumentation model is the missing design foundation of software rule engines, and **demonstrating it through implementation**.

### 8.2 Limitations and Threats to Validity

**Model and implementation**

- Rule bodies (host functions) are opaque to the engine, so partial evaluation and static analysis are not supported. Graph structure and TANGEUL expression rules serialize to `.geul`, but host function bodies do not. Although the four ports implement the same computation, rule functions must be rewritten for each language.
- We have not analyzed whether the exception gate, attack force, and exclusion satisfy the principles studied in the gradual semantics literature. The [-1, 1] transform is linear, and complex rule interactions may require more sophisticated aggregation.
- Because cycles of attack edges and of supersede edges are rejected as errors, cases that need cyclic graphs, such as mutually attacking arguments, cannot be handled.
- A rule verdict of 0 means "half applies." Returning 0 for "no" out of bool habit silently changes the judgment, and the engine cannot detect it.
- Exclusion is a reporting mark. It does not change verdicts and takes effect only if downstream aggregation follows the mark. Results of `Evaluate` with trace off carry no marks.
- In v0.9.0 neither the public builder API nor TANGEUL offers a way to mark a node as strict. Strict handling is implemented inside the engine, but users can express a strict rule only by giving it no attack edges.
- Claim modeling requires further work: relational assertions between two subjects exceed the expressiveness of simple (subject, property, value) triples.

**Evaluation**

- The case studies are limited to the author's own projects. Validation on external open-source projects is future work.
- The figures in §6.3 show the effect of applying the rule set and do not isolate the effect of the Toulmin conversion itself. The measurement commands and environment are not included in this edition.
- §7 is a design case study without quantitative evaluation.
- The argument in §3.1 about "fact" naming is conceptual and has not been tested empirically.

### 8.3 Relationship to Existing Systems

This engine does not replace Rego, Drools, or Semgrep; it provides a **superordinate abstraction**. Rego can serve as the backend of a runtime Ground Adapter, and Semgrep rules can be reused in a compile-time Ground Adapter. Existing systems become Ground Adapter implementations for specific phases.

---

## 9. Conclusion

Since the advent of production rule systems, software rule engines have treated validation targets as "facts" and built their designs on that assumption. Toulmin's argumentation model (1958) already provided a suitable structure, namely Claim (not fact), Ground, Warrant, Backing, Qualifier, and Rebuttal, but it has rarely been used in rule engine design.

We implemented `toulmin`, a rule engine that bridges this gap. Rules are functions that take a context and specs and return a rule verdict and evidence. Spec is the judgment criterion fixed at declaration time (Backing), and the context holds the per-request facts about the judgment target (Ground). The engine organizes rules into an attack graph and computes verdicts in [-1, 1] with an extension of the weighted h-Categoriser. Rebuttals are computed in the denominator and exceptions as an applicability gate, attack edges can carry a force, and exclusion reports a crisp priority without touching verdicts. The same computation exists in four ports (Go, TypeScript, Python, C), and TANGEUL adds an auditable declarative surface and the `.geul` source format on top of the core. The trace records which rules took part in the judgment, in what role, and with what spec and evidence.

The core needs no separate DSL. The v0.9.0 Go core (`go/pkg/toulmin`) consists of 95 Go files excluding tests, totaling 1,078 lines excluding blank lines and `//` comment lines (measured on Linux at commit `4a6e568` with `find go/pkg/toulmin -name '*.go' ! -name '*_test.go' | xargs cat | grep -v '^\s*//' | grep -v '^\s*$' | wc -l`). The same rule function can be reused across graphs with different attack relations or with different specs, separating rule logic from graph structure. The conversion of filefunc's 22 rules shows that the Toulmin model naturally accommodates rule strength classification, structured defeasibility, and explicit spec, and the fullend case shows a design in which Ground Adapters let the same rule be reused across phases.

The `toulmin` software is available under the MIT License.

---

## References

[1] Toulmin, S. *The Uses of Argument*. Cambridge University Press, 1958.

[2] Nute, D. "Defeasible Logic." In D.M. Gabbay, C.J. Hogger, J.A. Robinson (eds.), *Handbook of Logic in Artificial Intelligence and Logic Programming*, Vol. 3, Oxford University Press, 1994.

[3] Amgoud, L., Ben-Naim, J., Doder, D., and Vesic, S. "Acceptability Semantics for Weighted Argumentation Frameworks." *IJCAI 2017*, pp. 56–62, 2017.

[4] Gabriel, V.O., Panisson, A.R., Bordini, R.H., Adamatti, D.F., Billa, C.Z. "Reasoning in BDI agents using Toulmin's argumentation model." *Theoretical Computer Science*, 805:76–91, 2020.

[5] Modgil, S., and Prakken, H. "The ASPIC+ Framework for Structured Argumentation: A Tutorial." *Argument and Computation*, 5(1):31–62, 2014.

[6] Dung, P.M. "On the Acceptability of Arguments and its Fundamental Role in Nonmonotonic Reasoning, Logic Programming and n-person Games." *Artificial Intelligence*, 77(2):321–357, 1995.

[7] Ahrendt, W., et al. "StaRVOOrS: A Tool for Combined Static and Runtime Verification of Java." *Runtime Verification*, LNCS 9333, 2015.

[8] Ahrendt, W., et al. "A Specification Language for Static and Runtime Verification of Data and Control Properties." *FM 2015*, LNCS 9109, 2015.

[9] Cockburn, A. "Hexagonal Architecture (Ports and Adapters)." 2005.

[10] Gödel, K. "Über formal unentscheidbare Sätze der Principia Mathematica und verwandter Systeme I." *Monatshefte für Mathematik und Physik*, 38(1):173–198, 1931.

[11] Kelsen, H. *Reine Rechtslehre (Pure Theory of Law)*. Franz Deuticke, Wien, 1934.

[12] Pollock, J.L. "Defeasible Reasoning." *Cognitive Science*, 11(4):481–518, 1987.

[13] Raz, J. *Practical Reason and Norms*. Hutchinson, London, 1975. (2nd ed. Oxford University Press, 1999.)

[14] Caminada, M. "On the Issue of Reinstatement in Argumentation." *JELIA 2006*, LNCS 4160, Springer, 2006.

[15] Besnard, P., and Hunter, A. "A logic-based theory of deductive arguments." *Artificial Intelligence*, 128(1–2):203–235, 2001.

[16] Amgoud, L., and Ben-Naim, J. "Ranking-based Semantics for Argumentation Frameworks." *SUM 2013*, LNCS 8078, 2013.

[17] Dunne, P.E., Hunter, A., McBurney, P., Parsons, S., and Wooldridge, M. "Weighted argument systems: Basic definitions, algorithms, and complexity results." *Artificial Intelligence*, 175(2):457–486, 2011.

[18] Antoniou, G., Billington, D., Governatori, G., and Maher, M.J. "Representation results for defeasible logic." *ACM Transactions on Computational Logic*, 2(2):255–287, 2001.

[19] Verheij, B. "Evaluating arguments based on Toulmin's scheme." *Argumentation*, 19(3):347–371, 2005.

[20] Jones, M., Bradley, J., and Sakimura, N. "JSON Web Token (JWT)." RFC 7519, IETF, 2015.

[21] GEUL grammar: TOULMIN codec. https://geul.org/grammar/toulmin/
