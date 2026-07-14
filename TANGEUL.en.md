# TANGEUL

**Toulmin Amgoud Nute Graph - General Encoding Unified Language**

An executable language written by AI and read by humans. A markdown file is the executable code itself.

Core principle: **repetition is only safe on top of idempotency, and failure is never silent — it is compensated or escalated to a human.** In this language's domain of use (autonomous RPA handling money, e.g. tax/banking), re-firing every tick becomes a double transfer, and a half-executed process is the most dangerous state. TANGEUL blocks both of these at the language level.

## File format

Markdown (`.md`). List items inside `## tangl:` sections are executable code; everything else is documentation.

- `## tangl:SectionName` — start of an executable code section
- `- statement` — unordered executable code. Order-independence and parallelizability apply only to **judgment statements** (node registration, defeat edges). **The firing order of execution edges (`do`/`run`/`undo`) is deterministic by document declaration order** (see §Execution Order Determinism)
- `1. statement` — ordered executable code. Executes sequentially in numeric order. A non-contiguous number is a parser error. **Execution sequences where order carries meaning are recommended to be expressed as an ordered list**
- `  - ` / `  1. ` — nested list
- `- ~~statement~~` — strikethrough is a comment. Ignored by the parser. An edge referencing a struck-through node without strikethrough is an error
- `## General Heading` / other text — documentation (ignored by the parser)

## Naming rules

All user-defined names **must be backtick-quoted**. Only keywords (`is`, `a`, `rule`, `don't`, `do`, `run`, `undo`, `once`, `when`, etc.) are used without backticks.

Double quotes and backticks have one uniform meaning in every syntactic position: a matching outer pair of double quotes denotes a string value and is syntax (the quotes are not part of the value), while backticks always denote a name, variable, or path component. This applies equally to definition values, comparison literals, and `is in` items. Unquoted literals retain their existing meaning.

```
`authenticate` is a general rule
don't `authenticate` when `block ip`
do `arm`.`placeCup` once when `order received`
undo `bank`.`refund` when `balance sufficient`
run `can pour water` when `order received`
```

Backticks visually separate variables from keywords. Backticked = user-defined name; unbacktricked = keyword or package alias.

## The seven sections

| Section | Role | Required | Allowed syntax |
|---|---|---|---|
| `tangl:Subject` | Subject declaration | required | `this document is` |
| `tangl:See` | Reference declaration | optional | `see ... from` |
| `tangl:Definitions` | Term definitions | optional | `means ... as`, `means`, `has ... as` |
| `tangl:Rules` | Inline rules | optional | `"..." when`, `and`, `or`, `not`, `either` |
| `tangl:Cases` | Judgment + execution | required | `in case of`, `is required`, `is a .../counter/except rule`, `checking`, `don't ... when`, `do ... when`, **`do ... once when`**, **`undo ... when`**, `run ... when` |
| `tangl:Provides` | Endpoints | optional | `provides`, `is required`, `run`, `check` |
| `tangl:Internal` | Internal work (ticks/events) | optional | `on`, `every`, `until`, `run`, `check` |

Order: Subject → See → Definitions → Rules → Cases → Provides → Internal. Each statement type belongs only to its designated section.

## Separation of judgment and execution — the Evaluate / Run duality

The same case graph runs in **two modes**:

| Mode | Meaning | `do`/`run` fires? | Used for |
|---|---|---|---|
| **Evaluate** | Pure judgment. Produces only a verdict | ❌ no | `checking` (judgment composition), read-only endpoints (`check`), tests |
| **Run** | Judge-first → execute-after | ✅ Active nodes' `do`/`run` fires | Executing endpoints (`run`), Internal ticks |

This duality is **why it is safe to put execution inside Cases**:

- If a case is **referenced as a rule** (`checking`) or its verdict is merely **queried**, that is Evaluate → **zero side effects**. Judgment is purely reused and tested.
- If an endpoint **`run`s** a case, that is Run → the Active nodes' execution fires.

> **Prohibition**: side effects must never leak from the read path (Evaluate). `checking` must never fire execution — otherwise you get the footgun of "I only referenced a judgment, and a robot moved."

## Subject (tangl:Subject)

```
- this document is `<namespace>`
```
The document's namespace, used for human-readable graph labels. External case references do **not** use the target document's own namespace; instead a referring document reaches another document's cases through a See-declared **local alias**, `` `alias`.`caseName` `` (the same local-alias mechanism as `using`, analogous to Go's `import x "path"`). See §Cases (`checking`) and §References.

## References (tangl:See)

```
- see `<alias>` from `<package_path>`
```
Rule functions and Spec factories of the referenced package are used as `` `alias`.`Name` ``.

The **same** `see` declaration is read two ways, decided by where the alias is used:
- In a **`using`** clause the alias is a **Go import path** — the referenced package's rule functions and Spec factories are called as `` `alias`.`Name` `` (the existing behavior, unchanged).
- In a **`checking`** clause the alias names **another TANGEUL document**: `<package_path>` is resolved as a **sibling `.tangeul` file path relative to the declaring document's directory** (a missing `.tangeul` extension is appended), *not* as a Go import path. Cross-document composition is a Go **cross-package call fixed at `gen` time** (not a runtime dynamic loader): `gen` emits `` `alias`.`Checking<Case>` `` and derives that package's Go import path from where the sibling document lives (its nearest enclosing `go.mod`). See §Cases (`checking` · CASE_REF).

## Term definitions (tangl:Definitions)

Constants (`means`), variables (`means ... as`), structs (`means` + `has`). Referenced in Cases as `` with `<term>` ``. Basic types (`Text`, `Integer`, `Number`, `Boolean`, `Email`, `Date`, `Time`, `URL`, `Currency`, `Quantity`) are provided by `tangl/types`.

```
- `credit threshold` means 650 as `credit`.`Threshold`
- `applicant` means
  - has `annual income` as Currency
  - has `credit score` as Integer
```

## Rules (tangl:Rules)

Inline condition definitions.

