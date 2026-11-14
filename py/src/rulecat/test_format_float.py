import pytest

from rulecat.format_float import format_float


@pytest.mark.parametrize("v,want", [
    (float("nan"), "NaN"), (float("inf"), "+Inf"), (float("-inf"), "-Inf"),
    (2.0, "2"), (-2.0, "-2"), (1.5, "1.5"), (1e21, "1e+21"), (0.0, "0"),
])
def test_format_float(v, want):
    assert format_float(v) == want
