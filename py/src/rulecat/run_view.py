from __future__ import annotations

from rulecat.defeat_edge import DefeatEdge
from rulecat.short_name import short_name
from rulecat.types import NodeEvent


class _RunView:
    """전 노드 최종 이벤트의 읽기 전용 스냅샷. 접근자는 복사본을 반환한다."""

    def __init__(
        self,
        order: list[NodeEvent],
        by_name: dict[str, NodeEvent],
        attackers: dict[str, list[NodeEvent]],
    ) -> None:
        self._order = order              # list[NodeEvent]
        self._by_name = by_name          # dict[str, NodeEvent]
        self._attackers = attackers      # dict[str, list[NodeEvent]]

    def all(self) -> list[NodeEvent]:
        return list(self._order)         # 복사본

    def get(self, name: str) -> NodeEvent | None:
        return self._by_name.get(name)

    def attackers(self, name: str) -> list[NodeEvent]:
        return list(self._attackers.get(name, []))


def _build_attacker_events(
    defeats: list[DefeatEdge], by_name: dict[str, NodeEvent]
) -> dict[str, list[NodeEvent]]:
    """self.defeats를 to(short_name)별로 그룹화해 각 from_의 이벤트 목록을 모은다.

    간선 방향이 이미 ``to ← from_``이므로 뒤집지 않고 to로 group-by 한다.
    """
    attackers: dict[str, list[NodeEvent]] = {}
    for d in defeats:
        to = short_name(d.to)
        ev = by_name.get(short_name(d.from_))
        if ev is None:
            continue
        attackers.setdefault(to, []).append(ev)
    return attackers
