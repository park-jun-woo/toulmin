import math
import unittest

import pytest

from rulecat import Graph, MapContext


def _r(ctx, specs):
    return (True, None)


class RuleRunTest(unittest.TestCase):
    def test_run_sets_sub_graph_and_returns_self(self):
        sub = Graph("sub")
        parent = Graph("parent")
        rule = parent.rule(_r)
        returned = rule.run(sub)
        self.assertIs(returned, rule)
        self.assertIs(parent.rules[0].run_graph, sub)

    def test_run_none_raises(self):
        g = Graph("g")
        rule = g.rule(_r)
        with self.assertRaises(ValueError) as cm:
            rule.run(None)
        self.assertIn("non-None sub-graph", str(cm.exception))


class RuleRunOnTest(unittest.TestCase):
    def test_run_on_sets_handler_and_returns_self(self):
        g = Graph("g")
        rule = g.rule(_r)
        h = lambda ctx, self_, trace: None
        returned = rule.run_on(h)
        self.assertIs(returned, rule)
        self.assertIs(g.rules[0].run_on, h)



# --- attacks force 검증 ---


class Tag:
    def __init__(self, value: str):
        self.value = value

    def spec_name(self) -> str:
        return "Tag"

    def validate(self) -> None:
        pass


def _const(v, ev=None):
    return lambda ctx, specs: (v, ev)


def test_attack_force_default_is_one():
    g = Graph("t")
    w = g.rule(_const(True))
    c = g.counter(_const(True))
    c.attacks(w)
    assert g.defeats[0].force == 1.0
    assert g.evaluate(MapContext())[0].verdict == 0.0


@pytest.mark.parametrize("bad", [-0.1, 1.1, float("nan"), float("inf")])
def test_attack_force_rejects_invalid(bad):
    g = Graph("t")
    w = g.rule(_const(True))
    c = g.counter(_const(True))
    with pytest.raises(ValueError, match="force must be between 0.0 and 1.0"):
        c.attacks(w, bad)
    assert g.defeats == []


def test_attack_force_survives_with_spec_rename():
    g = Graph("t")
    w = g.rule(_const(True))
    c = g.counter(_const(True))
    c.attacks(w, 0.5)
    c.with_spec(Tag("x"))
    assert g.defeats[0].force == 0.5
    assert math.isclose(g.evaluate(MapContext())[0].verdict, 1.0 / 3.0)


if __name__ == "__main__":
    unittest.main()
