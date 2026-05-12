from __future__ import annotations


def detect_cycle(edges: dict[str, list[str]]) -> str | None:
    WHITE, GRAY, BLACK = 0, 1, 2
    color: dict[str, int] = {}

    def dfs(node: str) -> str | None:
        color[node] = GRAY
        for neighbor in edges.get(node, []):
            c = color.get(neighbor, WHITE)
            if c == GRAY:
                return f"cycle detected at {neighbor}"
            if c == WHITE:
                err = dfs(neighbor)
                if err:
                    return err
        color[node] = BLACK
        return None

    for node in edges:
        if color.get(node, WHITE) == WHITE:
            err = dfs(node)
            if err:
                return err
    return None
