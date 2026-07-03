# TANGL v0.3

**Toulmin Amgoud Nute Graph Language**

An executable language written by AI and read by humans. A markdown file is the executable code itself.

> Changes vs v0.2: **fills in the semantics of "repeated Run" and "failing Run."** The `once` modifier adds idempotency for tick re-execution to the language, and the `undo` edge adds compensation (saga) for partial failure. It codifies the temporal semantics of cascades and the firing order of execution, and locks in static enumeration of the execution surface (effect summary) as an implementer requirement. v0.2's "single Run" semantics (Evaluate/Run duality, node-attached execution) are carried forward unchanged.

## v0.3 change summary

| Area | v0.2 | v0.3 |
|---|---|---|
| Tick re-execution (repeated Run) | No semantics — `do` re-fires every tick the condition is true | **`do ... once when`** (at most once per ctx lifetime) + **idempotency requirement codified** |
| Partial failure (failing Run) | Run aborts immediately on handler error, leaving half-executed state | **`undo ... when`** compensation edge — reverse-firing-order (LIFO) compensation, escalates to REVIEW on compensation failure |
| Cascade temporal order | Implicit (described only as "judge-first, execute-after") | Codified that judge-first-execute-after applies **per level** — parent judge → parent execute → child judge → child execute |
| Execution firing order | `- ` unordered = order-independent, parallelizable (across the board) | Order-independence applies **only to judgment statements** — execution firing is **deterministic by document declaration order**, parallelization forbidden |
| Audit tooling | Trace (dynamic, post-execution) | **`tangl effects`** — static transitive closure of `do`/`undo` per endpoint (MUST) |

Core principle: **repetition is only safe on top of idempotency, and failure is never silent — it is compensated or escalated to a human.** In this language's domain of use (autonomous RPA handling money, e.g. tax/banking), re-firing every tick becomes a double transfer, and a half-executed process is the most dangerous state. v0.3 blocks both of these at the language level.

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

## Separation of judgment and execution — the Evaluate / Run duality (carried over from v0.2)

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
The document's namespace. Case names are referenced externally as `` namespace.`caseName` ``.

## References (tangl:See)

```
- see `<alias>` from `<package_path>`
```
Rule functions and Spec factories of the referenced package are used as `` `alias`.`Name` ``.

## Term definitions (tangl:Definitions)

Same as v0.1. Constants (`means`), variables (`means ... as`), structs (`means` + `has`). Referenced in Cases as `` with `<term>` ``. Basic types (`Text`, `Integer`, `Number`, `Boolean`, `Email`, `Date`, `Time`, `URL`, `Currency`, `Quantity`) are provided by `tangl/types`.

```
- `credit threshold` means 650 as `credit`.`Threshold`
- `applicant` means
  - has `annual income` as Currency
  - has `credit score` as Integer
```

## Rules (tangl:Rules)

Same as v0.1. Inline condition definitions.

```
- `<name>` when <field> <operator> <value>
1. `<name>` when
    - either
      - <condition1>
      - and <condition2>
    - or <condition3>
```

Operators: `is empty`, `is not empty`, `equals`, `is in`, `is greater than`, `is less than`, `is at most`, `is at least`. Logic: `and`, `or`, `not`, `either` (grouping).

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
- `<name>` is a general rule checking <other-case>
```

| Node | Engine | Meaning |
|---|---|---|
| `general rule` | Rule (Warrant) | claim |
| `counter rule` | Counter (Rebuttal) | rebuttal (reverses verdict) |
| `except rule` | Except (Defeater) | exception (attenuates verdict) |

`checking <case>` is **judgment composition** — using another case's verdict as a rule. It **compiles to Evaluate** and fires no side effects. Circular references are forbidden.

> **Implementation note (Active composition rule for case verdicts)**: a case may have multiple warrants (`general rule`). The "case's judgment" that `checking` references is composed by the following rule — among a case's Evaluate results, **if one or more warrants have verdict > 0, the case is considered established (active)**. `until` (§Internal) uses the same rule.

### Defeat edges (don't / do not)

```
- don't `<target>` when `<attacker>`
```
If the attacker is Active, it attacks the target. `do not` is equivalent. → `attacker.Attacks(target)`

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

### The once modifier — tick idempotency (new in v0.3)

```
- do `<func>` once when `<node>`
- do `<func>` once when `<node>` if at least 75% certain
```

**Fires at most once per ctx lifetime.** Since the Internal tick loop shares ctx across ticks, `once` means "once per loop lifetime (until `until` terminates)." For one-shot endpoints (Provides `run`), ctx is created fresh per request, so it is once per request.

- **When it's consumed**: the `once` guard is consumed **only after the firing succeeds**. If the handler errors or is blocked by the certainty gate, it is not consumed — it will be retried on the next tick.
- **Composable with the certainty clause**: `once` comes first, the certainty clause after. Consumption occurs only after passing the gate and firing/succeeding.
- **Implementation (no engine changes needed)**: the codegen inserts a ctx guard inside the generated RunOn handler. The key convention is `once:<subject>.<case>.<node>#<n>` (the n-th `do` of the same node, 0-based). `ctx.Set` after a successful fire. Since ctx flows into cascaded subgraphs, the guard is consistent across the entire cascade.

