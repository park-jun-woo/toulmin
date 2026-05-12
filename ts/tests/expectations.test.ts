import { describe, it, expect } from "vitest";
import { verdictAbove, verdictAtMost, verdictBetween, noResult, type EvalResult } from "../src/index.js";

describe("Expectations", () => {
  const result = (v: number): EvalResult[] => [{ name: "test", verdict: v }];

  it("verdictAbove passes when verdict > threshold", () => {
    expect(() => verdictAbove(0)(result(0.5))).not.toThrow();
    expect(() => verdictAbove(0)(result(-0.1))).toThrow();
    expect(() => verdictAbove(0)([])).toThrow();
  });

  it("verdictAtMost passes when verdict <= threshold", () => {
    expect(() => verdictAtMost(0)(result(0))).not.toThrow();
    expect(() => verdictAtMost(0)(result(0.1))).toThrow();
  });

  it("verdictBetween passes when lo < verdict <= hi", () => {
    expect(() => verdictBetween(0, 0.5)(result(0.3))).not.toThrow();
    expect(() => verdictBetween(0, 0.5)(result(0))).toThrow();
    expect(() => verdictBetween(0, 0.5)(result(0.6))).toThrow();
  });

  it("noResult passes when empty", () => {
    expect(() => noResult([])).not.toThrow();
    expect(() => noResult(result(1))).toThrow();
  });
});
