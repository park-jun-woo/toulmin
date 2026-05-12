import type { RuleFunc, Specs } from "./types.js";
import { funcID } from "./func-id.js";

export function ruleID(fn: RuleFunc, specs: Specs): string {
  const id = funcID(fn);
  if (specs.length === 0) return id;
  const names = specs.map(s => JSON.stringify(s)).sort();
  return id + "#" + names.join("+");
}
