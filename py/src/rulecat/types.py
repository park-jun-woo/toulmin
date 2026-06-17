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
    name: str = ""
    role: str = ""
    activated: bool = False
    qualifier: float = 1.0
    evidence: Any = None
    specs: list[Spec] = field(default_factory=list)
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


class NodeEventType(IntEnum):
    INACTIVE = 0  # 비활성실행
    ACTIVE = 1    # 활성실행 — verdict > 0
    DEFEATED = 2  # 무력화실행 — verdict <= 0


@dataclass
class NodeEvent:
    name: str = ""
    role: str = ""              # "rule" | "counter" | "except"
    type: NodeEventType = NodeEventType.INACTIVE
    verdict: float = 0.0        # INACTIVE면 의미 없음
    evidence: Any = None


@runtime_checkable
class RunView(Protocol):
    def all(self) -> list[NodeEvent]: ...               # 전 노드, 등록 순서
    def get(self, name: str) -> NodeEvent | None: ...   # 표시명(short_name)으로 조회
    def attackers(self, name: str) -> list[NodeEvent]: ...  # name을 공격한 노드 이벤트들


NodeHandler = Callable[[Context, "NodeEvent", "RunView"], None]


@dataclass
class RunResult:
    results: list[EvalResult] = field(default_factory=list)
    view: "RunView | None" = None
