from __future__ import annotations
import math


def format_float(v: float) -> str:
    """Go %v 형식으로 float를 문자열로 바꾼다(오류 메시지를 Go 판과 맞추기 위함).

    NaN → "NaN", ±∞ → "+Inf"/"-Inf", 정수값은 소수점 없이("2"), 나머지는 repr.
    """
    if math.isnan(v):
        return "NaN"
    if math.isinf(v):
        return "+Inf" if v > 0 else "-Inf"
    if v == int(v) and abs(v) < 1e21:
        return str(int(v))
    return repr(v)
