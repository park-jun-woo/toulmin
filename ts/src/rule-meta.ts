import type { Strength, Specs, RuleFunc, NodeHandler } from "./types.js";
export interface RuleMeta {
  name: string;
  qualifier: number;
  strength: Strength;
  defeats: string[];
  specs: Specs;
  fn: RuleFunc;
  onActive?: NodeHandler;   // 활성실행
  onDefeated?: NodeHandler; // 무력화실행
  onInactive?: NodeHandler; // 비활성실행
}
