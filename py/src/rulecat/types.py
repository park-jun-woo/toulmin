from __future__ import annotations
from dataclasses import dataclass, field
from enum import IntEnum
from typing import Any, Callable, Protocol, runtime_checkable


@runtime_checkable
class Context(Protocol):
    def get(self, key: str) -> Any: ...
    def set(self, key: str, value: Any) -> None: ...


@runtime_checkable
class Spec(Protocol):
    def spec_name(self) -> str: ...
    def validate(self) -> None: ...


RuleFunc = Callable[[Context, list[Spec]], tuple[bool, Any]]


class Strength(IntEnum):
    DEFEASIBLE = 0
    STRICT = 1
    DEFEATER = 2


class EvalMethod(IntEnum):
    MATRIX = 0


@dataclass
class EvalOption:
    method: EvalMethod = EvalMethod.MATRIX
    trace: bool = False
    duration: bool = False


@dataclass
class TraceEntry:
    name: str = ""           # = Claim
    role: str = ""           # rule | counter | except
    activated: bool = False
    qualifier: float = 1.0
    verdict: float = 0.0     # 노드 우세/패배 판별용
    evidence: Any = None
    ground: Any = None       # ctx 그대로
    specs: list[Spec] = field(default_factory=list)  # = Backing
    duration: float | None = None


@dataclass
class EvalResult:
    name: str = ""
    verdict: float = 0.0
    evidence: Any = None
    trace: list[TraceEntry] = field(default_factory=list)


# EvalResult 이후 -- 타입 별칭은 런타임 즉시 평가
Expectation = Callable[[list[EvalResult]], None]


@dataclass
class TestCase:
    name: str = ""
    context: Context | None = None
    option: EvalOption = field(default_factory=EvalOption)
    expect: Expectation = lambda r: None


def find_spec(specs: list[Spec], name: str) -> Spec | None:
    return next((s for s in specs if s.spec_name() == name), None)


# self = 발화한 노드의 TraceEntry, trace = 전 노드 TraceEntry 목록(조회용)
NodeHandler = Callable[[Context, "TraceEntry", list["TraceEntry"]], None]


@dataclass
class RunResult:
    results: list[EvalResult] = field(default_factory=list)
    trace: list[TraceEntry] = field(default_factory=list)
