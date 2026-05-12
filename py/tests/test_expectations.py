import pytest
from rulecat import EvalResult, verdict_above, verdict_at_most, verdict_between, no_result


def _result(v: float) -> list[EvalResult]:
    return [EvalResult(name="test", verdict=v)]


def test_verdict_above():
    verdict_above(0)(_result(0.5))
    with pytest.raises(AssertionError):
        verdict_above(0)(_result(-0.1))
    with pytest.raises(AssertionError):
        verdict_above(0)([])


def test_verdict_at_most():
    verdict_at_most(0)(_result(0))
    with pytest.raises(AssertionError):
        verdict_at_most(0)(_result(0.1))


def test_verdict_between():
    verdict_between(0, 0.5)(_result(0.3))
    with pytest.raises(AssertionError):
        verdict_between(0, 0.5)(_result(0))
    with pytest.raises(AssertionError):
        verdict_between(0, 0.5)(_result(0.6))


def test_no_result():
    no_result([])
    with pytest.raises(AssertionError):
        no_result(_result(1.0))
