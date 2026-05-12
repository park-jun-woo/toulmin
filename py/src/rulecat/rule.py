from __future__ import annotations
from typing import TYPE_CHECKING, Any

from rulecat.rule_id import rule_id

if TYPE_CHECKING:
    from rulecat.graph import Graph
    from rulecat.types import Spec


class Rule:
    def __init__(self, id: str, graph: Graph, idx: int, fn: Any) -> None:
        self.id = id
        self._graph = graph
        self._idx = idx
        self._fn = fn

    def attacks(self, target: Rule) -> None:
        from rulecat.defeat_edge import DefeatEdge
        self._graph.defeats.append(DefeatEdge(from_=self.id, to=target.id))

    def with_spec(self, spec: Spec) -> Rule:
        spec.validate()
        meta = self._graph.rules[self._idx]
        meta.specs.append(spec)
        old_id = meta.name
        new_id = rule_id(self._fn, meta.specs)
        meta.name = new_id
        self.id = new_id
        role = self._graph.roles.pop(old_id, None)
        if role is not None:
            self._graph.roles[new_id] = role
        for edge in self._graph.defeats:
            if edge.from_ == old_id:
                edge.from_ = new_id
            if edge.to == old_id:
                edge.to = new_id
        return self

    def qualifier(self, q: float) -> Rule:
        if q < 0.0 or q > 1.0:
            raise ValueError("qualifier must be between 0.0 and 1.0")
        self._graph.rules[self._idx].qualifier = q
        return self
