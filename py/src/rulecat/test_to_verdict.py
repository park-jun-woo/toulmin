import pytest

from rulecat.to_verdict import to_verdict


def test_to_verdict_bool_before_int():
    assert to_verdict(True) == 1.0
    assert to_verdict(False) == -1.0
    assert to_verdict(0) == 0.0
    assert isinstance(to_verdict(1), float)


@pytest.mark.parametrize("v", [-1.0, -0.25, 0.0, 0.5, 1.0])
def test_to_verdict_in_range(v):
    assert to_verdict(v) == v


@pytest.mark.parametrize("v,msg", [
    (1.0000001, "rule verdict 1.0000001 out of range [-1, 1]"),
    (-3, "rule verdict -3 out of range [-1, 1]"),
    (float("nan"), "rule verdict NaN out of range [-1, 1]"),
])
def test_to_verdict_out_of_range(v, msg):
    with pytest.raises(ValueError) as ei:
        to_verdict(v)
    assert str(ei.value) == msg


@pytest.mark.parametrize("v", [None, "1", [1.0]])
def test_to_verdict_type_error(v):
    with pytest.raises(TypeError, match="rule verdict must be bool or float"):
        to_verdict(v)
