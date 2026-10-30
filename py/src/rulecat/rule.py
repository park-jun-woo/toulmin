from __future__ import annotations
import math
from typing import TYPE_CHECKING, Any

from rulecat.attack_force import attack_force
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

    def attacks(self, target: Rule, force: float = 1.0) -> None:
        """이 규칙이 target을 공격한다는 간선을 선언한다.

        force F ∈ [0, 1](기본 1)는 공격의 세기다. 반박(counter·rule 공격자)은
        target의 분모에 F·(v+1)/2를 더하고, 예외(except)는 target의 적용도에
        1 − F·(v+1)/2를 곱한다. NaN·범위 밖이면 ValueError.
        """
        from rulecat.defeat_edge import DefeatEdge
        f = attack_force(force)
        self._graph.defeats.append(DefeatEdge(from_=self.id, to=target.id, force=f))

    def supersedes(self, target: Rule) -> None:
        """이 규칙이 target을 크리스프하게 배제(supersede)한다는 간선을 선언한다.

        이 규칙의 적용도 a_eff > 0일 때만 배제한다. a_eff = (v+1)/2 · Π_예외
        (1 − F·(v_j+1)/2)로 예외 게이트는 반영하고 반박(counter)은 반영하지 않으며,
        strict 규칙은 a_eff = (v+1)/2다. 예외가 이 규칙을 완전히 막으면 배제도 풀린다.
        간선은 h-Categoriser 순회에 들어가지 않고 TraceEntry.excluded만 정한다.
        """
        from rulecat.defeat_edge import DefeatEdge
        self._graph.supersedes.append(DefeatEdge(from_=self.id, to=target.id))

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
        for edge in [*self._graph.defeats, *self._graph.supersedes]:
            if edge.from_ == old_id:
                edge.from_ = new_id
            if edge.to == old_id:
                edge.to = new_id
        return self

    def qualifier(self, q: float) -> Rule:
        if math.isnan(q) or q < 0.0 or q > 1.0:
            raise ValueError("toulmin: qualifier must be between 0.0 and 1.0")
        self._graph.rules[self._idx].qualifier = q
        return self

    def run_on(self, h: Any) -> "Rule":
        self._graph.rules[self._idx].run_on = h
        return self

    def run(self, g: Any) -> "Rule":  # g: Graph — Active면 Run할 하위 그래프
        if g is None:
            raise ValueError("toulmin: run requires a non-None sub-graph")
        self._graph.rules[self._idx].run_graph = g
        return self
