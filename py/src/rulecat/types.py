from __future__ import annotations
from dataclasses import dataclass, field
from enum import IntEnum
from typing import Any, Callable, Protocol, runtime_checkable, TYPE_CHECKING

if TYPE_CHECKING:
    from rulecat.trace import Trace


@runtime_checkable
class Context(Protocol):
    def get(self, key: str) -> Any: ...
    def set(self, key: str, value: Any) -> None: ...


@runtime_checkable
class Spec(Protocol):
    def spec_name(self) -> str: ...
    def validate(self) -> None: ...


# 규칙 함수: (ctx, specs) -> (verdict, evidence), verdict ∈ [-1, 1].
# +1 = 완전히 해당, 0 = 반쯤 해당, -1 = 해당 없음.
# 해당 없음은 0이 아니라 -1이다. 0을 돌려주면 규칙이 Active로 남아 공격자까지
# 평가되고 판정에 계속 반영된다. 0은 합법 값이라 엔진이 이 실수를 검출하지 못한다.
# 옛 bool 반환도 받는다: True → +1, False → -1.
# NaN·범위 밖이면 평가 오류(RuntimeError)다.
RuleFunc = Callable[[Context, list[Spec]], tuple[float | bool, Any]]


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
    activated: bool = False  # 파생값: rule_verdict > -1
    rule_verdict: float = -1.0  # 규칙 함수 자신의 v ∈ [-1, 1](공격 반영 전). 실행 안 됐으면 -1
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
    # 유효 적용도 a_eff ∈ [0, 1] = (v+1)/2 · Π_except (1 − F·(v_j+1)/2)
    applicability: float = 0.0
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


# self = 발화 노드의 TraceEntry(엔진이 직접 전달), t = 전체 Trace
NodeHandler = Callable[["TraceEntry", "Trace"], None]


@dataclass
class RunResult:
    results: list[EvalResult] = field(default_factory=list)
    trace: "Trace | None" = None
