import type { Strength, Specs, RuleFunc } from "./types.js";
export interface RuleMeta {
  name: string;
  qualifier: number;
  strength: Strength;
  defeats: string[];
  specs: Specs;
  fn: RuleFunc;
}
