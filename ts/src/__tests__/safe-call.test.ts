import { describe, it, expect } from "vitest";
import { safeCall } from "../safe-call.js";
import { newContext } from "../map-context.js";
import type { RuleFunc } from "../types.js";

const call = (fn: RuleFunc) => safeCall(fn, newContext(), []);

describe("safeCall", () => {
  it("maps boolean true → +1 and false → -1, passing evidence through", () => {
    expect(call(() => [true, "e1"])).toEqual({ ok: true, verdict: 1, evidence: "e1" });
    expect(call(() => [false, "e2"])).toEqual({ ok: true, verdict: -1, evidence: "e2" });
  });

  it("passes a number in [-1, 1] through unchanged (0 = half applies)", () => {
    for (const v of [-1, -0.5, 0, 0.25, 1]) {
      expect(call(() => [v, null])).toEqual({ ok: true, verdict: v, evidence: null });
    }
  });

  it("rejects out-of-range and NaN verdicts with the Go-compatible message", () => {
    expect(call(() => [1.5, null])).toEqual({ ok: false, error: "rule verdict 1.5 out of range [-1, 1]" });
    expect(call(() => [-2, null])).toEqual({ ok: false, error: "rule verdict -2 out of range [-1, 1]" });
    expect(call(() => [NaN, null])).toEqual({ ok: false, error: "rule verdict NaN out of range [-1, 1]" });
  });

  it("rejects a verdict that is neither number nor boolean", () => {
    const fn = (() => [undefined, null]) as unknown as RuleFunc;
    expect(call(fn)).toEqual({ ok: false, error: "rule verdict must be a number or boolean, got undefined" });
  });

  it("captures a thrown error, and a non-tuple return, as an error string", () => {
    expect(call(() => { throw new Error("boom"); })).toEqual({ ok: false, error: "Error: boom" });
    const bad = (() => undefined) as unknown as RuleFunc;
    const out = call(bad);
    expect(out.ok).toBe(false);
    if (!out.ok) expect(out.error).toMatch(/TypeError/);
  });

  it("passes ctx and specs to the rule function", () => {
    const ctx = newContext();
    ctx.set("k", 0.5);
    const specs = [{ specName: () => "s", validate: () => {} }];
    let seen: unknown[] = [];
    const out = safeCall((c, s) => { seen = [c, s]; return [c.get("k") as number, null]; }, ctx, specs);
    expect(out).toEqual({ ok: true, verdict: 0.5, evidence: null });
    expect(seen[0]).toBe(ctx);
    expect(seen[1]).toBe(specs);
  });
});
