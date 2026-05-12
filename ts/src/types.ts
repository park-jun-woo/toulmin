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

export type Expectation = (results: EvalResult[]) => void;

export interface TestCase {
  name: string;
  context?: Context;
  option?: EvalOption;
  expect: Expectation;
}
