import { describe, it, expect } from "vitest";
import { findSpec, type Spec } from "../src/index.js";

class TestSpec implements Spec {
  constructor(public value: string) {}
  specName() { return "TestSpec"; }
  validate() {}
}

describe("findSpec", () => {
  it("finds by specName", () => {
    const specs = [new TestSpec("a")];
    const found = findSpec(specs, "TestSpec");
    expect(found).toBeDefined();
    expect((found as TestSpec).value).toBe("a");
  });

  it("returns undefined when not found", () => {
    expect(findSpec([], "nope")).toBeUndefined();
  });
});
