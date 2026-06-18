from __future__ import annotations
from typing import Any


def detect_run_cycle(root: Any) -> str | None:  # root: Graph
    WHITE, GRAY, BLACK = 0, 1, 2
    color: dict[Any, int] = {}  # Graph 객체(identity 해시)로 키잉

    def dfs(g: Any) -> str | None:
        color[g] = GRAY
        for r in g.rules:
            sub = r.run_graph
            if sub is None:
                continue
            c = color.get(sub, WHITE)
            if c == GRAY:
                return f'run cycle detected at graph "{sub.name}"'
            if c == WHITE:
                err = dfs(sub)
                if err:
                    return err
        color[g] = BLACK
        return None

    return dfs(root)
