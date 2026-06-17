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
  name: string;
  role: string;
  activated: boolean;
  qualifier: number;
  evidence?: unknown;
  specs?: Specs;
  duration?: number;
}

export enum NodeEventType {
  Inactive = 0, // 비활성실행
  Active = 1,   // 활성실행 — verdict > 0
  Defeated = 2, // 무력화실행 — verdict <= 0
}

export interface NodeEvent {
  name: string;
  role: string;            // "rule" | "counter" | "except"
  type: NodeEventType;
  verdict: number;         // Inactive면 의미 없음
  evidence?: unknown;
}

export interface RunView {
  all(): NodeEvent[];                       // 전 노드, 등록 순서 (Inactive 포함)
  get(name: string): NodeEvent | undefined; // 표시명(shortName)으로 조회
  attackers(name: string): NodeEvent[];     // name을 공격한 노드 이벤트들
}

export type NodeHandler = (ctx: Context, ev: NodeEvent, view: RunView) => void;

export interface RunResult {
  results: EvalResult[];
  view: RunView;           // full pass 후 전 노드 최종 상태 스냅샷
}

export type Expectation = (results: EvalResult[]) => void;

export interface TestCase {
  name: string;
  context?: Context;
  option?: EvalOption;
  expect: Expectation;
}
