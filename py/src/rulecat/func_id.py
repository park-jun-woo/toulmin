from __future__ import annotations
from typing import Any

_registry: dict[int, str] = {}
_counter = 0


def func_id(fn: Any) -> str:
    global _counter
    obj_id = id(fn)
    if obj_id in _registry:
        return _registry[obj_id]
    _counter += 1
    name = getattr(fn, "__qualname__", None) or f"fn_{_counter}"
    uid = f"{name}_{_counter}"
    _registry[obj_id] = uid
    return uid
