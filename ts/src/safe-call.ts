import type { Context, RuleFunc, Specs } from "./types.js";

export type CallOutcome =
  | { ok: true; verdict: number; evidence: unknown }
  | { ok: false; error: string };

/**
 * Runs a rule function and normalizes its verdict.
 * boolean → ±1 (true → +1, false → -1). A number must lie in [-1, 1] (NaN rejected).
 * A thrown error, an out-of-range verdict, or any other verdict type is an error;
 * the caller wraps it as `rule "<id>": <error>`.
 */
export function safeCall(fn: RuleFunc, ctx: Context, specs: Specs): CallOutcome {
  let raw: unknown;
  let evidence: unknown;
  try {
    [raw, evidence] = fn(ctx, specs);
  } catch (e) {
    return { ok: false, error: String(e) };
  }
  if (typeof raw === "boolean") return { ok: true, verdict: raw ? 1.0 : -1.0, evidence };
  if (typeof raw !== "number") {
    return { ok: false, error: `rule verdict must be a number or boolean, got ${typeof raw}` };
  }
  if (Number.isNaN(raw) || raw < -1.0 || raw > 1.0) {
    return { ok: false, error: `rule verdict ${raw} out of range [-1, 1]` };
  }
  return { ok: true, verdict: raw, evidence };
}
