import type { Context } from "./types.js";
export class MapContext implements Context {
  private data = new Map<string, unknown>();
  get(key: string): unknown { return this.data.get(key); }
  set(key: string, value: unknown): void { this.data.set(key, value); }
}
export function newContext(): MapContext { return new MapContext(); }