Why `once` is needed: `every 30s` **re-Runs** a case wholesale every tick. `do arm.placeCup when order received` would **place a cup every tick** while the order stays active. In a coffee robot, sensors naturally update the world so the next stage's guard catches it, but in transfers or filings there is no such sensor — every tick becomes a duplicate transfer or duplicate filing. `once` blocks this accident with a single grammar word.

### The undo edge — partial failure and compensation (new in v0.3)

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

> **Distinction**: `do <func> when <counter-node>` (existing since v0.2) and `undo` are different things. The former is an action for **when a judgment is rejected** (e.g. a deny node becomes Active → `Deny`), while the latter is compensation for **when execution dies partway through**. Judgment rejection is one branch of the normal path; execution failure is an abnormal termination.

### Execution order determinism (refined in v0.3)

v0.2's claim that "`- ` unordered = order-independent, runtime may parallelize" is true only for **judgment statements** (node registration, `don't` edges). Since judgment is pure, the verdict is the same regardless of evaluation order.

Execution is different:

- **Execution firing order is deterministic by the node's document declaration (registration) order.** If two nodes are simultaneously Active, their `do`s fire in document order (engine: handlers fire in rule registration order).
- **Parallelizing execution is forbidden.** If the order of side effects becomes nondeterministic, claims of determinism and auditability break — for reproducibility and auditability to hold, the same document and the same world must always produce the same execution sequence.
- Multiple `do`s on the same node compose into sequential execution within a single handler, in document order (carried over from the v0.2 implementation note).
- **Execution sequences where order carries meaning are recommended to be made explicit via `1.` ordered lists** — this reveals to the reader that meaning is carried by declaration order.

## Cascade temporal semantics — judge-first-execute-after is per level (codified in v0.3)

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

### Go code generation — compensation wrapper (updated in v0.3)

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

### Tick idempotency requirement (new in v0.3, MUST)

**Every `do` of a case reachable from Internal** (either a direct `run` target or belonging to its cascade closure) **must be one of the following:**

- **(a) Idempotent** — the world is the same no matter how many times it executes (e.g., overwriting state, log-like recording)
- **(b) `once`** — guarded to once per ctx lifetime
- **(c) Protected by a counter/sensor guard reflecting completion state** — structured so that the trigger node is no longer Active once execution has finished in the resulting world

Rationale: the tick loop re-Runs a case wholesale every tick. A node whose condition is still true will have its `do` fire again every tick. In a coffee robot, (c) naturally holds — the trigger for `pourWater` is `cup placed`, and pouring water causes the next stage's sensor to update the world. But in domains where **the world does not reflect completion of execution**, such as transfers, filings, or dispatches, neither (a) nor (c) applies — you must use (b) `once` or design a completion counter.

> **Lint recommendation**: an unguarded `do` (with no idempotency marker, no `once`, and no completion guard) on a case reachable from Internal is subject to a parser/linter **warning**.

### verdict → percent conversion

`% = (verdict + 1) × 50` — `+1.0` = 100%, `0.0` = 50%, `-1.0` = 0%. Used in the certainty gate (`if at least N% certain`).

## Example: robotic action control (v0.3)

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

## Example: access control (carried over from v0.2)

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

## Example: bank transfer (compensation) (new in v0.3)

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

## Audit tooling — effect summary (new in v0.3, MUST)

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

cases_stmt = case_decl | required_decl | node_decl | edge_decl | exec_decl
case_decl  = "in" "case" "of" NAME
node_decl  = NAME "is" ("a"|"an") ROLE_EXPR ["using" FUNC_REF [with_clause] [qual_clause] | "checking" NAME]
edge_decl  = DONT NAME "when" NAME                          # defeat edge
exec_decl  = do_edge | run_edge | undo_edge                # v0.3: undo_edge added
do_edge    = "do" FUNC_REF ["once"] "when" NAME [certainty_clause]   # v0.3: optional "once"
undo_edge  = "undo" FUNC_REF "when" NAME                   # v0.3 compensation edge
run_edge   = "run" NAME "when" NAME                        # run <case> when <node>
certainty_clause = "if" THRESHOLD_OP INT "%" "certain"
with_clause = "with" with_arg [("and") with_arg]*
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
```