```
- `<name>` when <field> <operator> <value>
1. `<name>` when
    - either
      - <condition1>
      - and <condition2>
    - or <condition3>
```

Operators: `is empty`, `is not empty`, `equals`, `is in`, `is greater than`, `is less than`, `is at most`, `is at least`. Logic: `and`, `or`, `not`, `either` (grouping).

### Operands

Beyond a plain `` `field` ``, either side of a comparison may be an **operand**:

```
- `trusted author` when `author`.`trustScore` is at least 0.7
- `is internal` when `attributes`.`internal` equals true
- `within budget` when spec `discount`.`min price` is at most `basePrice`
- `high confidence` when result of `nlp`.`score` is at least 0.7
- `flagged` when `check`.`isSpam` holds
- `has coupon` when any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`
1. `active promotion` when any of `promotions` as `p` satisfies
    - `p`.`active` equals true
    - and `p`.`discount` is at least 10
```

- **Data path** — a backtick name chain joined by dots (`` `author`.`trustScore` ``) walks struct fields and map keys of the judgment context. A **bare backtick dot-chain is always a data path**; it is never a function call. Host calls must carry an explicit marker (`result of` or `holds`), so an operand can be classified locally without a symbol table.
- **Spec reference** — `` spec `<term>` `` (optionally path-extended, `` spec `discount`.`min price` ``) reads an argument injected by a node's `with` clause.
- **Host value call** — `` result of `<func>` `` calls a host function and uses its return value as the operand (either side of a comparison).
- **Host predicate call** — `` `<func>` holds `` is a condition term that is true when the host predicate returns true.
- **Existential quantifier** — `` any of `<source>` as `<binder>` satisfies <condition> `` holds when at least one element of the source collection satisfies the condition. The binder (`` as `c` ``) names the current element and is referenced as the first step of a data path inside the body (`` `c`.`minPrice` ``); it lexically shadows a context key of the same name (the binder wins). A binder must not reuse the name of an enclosing binder in a nested `any of` (load error). When the text after `satisfies` is empty, the body is the item's nested `and`/`or`/`either` condition list; otherwise it is a one-line comparison. An `any of` head may also stand alone as a rule head (`` `has coupon` when any of ... ``).
- **Right side** — a comparison right side that begins with a backtick, `spec`, or `result` is parsed as an operand; anything else (`650`, `"VIP"`, `a, b, c`) stays a literal value. One matching pair of surrounding double quotes is syntax and is removed (`"VIP"` becomes `VIP`); for `is in`, this normalization is applied to each comma-separated item (`"VIP", "GOLD"` becomes `VIP, GOLD`). Unquoted literals retain their existing meaning. Commas delimit items and are not supported inside an item. Consequently a **backtick right side denotes a data path**, not a literal — writing a name in backticks to mean a literal string is not supported.
- **Reserved quantifiers** — `all of ...`, `count of ...`, and `sum of ...` are reserved words; the parser rejects them with a clear "reserved" error until a later version implements them.
- **Two-stage name lookup** — a bare `` `name` `` reference resolves **document-local `tangl:Rules` first, then externally** (a Go symbol or a consumer registry). The code generator detects a document/host name collision as a Go duplicate-identifier compile error; a runtime loader follows the same document-local-first rule and may additionally escalate a collision to an error.
- **Scope of application** — the condition grammar is **predicate-only**: it produces a boolean judgment and has no execution/action statements. Effects live in `tangl:Cases` (`do`/`run`/`undo`), never inside a rule condition.
- **Determinism** — the expression grammar itself has no clock, randomness, or I/O. A clock-dependent predicate is kept pure by injecting the current time as a context key (`now`); a host call is the deliberate escape hatch whose determinism the host owns.

#### Execution binding

- **Execution binding** — the implementation exposes the condition AST and provides a **reference evaluator** (`pkg/tangl/eval`); predicate execution is **not built into the engine** (`pkg/toulmin`). A consumer loads a document, injects a host-function resolver, and compiles each `tangl:Rules` entry to a `func(toulmin.Context, toulmin.Specs) (bool, any)`. The code generator (`pkg/tangl/gen`) is the other path; both produce the same `(bool, evidence)` for the same document and context.

#### Evidence determinism (the witness leaf)

- **Evidence determinism** — evaluation is strict, document-order, left-to-right short-circuit (`and` stops at the first false; `or`/`either` at the first true; `not` inverts; `any of` stops at the first satisfying element). The **leaf where evaluation stopped** — the leaf whose truth value settled the whole result — is the **witness**, and its auto-generated `toulmin.Fact` is returned as the evidence. This makes the evidence deterministic: same (document, context) ⇒ same `(bool, Fact)` down to the byte.

The witness leaf's `toulmin.Fact` is populated per leaf kind:

| Leaf | Where | Expected | Actual |
|---|---|---|---|
| Compare (flat / path / spec) | left-side notation — the `Field`, or the path/spec steps joined by `.` (an `any of` element includes its index, e.g. `coupons[2].minPrice`) | `"<op> <right-side notation>"` (e.g. `is at least 0.7`) | the measured left value (nil when there is no value; when the left side is a host value call that returned an error, that error value) |
| Host predicate call (`holds`) | `alias.name` (or bare `name`) | `holds` | the boolean result — **except** when the host returns a non-nil evidence, which passes through unchanged (the auto-Fact never overwrites what the host already reported) |
| `any of` — no element satisfies (false) | the source notation | `any of <source> as <binder> satisfies` | the number of elements examined |
| `any of` — an element satisfies (true) | the satisfying element's inner witness, propagated upward | 〃 | 〃 |
| `not` | the inner witness, unchanged (the verdict conveys the negation; Expected is not transformed) | | |

#### Comparison & path semantics

- **No-value handling** — a missing ctx key / missing struct field / missing map key / nil intermediate yields "no value"; under the eight operators this behaves exactly as the pre-existing flat comparison did — `is empty` is true on no-value, every other operator is false on no-value.
- **Numeric comparison** — `is greater than` / `is less than` / `is at most` / `is at least` stringify both sides (`fmt.Sprintf("%v", …)`) and parse them as `float64`; if either fails to parse the comparison is false. `equals` / `is in` compare stringified forms.
- **Binder scope** — inside `` any of ... as `c` satisfies ... ``, a data path whose first step is `c` resolves to the bound element; the binder shadows a same-named ctx key.
- **`any of` source** — must be a slice or array (iterated in index order — deterministic); a map or non-collection source is treated as zero elements (Go map iteration order is nondeterministic and therefore unsupported).
- **Determinism (evaluator)** — extending the grammar-level **Determinism** bullet above: the evaluator carries no clock, randomness, or I/O, so a clock-dependent predicate is purified by injecting `now` as a context key, and host calls remain the deliberate escape hatch whose determinism the host owns.

#### Canonical rendering (mode 2)

Mode-2 (normalized) rendering re-emits a rule in its **minimal** shape, the contrapositive of the parser's one-line rule-head dispatch (`inline_term = exists_term | call_term | compare`):

- **Rule head** — a rule whose condition is a **lone comparison** (operand comparisons included), a **lone `any of`** with a lone-comparison body, or a **lone predicate call** renders on one line (`` `name` when <compare> `` / `` `name` when any of ... satisfies <compare> `` / `` `name` when `f` holds ``). An `and`/`or` chain, an `either`, a `not`, or an `any of` with a nested body renders as the `` `name` when `` head plus a nested list.
- **`any of`** — a lone-comparison body follows `satisfies` on the same line; any other body is the head line plus a nested list.
- **Predicate call (`holds`)** — always a one-line term, both as a condition term and as a rule head.

For example, a lone-body `any of` stays one line while a nested body expands:

```
- `has coupon` when any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`
- `active promotion` when any of `promotions` as `p` satisfies
  - `p`.`active` equals true
  - and `p`.`discount` is at least 10
```

