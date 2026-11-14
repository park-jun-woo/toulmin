from dataclasses import dataclass


@dataclass
class DefeatEdge:
    from_: str
    to: str
    force: float = 1.0  # 공격 강도 F ∈ [0, 1]
