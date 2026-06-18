import type { Strength, Specs, RuleFunc, NodeHandler } from "./types.js";
import type { Graph } from "./graph.js";   // 타입 전용(순환 import 안전)
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
  runGraph?: Graph;         // 이 노드가 Active면 Run할 하위 그래프
}