## Cases (tangl:Cases) — judgment + execution

### Case declaration

```
- in case of `<name>`
  - `<field>` is required
```
Case names conventionally follow the `can ...` question form. `is required` declares required context fields (parser error if absent).

### Node registration

```
- `<name>` is a general rule
- `<name>` is a general rule using <func>
- `<name>` is a counter rule using <func>
- `<name>` is an except rule using <func>
- `<name>` is a general rule using <func> with <Spec> qualified <float>
- `<name>` is a counter rule using <func> grounded on <ground>
- `<name>` is a general rule checking <other-case>
- `<name>` is a general rule checking `<alias>`.`<other-case>`
- `<name>` is a counter rule using <func>, level fail
- `<name>` is an except rule using <func>, level review
```

| Node | Engine | Meaning |
|---|---|---|
| `general rule` | Rule (Warrant) | claim |
| `counter rule` | Counter (Rebuttal) | rebuttal (reverses verdict) |
| `except rule` | Except (Defeater) | exception (attenuates verdict) |

`, level fail` / `, level review` is an optional trailing modifier on a node
(after `using`/`checking`) that pins the node's **crisp reporting level** for
reins-style aggregation (any Fail → FAIL, else any Review → REVIEW, else PASS).
It compiles to `.Level(toulmin.Fail)` / `.Level(toulmin.Review)`. When omitted,
nothing is emitted and the engine's role-based default applies — so documents
without the modifier are byte-for-byte unchanged.

`grounded on <ground> [and <ground>]*` is an optional trailing modifier (after
`using ... [with ...][qualified ...]`) that records the **external ground-fact
names** this rule's function resolves at evaluation time — e.g. a network lookup
behind a predicate. It compiles to `.Grounds("<ground>", ...)`. It is a
**different axis** from `with`: `with` fixes a judgment-criteria spec at
declaration time, whereas `grounded on` names an external fact resolved at
evaluation time.

| Clause | Nature | Timing |
|---|---|---|
| `with <Spec>` | judgment-criteria spec | fixed at declaration time |
| `grounded on <ground>` | external fact resolved at evaluation time | evaluation time — a different axis |

