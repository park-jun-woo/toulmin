import type { Strength, Specs, RuleFunc, NodeHandler } from "./types.js";
import type { Graph } from "./graph.js";   // 타입 전용(순환 import 안전)
export interface RuleMeta {
  name: string;
  qualifier: number;
  strength: Strength;
  defeats: string[];
  specs: Specs;
  fn: RuleFunc;
  runOn?: NodeHandler;      // Active(activated && verdict>0) 시 발화하는 핸들러
  runGraph?: Graph;         // 이 노드가 Active면 Run할 하위 그래프
}
