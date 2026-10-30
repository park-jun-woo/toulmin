/**
 * Resolves the optional force argument of Rule.attacks.
 * No value → 1. Throws if more than one value is given, or the value is NaN
 * or outside [0.0, 1.0] (same convention as qualifier).
 */
export function attackForce(force: number[]): number {
  if (force.length === 0) return 1.0;
  if (force.length > 1) throw new Error("toulmin: attacks accepts at most one force value");
  const f = force[0];
  if (typeof f !== "number" || Number.isNaN(f) || f < 0.0 || f > 1.0) {
    throw new Error("toulmin: force must be between 0.0 and 1.0");
  }
  return f;
}
