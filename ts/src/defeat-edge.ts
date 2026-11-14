export interface DefeatEdge {
  from: string;
  to: string;
  force: number;   // attack force F in [0, 1] (default 1)
}

/** A supersede edge: from crisply excludes to (reporting-only, never traversed). */
export interface SupersedeEdge {
  from: string;
  to: string;
}
