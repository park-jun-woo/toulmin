import pytest
from rulecat import Graph, MapContext, Rule


class TestSpec:
    def __init__(self, value: str):
        self.value = value
    def spec_name(self) -> str:
        return "TestSpec"
    def validate(self) -> None:
        pass


def test_register_rule_counter_except():
    fn1 = lambda ctx, specs: (True, None)
    fn2 = lambda ctx, specs: (True, None)
    fn3 = lambda ctx, specs: (True, None)
    g = Graph("test")
    g.rule(fn1)
    g.counter(fn2)
    g.except_(fn3)
    assert len(g.rules) == 3


def test_duplicate_registration():
    fn = lambda ctx, specs: (True, None)
    g = Graph("test")
    g.rule(fn)
    with pytest.raises(ValueError, match="duplicate"):
        g.rule(fn)


def test_with_spec_updates_defeat_edges():
    fn1 = lambda ctx, specs: (True, None)
    fn2 = lambda ctx, specs: (True, None)
    g = Graph("test")
    r1 = g.rule(fn1)
    r2 = g.counter(fn2)
    r2.attacks(r1)
    r1.with_spec(TestSpec("admin"))
    assert "#" in g.defeats[0].to


def test_same_func_different_spec():
    fn = lambda ctx, specs: (True, None)
    g = Graph("test")
    g.rule(fn).with_spec(TestSpec("admin"))
    g.rule(fn).with_spec(TestSpec("editor"))
    assert len(g.rules) == 2
    assert g.rules[0].name != g.rules[1].name


def test_qualifier_validates_range():
    fn1 = lambda ctx, specs: (True, None)
    fn2 = lambda ctx, specs: (True, None)
    g = Graph("test")
    with pytest.raises(ValueError, match="qualifier"):
        g.rule(fn1).qualifier(1.5)
    with pytest.raises(ValueError, match="qualifier"):
        g.rule(fn2).qualifier(-0.1)
