import type { Expectation } from "./types.js";

export function verdictAbove(v: number): Expectation {
  return (results) => {
    if (results.length === 0) throw new Error("expected at least one result");
    if (results[0].verdict <= v) throw new Error(`expected verdict > ${v}, got ${results[0].verdict}`);
  };
}

export function verdictAtMost(v: number): Expectation {
  return (results) => {
    if (results.length === 0) throw new Error("expected at least one result");
    if (results[0].verdict > v) throw new Error(`expected verdict <= ${v}, got ${results[0].verdict}`);
  };
}

export function verdictBetween(lo: number, hi: number): Expectation {
  return (results) => {
    if (results.length === 0) throw new Error("expected at least one result");
    const verdict = results[0].verdict;
    if (verdict <= lo || verdict > hi) throw new Error(`expected ${lo} < verdict <= ${hi}, got ${verdict}`);
  };
}

export const noResult: Expectation = (results) => {
  if (results.length !== 0) throw new Error(`expected no results, got ${results.length}`);
};
