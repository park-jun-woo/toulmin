export interface Context {
  get(key: string): unknown;
  set(key: string, value: unknown): void;
}

export interface Spec {
  specName(): string;
  validate(): void;
}

export type Specs = Spec[];

export function findSpec(specs: Specs, name: string): Spec | undefined {
  return specs.find(s => s.specName() === name);
}

/**
 * A rule function returns [verdict, evidence].
 *
 * verdict is a number in [-1, 1]: -1 = not applicable, +1 = fully applies.
 * **Not applicable is -1, not 0. 0 means "half applies".** The engine cannot
 * detect a rule that returns 0 by mistake; it is evaluated as half applicable.
 *
 * A legacy boolean is still accepted: true → +1, false → -1 (same verdicts as before).
 * A verdict outside [-1, 1], NaN, or a non-number/non-boolean is an evaluation error.
 */
export type RuleFunc = (ctx: Context, specs: Specs) => [number | boolean, unknown];

export enum Strength {
  Defeasible = 0,
  Strict = 1,
  Defeater = 2,
}

export enum EvalMethod {
  Matrix = 0,
}

export interface EvalOption {
  method?: EvalMethod;
  trace?: boolean;
  duration?: boolean;
}

export interface EvalResult {
  name: string;
  verdict: number;
  evidence?: unknown;
  /**
   * Effective applicability a_eff in [0, 1]: (v+1)/2 of the warrant's own rule
   * verdict, multiplied by every Except gate (1 − F·(v_j+1)/2).
   */
  applicability: number;
  trace?: TraceEntry[];
}

export interface TraceEntry {
  name: string;        // = Claim
  role: string;        // rule | counter | except
  activated: boolean;  // derived: ruleVerdict > -1
  ruleVerdict: number; // the rule function's own v in [-1, 1] (before attacks); -1 if never ran
  qualifier: number;
  verdict: number;     // ← 추가
  evidence?: unknown;
  ground?: unknown;    // ← 추가: ctx 그대로
  specs?: Specs;       // = Backing
  duration?: number;
  /**
   * Crisply superseded out of the tally: some superseder of this node has
   * effective applicability a_eff > 0. Reporting-only — the verdict is unchanged.
   */
  excluded?: boolean;
}

export interface Trace {
  all(): TraceEntry[];                       // 전 노드(등록순)
  get(name: string): TraceEntry | undefined; // 짧은 이름 조회
  ctx(): Context;                            // 이 Run의 컨텍스트
}

export type NodeHandler = (self: TraceEntry, t: Trace) => void;

export interface RunResult {
  results: EvalResult[];
  trace: Trace;            // full pass 후 전 노드 Trace (등록 순서)
}

export type Expectation = (results: EvalResult[]) => void;

export interface TestCase {
  name: string;
  context?: Context;
  option?: EvalOption;
  expect: Expectation;
}
