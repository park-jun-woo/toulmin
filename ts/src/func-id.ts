import type { RuleFunc } from "./types.js";
const funcRegistry = new WeakMap<Function, string>();
let funcCounter = 0;

export function funcID(fn: RuleFunc): string {
  let id = funcRegistry.get(fn);
  if (!id) {
    id = `${fn.name || "fn"}_${++funcCounter}`;
    funcRegistry.set(fn, id);
  }
  return id;
}
