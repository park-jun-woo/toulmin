from __future__ import annotations
import json
from typing import Any

from rulecat.func_id import func_id


def rule_id(fn: Any, specs: list[Any]) -> str:
    fid = func_id(fn)
    if not specs:
        return fid
    names = sorted(json.dumps(s.__dict__, default=str) for s in specs)
    return fid + "#" + "+".join(names)
