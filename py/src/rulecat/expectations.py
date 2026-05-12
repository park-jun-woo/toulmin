from __future__ import annotations
from rulecat.types import EvalResult, Expectation


def verdict_above(v: float) -> Expectation:
    def check(results: list[EvalResult]) -> None:
        if not results:
            raise AssertionError("expected at least one result")
        if results[0].verdict <= v:
            raise AssertionError(f"expected verdict > {v}, got {results[0].verdict}")
    return check


def verdict_at_most(v: float) -> Expectation:
    def check(results: list[EvalResult]) -> None:
        if not results:
            raise AssertionError("expected at least one result")
        if results[0].verdict > v:
            raise AssertionError(f"expected verdict <= {v}, got {results[0].verdict}")
    return check


def verdict_between(lo: float, hi: float) -> Expectation:
    def check(results: list[EvalResult]) -> None:
        if not results:
            raise AssertionError("expected at least one result")
        verdict = results[0].verdict
        if verdict <= lo or verdict > hi:
            raise AssertionError(f"expected {lo} < verdict <= {hi}, got {verdict}")
    return check


def no_result(results: list[EvalResult]) -> None:
    if results:
        raise AssertionError(f"expected no results, got {len(results)}")
