from __future__ import annotations
import math
import numbers
from typing import Any

from rulecat.format_float import format_float


def to_verdict(v: Any) -> float:
    """규칙 함수의 반환값을 verdict [-1, 1]로 바꾼다.

    - bool(옛 시그니처): True → +1(완전히 해당), False → -1(해당 없음).
      bool은 int의 하위 타입이므로 숫자 검사보다 먼저 본다.
    - 실수(int·float): float로 바꾼다. 0은 "반쯤 해당"이다(해당 없음은 -1).
    - NaN·범위 밖은 ValueError, 숫자가 아니면 TypeError.
    """
    if isinstance(v, bool):
        return 1.0 if v else -1.0
    if not isinstance(v, numbers.Real):
        raise TypeError(
            f"rule verdict must be bool or float, got {type(v).__name__}"
        )
    f = float(v)
    if math.isnan(f) or f < -1.0 or f > 1.0:
        raise ValueError(f"rule verdict {format_float(f)} out of range [-1, 1]")
    return f