`checking <case>` is **judgment composition** — using another case's verdict as a rule. It **compiles to Evaluate** and fires no side effects. Circular references are forbidden. The target is a `CASE_REF` (the same grammar as `FUNC_REF`): a bare local `` `case` `` names a case in the same document, and an alias-qualified `` `alias`.`case` `` names a case reached through a See-declared local alias (cross-document). A cross-document reference **compiles to a Go cross-package call fixed at `gen` time**: `gen` emits `` `alias`.`Checking<Case>` `` — a call into the exported wrapper the referenced document's own `gen` output publishes for that case — and adds an import of the referenced document's generated package (its import path derived from the sibling `.tangeul`'s location; see §References). Every case is externally addressable, so a document exports a `Checking<Case>` wrapper for **each** of its cases, whether or not it checks them itself. Cross-document cycles (A's case checks B's case checks A's case) are rejected by the multi-document loader before `gen` runs.

> **Implementation note (Active composition rule for case verdicts)**: a case may have multiple warrants (`general rule`). The "case's judgment" that `checking` references is composed by the following rule — among a case's Evaluate results, **if one or more warrants have verdict > 0, the case is considered established (active)**. `until` (§Internal) uses the same rule.

### Defeat edges (don't / do not)

```
- don't `<target>` when `<attacker>`
```
If the attacker is Active, it attacks the target. `do not` is equivalent. → `attacker.Attacks(target)`

### Supersede edges (supersedes)

```
- `<superseder>` supersedes `<superseded>`
```
If the superseder is Active, the superseded node is **crisply excluded from the
tally** (it never enters the h-Categoriser traversal — only its trace is marked
Excluded). → `superseder.Supersedes(superseded)`

| Edge | Effect | Engine |
|---|---|---|
| `don't <y> when <x>` | attenuates `<y>`'s verdict (graded, enters the tally) | `x.Attacks(y)` |
| `<x> supersedes <y>` | crisply excludes `<y>` from the tally (deterministic priority) | `x.Supersedes(y)` |

Like defeat edges, supersede edges must reference nodes registered in the same
case and must not form a cycle (validated independently of the Attacks graph).

### Execution edges (do / run)

Execution is attached with **node Active as the trigger**. It is the execution counterpart of the defeat edge (`don't ... when`).

```
- do `<func>` when `<node>`
- do `<func>` once when `<node>`
- do `<func>` once when `<node>` if at least <N>% certain
- run `<case>` when `<node>`
```

| Syntax | Engine | Meaning |
|---|---|---|
| `do <func> when <node>` | `node.RunOn(...)` | If node is Active, executes a **leaf action** (func) |
| `do <func> once when <node>` | `node.RunOn(...)` + ctx guard | Same as above, but **at most once per ctx lifetime** (§once) |
| `run <case> when <node>` | `node.Run(<case graph>)` | If node is Active, **Runs a sub-case** (execution composition, cascade + Trace) |

- **`do` = leaf** (calls another subject's func/`provides`). **`run` = composite** (Runs a sub-case graph, whose `do`/`run` fire in turn).
- `<node>` is a node name within the same case. Fires only when the node is Active (`Activated && verdict > 0`).
- Certainty gate: `if at least/above/less than/at most <N>% certain` → the handler branches on `self.Verdict`.

> **Implementation note (function signature)**: functions referenced by `do`/`undo` follow the `func(ctx) error` signature (Go implementation: `func(toulmin.Context) error`). They take ctx as an argument and return an error — the returned error is the signal that triggers Run abortion and compensation (undo).

> **Acyclicity constraint**: Cases must always be a **acyclic DAG** (enforced by the engine's `detectCycle`/`detectRunCycle`). Case→case cascading via `run` is allowed, but **cycles are forbidden**. **Repetition does not belong in Cases** — it is expressed via Internal's tick loop.

> **Implementation note (current engine constraint)**: a node has one RunOn handler and one RunGraph. If a `<node>` has multiple `do`s, they compose into **sequential execution in document order** within a single handler. If a `<node>` has multiple `run`s, a composite subgraph is needed (engine extension or wrapper). If fan-out is needed, split the trigger node.

### The once modifier — tick idempotency

```
- do `<func>` once when `<node>`
- do `<func>` once when `<node>` if at least 75% certain
```

**Fires at most once per ctx lifetime.** Since the Internal tick loop shares ctx across ticks, `once` means "once per loop lifetime (until `until` terminates)." For one-shot endpoints (Provides `run`), ctx is created fresh per request, so it is once per request.

- **When it's consumed**: the `once` guard is consumed **only after the firing succeeds**. If the handler errors or is blocked by the certainty gate, it is not consumed — it will be retried on the next tick.
- **Composable with the certainty clause**: `once` comes first, the certainty clause after. Consumption occurs only after passing the gate and firing/succeeding.
- **Implementation (no engine changes needed)**: the codegen inserts a ctx guard inside the generated RunOn handler. The key convention is `once:<subject>.<case>.<node>#<n>` (the n-th `do` of the same node, 0-based). `ctx.Set` after a successful fire. Since ctx flows into cascaded subgraphs, the guard is consistent across the entire cascade.

Why `once` is needed: `every 30s` **re-Runs** a case wholesale every tick. `do arm.placeCup when order received` would **place a cup every tick** while the order stays active. In a coffee robot, sensors naturally update the world so the next stage's guard catches it, but in transfers or filings there is no such sensor — every tick becomes a duplicate transfer or duplicate filing. `once` blocks this accident with a single grammar word.

### The undo edge — partial failure and compensation

```
- undo `<func>` when `<node>`
```

After the node's `do` **fires successfully**, if a handler error occurs at a later point within the **same top-level Run pass** (including the whole cascade), the undos of already-fired nodes are executed in **reverse firing order (LIFO)**. This is the compensation pattern of sagas.

| Rule | Content |
|---|---|
| Arming point | Compensation is armed immediately after the node's `do` completes successfully. `undo` on a node with no `do` is a parser error |
| Trigger condition | A handler error occurring at a **later point** in the same top-level Run pass (the engine aborts the Run immediately on handler error — the wrapper then runs compensation afterward) |
| Execution order | Reverse arming order (LIFO). If a node has multiple `undo`s, they are armed in document order → executed in reverse |
| Scope | **Once per top-level Run pass.** The compensation stack is initialized at pass start, discarded (committed) on pass success, and executed in reverse on failure. In Internal, each tick is a new pass — note that its scope differs from the `once` guard (ctx lifetime) |
| Pass owner | **Two**: the Provides `run` wrapper (per request) and Internal's generated tick runner (per tick). The tick runner also initializes the compensation stack at the start of each tick → commits on tick success → LIFO compensation on tick failure → escalates to REVIEW on compensation failure, following the same contract as the wrapper |
| Compensation failure | If the undo itself errors, **the remaining compensations are aborted and the whole thing is escalated to REVIEW (a human gate)** (MUST). Compensation failure is never silently papered over. Logging is mandatory, silence is forbidden |

- **Implementation (no engine changes needed)**: at the codegen level. Each `do` handler pushes a compensation closure onto **the compensation stack inside ctx** after success, and the Provides wrapper (`run` endpoint) executes the stack in reverse when `g.Run(ctx)` errors. Since ctx flows into subgraphs, **the entire cascade shares one stack** — a parent's withdrawal and a child's deposit form a single compensation unit.

> **Distinction**: `do <func> when <counter-node>` and `undo` are different things. The former is an action for **when a judgment is rejected** (e.g. a deny node becomes Active → `Deny`), while the latter is compensation for **when execution dies partway through**. Judgment rejection is one branch of the normal path; execution failure is an abnormal termination.

### Execution order determinism

The rule that "`- ` unordered = order-independent, runtime may parallelize" holds only for **judgment statements** (node registration, `don't` edges). Since judgment is pure, the verdict is the same regardless of evaluation order.

Execution is different:

- **Execution firing order is deterministic by the node's document declaration (registration) order.** If two nodes are simultaneously Active, their `do`s fire in document order (engine: handlers fire in rule registration order).
- **Parallelizing execution is forbidden.** If the order of side effects becomes nondeterministic, claims of determinism and auditability break — for reproducibility and auditability to hold, the same document and the same world must always produce the same execution sequence.
- Multiple `do`s on the same node compose into sequential execution within a single handler, in document order.
- **Execution sequences where order carries meaning are recommended to be made explicit via `1.` ordered lists** — this reveals to the reader that meaning is carried by declaration order.

## Cascade temporal semantics — judge-first-execute-after is per level

"Judge-first, execute-after" is true **per level (graph)**, not globally. The actual temporal order of a `run <case> when <node>` cascade:

```
parent judge → parent execute → child judge → child execute → (grandchild judge → grandchild execute → …)
```

The child's judgment **sees the ctx/world as mutated by the parent's handler.** The coffee example depends exactly on this interleaving — the timeline within one tick (one top-level Run pass):

```
T1  can place cup judged     → order received Active (no cup yet)
T2  placeCup executes        → world changes: cup is placed
T3  can pour water judged    → cupIsPlaced sensor sees T2's result → Active
T4  pourWater executes
T5  can brew espresso judged → waterIsPoured is true
T6  extractShot executes
T7  can serve coffee judged
T8  serve executes
```

T3's judgment comes **after** T2's execution. Auditors must not read this as "all judgments precede all executions" — within each graph, all of that graph's judgments precede all of that graph's executions, and cascades interleave judgment and execution level by level. This is not a bug but a design: children must be able to judge based on the results of the parent's execution for a staged process to be expressed in a single pass. (If a sensor update is slow and T3 is still false, that tick simply stops there, and the next tick picks up where it left off — which is why the Internal tick loop and `once` are a pair.)

## Work (tangl:Provides) — endpoints

The **entry points** this Subject exposes externally. Judgment and execution live in Cases; Provides declares only **which case is triggered through which interface**. The document's location determines the mode of provision:

| Location | Provision mode | Go translation |
|---|---|---|
| `commands/` | CLI command | Cobra subcommand |
| `services/` | HTTP endpoint | `http.HandleFunc` |
| `documents/` | package public func | exported func |

### Endpoint declaration

```
- provides `<name>`
  - `<field>` is required
  - `<field>` is required as <Type>
  - run `<case>`
  - check `<case>`
```

| Syntax | Mode | Meaning |
|---|---|---|
| `` provides `<name>` `` | — | Endpoint declaration. Called externally as `` subject.`name` `` |
| `` `<field>` is required `` | — | Required input. Rejected if absent (`tangl.ErrRequired`) |
| `` run `<case>` `` | **Run** | Judges then executes the case. Cases' `do`/`run` fire |
| `` check `<case>` `` | **Evaluate** | Judges the case only. Returns a verdict, **no side effects** (for queries/inquiries) |

`run` and `check` surface the Evaluate/Run duality on the endpoint face — an auditor can immediately read "does this endpoint execute, or only query?"

### Go code generation — compensation wrapper

The `run` endpoint wrapper is the owner of the compensation stack: it initializes at pass start, runs the stack in reverse on `g.Run(ctx)` error, and discards it on success. Internal's generated tick runner performs the same cycle per tick (Init → Run → Compensate/Commit → Review on compensation failure) — even in tick-driven processes that never pass through an endpoint, there is still an owner for triggering `undo`.

```go
// provides `transfer` in documents/transfer/transfer.md (Subject: "transfer")
func Transfer(ctx tangl.Context) error {
    tangl.InitCompensation(ctx)              // initialize compensation stack (pass-scoped)
    _, trace, err := canWithdrawGraph.Run(ctx) // judge-first, execute-after (including cascade)
    if err != nil {
        if cerr := tangl.Compensate(ctx); cerr != nil { // LIFO reverse compensation
            return tangl.Review(ctx, err, cerr)          // compensation failure → escalate to REVIEW
        }
        return err                            // compensation completed — report original error
    }
    tangl.CommitCompensation(ctx)             // success — discard compensation
    return nil
}

// services/ → wrapped as an HTTP handler, commands/ → wrapped as a Cobra subcommand
```

## Internal work (tangl:Internal) — ticks / events (loop layer)

Work that self-triggers without an external request. **Repetition and infinity live here** (since Cases must be acyclic). Translated into Go unexported funcs.

```
- on <event>
  - run `<case>`
- every <interval>
  - run `<case>`
- every <interval> until `<case>`
  - run `<case>`
```

| Syntax | Meaning |
|---|---|
| `on <event>` | Runs once upon receiving the event |
| `every <interval>` | Repeats Run at the interval (`30s`, `1hour`, `7:00`) |
| `` every <interval> until `<case>` `` | Runs every tick until the termination condition is satisfied. **ctx is shared across ticks** — the lifetime of the `once` guard is exactly this loop. The termination condition uses the same Active composition rule as §Node Registration — terminates if one or more of `<case>`'s warrants have verdict > 0 |

Event sources: referenced packages (`` `sensor`.`buttonPress` ``), completion of other work (`` `access`.complete ``), the system (`start`).

> **Implementation note (`every`'s time-of-day notation)**: a clock-of-day schedule such as `7:00` is valid at the language level, but the current reference implementation (codegen) supports **only duration intervals** (extended formats such as `30s`, `1hour` parsed via `time.ParseDuration`). A time-of-day schedule parses syntactically but is treated as unsupported at code-generation time — a cron-style scheduler requires a separate implementation extension.

**Infinite state machine**: every tick is one Run pass of a finite, acyclic case (bounded latency, auditable via Trace). Infinity lives in the `every`/`on`/`until` loop, state lives in ctx, and judgment lives in a finite DAG — this separation simultaneously preserves real-time behavior, safety certifiability, and auditability.

### Tick idempotency requirement (MUST)

**Every `do` of a case reachable from Internal** (either a direct `run` target or belonging to its cascade closure) **must be one of the following:**

- **(a) Idempotent** — the world is the same no matter how many times it executes (e.g., overwriting state, log-like recording)
- **(b) `once`** — guarded to once per ctx lifetime
- **(c) Protected by a counter/sensor guard reflecting completion state** — structured so that the trigger node is no longer Active once execution has finished in the resulting world

Rationale: the tick loop re-Runs a case wholesale every tick. A node whose condition is still true will have its `do` fire again every tick. In a coffee robot, (c) naturally holds — the trigger for `pourWater` is `cup placed`, and pouring water causes the next stage's sensor to update the world. But in domains where **the world does not reflect completion of execution**, such as transfers, filings, or dispatches, neither (a) nor (c) applies — you must use (b) `once` or design a completion counter.

> **Lint recommendation**: an unguarded `do` (with no idempotency marker, no `once`, and no completion guard) on a case reachable from Internal is subject to a parser/linter **warning**.

### verdict → percent conversion

`% = (verdict + 1) × 50` — `+1.0` = 100%, `0.0` = 50%, `-1.0` = 0%. Used in the certainty gate (`if at least N% certain`).

## Example: robotic action control

```markdown
# Americano preparation

When a customer orders, place a cup, pour water, extract, and serve. If a sensor update
is slow, the next tick picks it up — the arm motions are protected from tick re-execution by once.

## tangl:Subject
- this document is `americano`

## tangl:See
- see `arm` from `robot/arm`
- see `sensor` from `robot/sensor`
- see `log` from `tangl/log`

## tangl:Definitions
- `serve temperature` means 65°C as `sensor`.`Temperature`

## tangl:Cases
- in case of `can place cup`
  - `order received` is a general rule using `sensor`.`isOrdered`
  - `no cup` is a counter rule using `sensor`.`noCupAvailable`
  - don't `order received` when `no cup`
  - do `arm`.`placeCup` once when `order received`
  - run `can pour water` when `order received`

- in case of `can pour water`
  - `cup placed` is a general rule using `sensor`.`cupIsPlaced`
  - do `arm`.`pourWater` once when `cup placed`
  - run `can brew espresso` when `cup placed`

- in case of `can brew espresso`
  - `water poured` is a general rule using `sensor`.`waterIsPoured`
  - do `arm`.`extractShot` once when `water poured`
  - run `can serve coffee` when `water poured`

- in case of `can serve coffee`
  - `brewing complete` is a general rule using `sensor`.`isComplete`
  - `low temperature` is an except rule using `sensor`.`temperatureLow` with `serve temperature`
  - don't `brewing complete` when `low temperature`
  - do `arm`.`serve` once when `brewing complete`
  - do `log`.`CoffeeStatus` when `brewing complete`

## tangl:Provides
- provides `make coffee`
  - run `can place cup`

## tangl:Internal
- every 1s until `can serve coffee`
  - run `can place cup`
```

Read by an auditor: "Once the order is received — place a cup **once**, then move on to pouring water. If there's no cup, the reception is void. … If the temperature is low, don't accept brewing as complete. Retry every second until serving is done." `order received` is Active every tick while the order stays active, so without `once` a cup would be placed every tick. `log`.`CoffeeStatus` falls under idempotency requirement (a) — since it's log-like, it is left unguarded.

## Example: access control

```markdown
## tangl:Subject
- this document is `api`

## tangl:See
- see `policy` from `bank/policy`
- see `log` from `tangl/log`

## tangl:Definitions
- `blocklist` means "blocklist" as `policy`.`IPList`

## tangl:Cases
- in case of `can access`
  - `user` is required
  - `authenticate` is a general rule
  - `block ip` is a counter rule using `policy`.`IsIPBlocked` with `blocklist`
  - `exempt internal ip` is an except rule using `policy`.`IsInternalIP`
  - don't `authenticate` when `block ip`
  - don't `block ip` when `exempt internal ip`
  - do `policy`.`Allow` when `authenticate`
  - do `log`.`AccessResult` when `authenticate`
  - do `policy`.`Deny` when `block ip`

## tangl:Provides
- provides `access`
  - `user` is required
  - run `can access`
```

`` do `policy`.`Deny` when `block ip` `` is **an action taken on judgment rejection** (attached to a counter node) — it is not an `undo`. Since ctx is freshly created per request and Allow/Deny are idempotent, `once` is also unnecessary here.

## Example: bank transfer (compensation)

```markdown
# Bank transfer

Withdraw, then deposit. If the deposit fails, refund the withdrawal.
No matter at which point failure occurs, the total across accounts is preserved, or escalated to human review.

## tangl:Subject
- this document is `transfer`

## tangl:See
- see `bank` from `bank/core`
- see `log` from `tangl/log`

## tangl:Cases
- in case of `can withdraw`
  - `from account` is required
  - `amount` is required
  - `balance sufficient` is a general rule using `bank`.`hasBalance`
  - `account frozen` is a counter rule using `bank`.`isFrozen`
  - don't `balance sufficient` when `account frozen`
  - do `bank`.`withdraw` once when `balance sufficient`
  - undo `bank`.`refund` when `balance sufficient`
  - run `can deposit` when `balance sufficient`

- in case of `can deposit`
  - `to account` is required
  - `recipient valid` is a general rule using `bank`.`isValidAccount`
  - do `bank`.`deposit` once when `recipient valid`
  - do `log`.`TransferComplete` when `recipient valid`

## tangl:Provides
- provides `transfer`
  - `from account` is required
  - `to account` is required
  - `amount` is required
  - run `can withdraw`
```

Read by an auditor: "If the balance is sufficient — execute the withdrawal **once**, and if it fails, roll it back via refund, then move on to depositing. If the account is frozen, the balance judgment is void. The deposit executes **once** if the recipient account is valid, and records completion."

Temporal order of the failure scenario:

```
1. can withdraw judged        → balance sufficient Active
2. withdraw executes, succeeds → refund compensation armed (stack: [refund])
3. can deposit judged          → recipient valid Active
4. deposit executes, errors    → Run aborts immediately
5. wrapper runs compensation stack in reverse → refund fires (the money comes back)
6. original error reported — the transfer failed, and the world is as it was before the transfer
```

If even `refund` fails (an error at step 5): the remaining compensations are aborted and the whole thing is escalated to **REVIEW** — a half-executed transfer is not a state that can be automatically papered over. A human looks at the ledger and decides. Both errors (the original error and the compensation error) are logged.

## Audit tooling — effect summary (MUST)

Because of `run` cascades, "what can this endpoint execute" cannot be answered from a single document — you must follow `run` to read the **transitive closure**. What's more dangerous is silent expansion: if a `do` is added to a case that was being referenced via `checking`, the checker side remains pure (as intended by design), but **the execution surface of every upstream process that `run`s that case silently widens.**

**Implementer requirement (MUST)**: `tangl effects <subject>.<endpoint>` — prints the **static transitive closure** (effect summary) of every `do`/`undo` reachable from the endpoint.

- The closure starts at the endpoint's `run <case>` and follows Cases' `run <case> when <node>` edges.
- **`checking` references are not included in the closure** — they compile to Evaluate and are pure.
- The effect summary of a `check <case>` endpoint is always the **empty set** — otherwise it is an implementation bug.

```
$ tangl effects transfer.`transfer`
do   bank.withdraw   once   (can withdraw / balance sufficient)
undo bank.refund            (can withdraw / balance sufficient)
do   bank.deposit    once   (can deposit / recipient valid)
do   log.TransferComplete   (can deposit / recipient valid)
```

> **Recommendation**: automatically insert a per-endpoint effect table when rendering documents. When an auditor opens the document, "everything in the world this entry point can touch" should be visible in a single table. Using the diff of the effect summary as a review gate in CI is also recommended — so that expansion of the execution surface cannot pass through silently.

> **Implementation note (output alignment)**: the output example above is hand-aligned. The normative element is the **content and order of the 4 columns** (kind / func / once / (case / node)), not character spacing — an implementation satisfies the spec by printing the same 4 columns with column-wise tab alignment (e.g., `text/tabwriter`).

## Parsing grammar

```
file = section*
section = subject_section | see_section | def_section | rules_section
        | cases_section | provides_section | internal_section | doc_section

rules_stmt = NAME "when" (inline_term | NEWLINE INDENT cond_list)
inline_term = exists_term | call_term | compare      # rule-head one-line dispatch (either/not are nested-only)
cond_list  = cond_term (NEWLINE ("and" | "or") cond_term)*
cond_term  = exists_term | either_term | not_term | call_term | compare
either_term = "either" NEWLINE INDENT cond_list
not_term   = "not" (compare | NEWLINE INDENT cond_list)
exists_term = "any" "of" operand "as" NAME "satisfies" (compare | NEWLINE INDENT cond_list)
call_term  = FUNC_REF "holds"
compare    = operand COMPARE_OP [literal | operand]
operand    = path | "spec" path | "result" "of" FUNC_REF
path       = NAME ("." NAME)*
COMPARE_OP = "is" "empty" | "is" "not" "empty" | "equals" | "is" "in"
           | "is" "greater" "than" | "is" "less" "than" | "is" "at" "most" | "is" "at" "least"

cases_stmt = case_decl | required_decl | node_decl | edge_decl | exec_decl
case_decl  = "in" "case" "of" NAME
node_decl  = NAME "is" ("a"|"an") ROLE_EXPR (["using" FUNC_REF [with_clause] [qual_clause] [grounded_clause] | "checking" CASE_REF] [level_clause])
level_clause = "," "level" ("fail" | "review")             # crisp reporting level
edge_decl  = DONT NAME "when" NAME                          # defeat edge
           | NAME "supersedes" NAME                         # supersede edge
exec_decl  = do_edge | run_edge | undo_edge
do_edge    = "do" FUNC_REF ["once"] "when" NAME [certainty_clause]
undo_edge  = "undo" FUNC_REF "when" NAME                   # compensation edge
run_edge   = "run" NAME "when" NAME                        # run <case> when <node>
certainty_clause = "if" THRESHOLD_OP INT "%" "certain"
with_clause = "with" with_arg [("and") with_arg]*
grounded_clause = "grounded" "on" with_arg [("and") with_arg]*
qual_clause = "qualified" FLOAT
ROLE_EXPR  = "general" "rule" | "counter" "rule" | "except" "rule"
DONT       = "don't" | "do" "not"

provides_stmt = provides_decl
provides_decl = "provides" NAME NEWLINE INDENT provides_body+
provides_body = required_decl | endpoint_decl
required_decl = NAME "is" "required" ["as" TYPE_REF]
endpoint_decl = ("run" | "check") NAME                     # run=Run, check=Evaluate

internal_stmt = on_decl | every_decl
on_decl    = "on" EVENT NEWLINE INDENT internal_body+
every_decl = "every" INTERVAL ["until" NAME] NEWLINE INDENT internal_body+
internal_body = ("run" | "check") NAME

THRESHOLD_OP = "at" "least" | "above" | "less" "than" | "at" "most"
NAME = BACKTICK_STRING
FUNC_REF = NAME | NAME "." NAME
CASE_REF = NAME | NAME "." NAME        # same grammar as FUNC_REF: bare local case, or `alias`.`case`
```

- The NAME node of an `undo_edge` must be preceded by at least one `do_edge` within the same case (parser-validated).
- Keywords: `this`, `document`, `see`, `means`, `as`, `is`, `a`, `an`, `from`, `general`, `rule`, `counter`, `except`, `case`, `using`, `checking`, `with`, `qualified`, `provides`, `required`, `run`, `check`, `do`, `don't`, `not`, `when`, **`once`**, **`undo`**, **`supersedes`**, **`level`**, **`fail`**, **`review`**, **`grounded`**, `on`, `every`, `until`, `has`, `if`, `at`, `least`, `above`, `less`, `than`, `most`, `certain`, `and`, `or`, `either`, `empty`, `equals`, `in`, `greater`, **`any`**, **`of`**, **`satisfies`**, **`holds`**, **`result`**, **`spec`**. Reserved (parsed but not yet implemented): **`all`**, **`count`**, **`sum`**.

## Go API correspondence table

| TANGEUL | Go |
|---|---|
| `` this document is `x` `` | `package x` |
| `` see `p` from `path` `` | `import p "path"` |
| `` in case of `x` `` | `x := toulmin.NewGraph("subject.x")` |
| `` `x` is a general rule using fn `` | `x := g.Rule(fn)` |
| `` `x` is a counter rule using fn `` | `x := g.Counter(fn)` |
| `` `x` is an except rule using fn `` | `x := g.Except(fn)` |
| `` `x` is a general rule checking `y` `` | `x := g.Rule(checkCase(y))` — **`y` is Evaluate (pure)** |
| `` `x` is a general rule checking `alias`.`y` `` | `x := g.Rule(alias.CheckingY)` — cross-package call; `alias` resolves to a Go import path derived from the sibling `.tangeul`'s location, and `CheckingY` is the exported wrapper the referenced document emits for case `y` (Evaluate, pure) |
| `` `n` when <expr> `` (a `tangl:Rules` entry) | `func n(toulmin.Context, toulmin.Specs) (bool, any)` — the second return is an auto-generated witness `toulmin.Fact` |
| `` result of `f` `` (host value call) | `func f(toulmin.Context, toulmin.Specs) (any, error)` |
| `` `f` holds `` (host predicate call) | `func f(toulmin.Context, toulmin.Specs) (bool, any)` |
| `` with `term` `` / `qualified 0.8` | `.With(defs["term"])` / `.Qualifier(0.8)` |
| `` grounded on `source-body` `` | `.Grounds("source-body")` |
| `` don't `x` when `y` `` | `y.Attacks(x)` |
| `` `x` supersedes `y` `` | `x.Supersedes(y)` |
| `` , level fail `` / `` , level review `` | `.Level(toulmin.Fail)` / `.Level(toulmin.Review)` |
| `` do `fn` when `n` `` | `n.RunOn(func(self, t){ ...fn... })` |
| `` do `fn` when `n` if at least 75% certain `` | `n.RunOn(func(self,t){ if self.Verdict >= 0.5 { fn } })` |
| **`` do `fn` once when `n` ``** | **codegen ctx-guard handler — `n.RunOn(func(self,t){ if _,ok := t.Ctx().Get("once:subj.case.n#0"); !ok { if err := fn(); err != nil { return err }; t.Ctx().Set("once:subj.case.n#0", true) } })`** |
| **`` undo `fn` when `n` ``** | **`n`'s do handler pushes a compensation closure (`fn`) onto the ctx compensation stack after success; the Provides `run` wrapper runs the stack in LIFO reverse when `g.Run(ctx)` errors, escalating to REVIEW on undo error** |
| `` run `c` when `n` `` | `n.Run(cGraph)` (RunGraph execution composition, ctx shared downward, verdict isolated) |
| **`` provides `x` `` + `` run `c` ``** | **`func X(ctx){ InitCompensation(ctx); _, tr, err := cGraph.Run(ctx); if err != nil { Compensate(ctx) ... } }`** |
| `` provides `x` `` + `` check `c` `` | `func X(ctx){ r, err := cGraph.Evaluate(ctx); ... }` (no side effects) |
| `` `f` is required `` | `ctx.Get("f")` + nil check, `tangl.ErrRequired` on failure |
| `on event` | event handler registration |
| `every 30s` | `time.NewTicker` + goroutine, `cGraph.Run(ctx)` every tick — **ctx shared across ticks (the once guard's lifetime), per-tick Init/Compensate/Commit/Review compensation cycle** |
| `` every 30s until `c` `` | ticker + Evaluate `c` every tick, terminates on pass |

## Design principles

1. **Markdown is the code itself** — documentation and policy in one file. GitHub renders it, the parser executes it.
2. **Strict section structure** — each statement type belongs only to its designated section.
3. **Separation of judgment and execution via call mode** — Cases hold both the judgment graph and node-attached execution together, gating side effects via the Evaluate (pure)/Run (execution) duality. `checking` is pure; only endpoint `run` arms execution.
4. **Three-way split of declaration, connection, and execution** — node registration (`is a ... rule`), defeat edges (`don't ... when`), and execution edges (`do`/`undo`/`run ... when`) are separated. Execution attaches with node Active as the trigger (1:1 with engine RunOn/RunGraph).
5. **Endpoints stay thin** — Provides is only the provision surface (CLI/HTTP/func) + required inputs + `run`/`check`. Policy lives in Cases, interface lives in Provides. Ownership of the compensation stack also belongs to the wrapper (endpoint).
6. **Infinity lives at the outer tick** — Cases is an acyclic finite DAG. Repetition and infinite loops live in Internal (`every`/`on`/`until`). Every tick is one bounded, traceable Run pass.
7. **Continuous-valued verdict** — counter rejects, except warns. Graded execution via the certainty gate (`if N% certain`).
8. **No concealment** — all policy is open to the auditor.
9. **1:1 Go API correspondence** — `run→Run`, `do→RunOn`, `run case→r.Run(subgraph)`, `checking→Evaluate`, `once→ctx guard`, `undo→compensation stack`. No ambiguity.
10. **Written by AI, read by humans** — readability is maximized.
11. **Repetition is only safe on top of idempotency** — every `do` reachable from Internal is either idempotent, `once`, or protected by a completion guard. The tick loop runs only on top of this contract.
12. **Failure is never silent** — partial execution is compensated via `undo` or escalated to a human via REVIEW. Compensation failure is never automatically papered over. Logging always remains.
13. **The execution surface must be statically enumerable** — the complete set of `do`/`undo` reachable from an endpoint (effect summary) can be listed without executing anything. Execution order is also deterministic by document order — same document, same world, same sequence.
