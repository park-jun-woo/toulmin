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

export type RuleFunc = (ctx: Context, specs: Specs) => [boolean, unknown];

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
  trace?: TraceEntry[];
}

export interface TraceEntry {
  name: string;        // = Claim
  role: string;        // rule | counter | except
  activated: boolean;
  qualifier: number;
  verdict: number;     // ← 추가
  evidence?: unknown;
  ground?: unknown;    // ← 추가: ctx 그대로
  specs?: Specs;       // = Backing
  duration?: number;
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
