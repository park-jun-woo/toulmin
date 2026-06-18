from __future__ import annotations
from dataclasses import dataclass, field
from typing import Any


@dataclass
class RuleMeta:
    name: str
    qualifier: float
    strength: int  # Strength enum value
    defeats: list[str] = field(default_factory=list)
    specs: list[Any] = field(default_factory=list)
    fn: Any = None  # RuleFunc
    run_on: Any = None        # NodeHandler | None — Active면 발화
    run_graph: Any = None     # "Graph | None" — Active면 Run할 하위 그래프 (순환 import 회피로 Any)