- The NAME node of an `undo_edge` must be preceded by at least one `do_edge` within the same case (parser-validated).
- Keywords: `this`, `document`, `see`, `means`, `as`, `is`, `a`, `an`, `from`, `general`, `rule`, `counter`, `except`, `case`, `using`, `checking`, `with`, `qualified`, `provides`, `required`, `run`, `check`, `do`, `don't`, `not`, `when`, **`once`**, **`undo`**, `on`, `every`, `until`, `has`, `if`, `at`, `least`, `above`, `less`, `than`, `most`, `certain`, `and`, `or`, `either`, `empty`, `equals`, `in`, `greater`.

## Go API correspondence table (v0.3)

| TANGL | Go |
|---|---|
| `` this document is `x` `` | `package x` |
| `` see `p` from `path` `` | `import p "path"` |
| `` in case of `x` `` | `x := toulmin.NewGraph("subject.x")` |
| `` `x` is a general rule using fn `` | `x := g.Rule(fn)` |
| `` `x` is a counter rule using fn `` | `x := g.Counter(fn)` |
| `` `x` is an except rule using fn `` | `x := g.Except(fn)` |
| `` `x` is a general rule checking `y` `` | `x := g.Rule(checkCase(y))` — **`y` is Evaluate (pure)** |
| `` with `term` `` / `qualified 0.8` | `.With(defs["term"])` / `.Qualifier(0.8)` |
| `` don't `x` when `y` `` | `y.Attacks(x)` |
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
11. **Repetition is only safe on top of idempotency** (v0.3) — every `do` reachable from Internal is either idempotent, `once`, or protected by a completion guard. The tick loop runs only on top of this contract.
12. **Failure is never silent** (v0.3) — partial execution is compensated via `undo` or escalated to a human via REVIEW. Compensation failure is never automatically papered over. Logging always remains.
13. **The execution surface must be statically enumerable** (v0.3) — the complete set of `do`/`undo` reachable from an endpoint (effect summary) can be listed without executing anything. Execution order is also deterministic by document order — same document, same world, same sequence.

## v0.2 → v0.3 migration

| v0.2 | v0.3 |
|---|---|
| Unguarded `do` on a case reachable from Internal (re-fires every tick) | If not idempotent, use `do ... once when` or a completion counter/sensor guard (subject to linter warning) |
| "Execute only once" implicitly relied on world state | Explicit `do <func> once when <node>` |
| Handler error during execution → Run aborts, half-executed state left in place | For every reversible side effect, `undo <func> when <node>` — the wrapper does LIFO compensation, compensation failure goes to REVIEW |
| Both judgment-rejection actions and execution-failure responses routed around via `do ... when <counter-node>` | Rejection actions remain `do ... when <counter-node>`; **only execution-failure compensation moves to `undo`** |
| "Judge-first, execute-after" could be misread as global | Per-level semantics codified (§Cascade Temporal Semantics) — check cascade interleaving during document review |
| Unordered list interpreted as fully parallelizable | Order-independence applies only to judgment statements — recommend converting execution sequences where order matters to `1.` ordered lists |
| Execution surface traced by hand | Static enumeration via `tangl effects <subject>.<endpoint>`, diff gate in CI recommended |
