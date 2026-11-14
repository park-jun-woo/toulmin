import { describe, it, expect } from "vitest";
import { attackForce } from "../attack-force.js";

describe("attackForce", () => {
  it("defaults to 1 when no value is given", () => {
    expect(attackForce([])).toBe(1);
  });

  it("returns the single value inside [0, 1], bounds included", () => {
    expect(attackForce([0])).toBe(0);
    expect(attackForce([0.3])).toBe(0.3);
    expect(attackForce([1])).toBe(1);
  });

  it("throws on more than one value", () => {
    expect(() => attackForce([0.5, 0.5])).toThrow("toulmin: attacks accepts at most one force value");
  });

  it("throws on NaN, out-of-range, and non-number values", () => {
    for (const bad of [NaN, -0.1, 1.1, Infinity, "0.5" as unknown as number]) {
      expect(() => attackForce([bad])).toThrow("toulmin: force must be between 0.0 and 1.0");
    }
  });
});
