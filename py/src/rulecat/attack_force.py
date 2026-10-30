from __future__ import annotations
import math
import numbers
from typing import Any


def attack_force(force: Any) -> float:
    """attacks()의 force 인자를 검증한다. [0, 1] 밖·NaN은 ValueError(Go 판의 panic)."""
    if isinstance(force, bool) or not isinstance(force, numbers.Real):
        raise TypeError(
            f"toulmin: force must be a float, got {type(force).__name__}"
        )
    f = float(force)
    if math.isnan(f) or f < 0.0 or f > 1.0:
        raise ValueError("toulmin: force must be between 0.0 and 1.0")
    return f
