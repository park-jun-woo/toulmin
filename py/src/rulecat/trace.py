from __future__ import annotations
from rulecat.types import Context, TraceEntry


class Trace:
    """한 Run의 읽기 전용 뷰: 전 노드 entry + ctx."""

    def __init__(self, nodes: list[TraceEntry], ctx: Context) -> None:
        self._nodes = nodes
        self._ctx = ctx

    def all(self) -> list[TraceEntry]:
        return self._nodes

    def get(self, name: str) -> TraceEntry | None:
        for e in self._nodes:
            if e.name == name:
                return e
        return None

    def ctx(self) -> Context:
        return self._ctx
