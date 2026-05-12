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
