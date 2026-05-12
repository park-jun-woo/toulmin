export function shortName(full: string): string {
  const hashIdx = full.indexOf("#");
  let base = hashIdx >= 0 ? full.substring(0, hashIdx) : full;
  const spec = hashIdx >= 0 ? full.substring(hashIdx) : "";
  const dotIdx = base.lastIndexOf(".");
  if (dotIdx >= 0) base = base.substring(dotIdx + 1);
  return base + spec;
}
