import pytest

from rulecat.attack_force import attack_force


@pytest.mark.parametrize("f", [0, 0.0, 0.3, 1, 1.0])
def test_attack_force_valid(f):
    got = attack_force(f)
    assert got == float(f) and isinstance(got, float)


@pytest.mark.parametrize("f", [-0.01, 1.01, float("nan"), float("-inf")])
def test_attack_force_invalid(f):
    with pytest.raises(ValueError, match=r"^toulmin: force must be between 0.0 and 1.0$"):
        attack_force(f)


@pytest.mark.parametrize("f", [True, None, "0.5"])
def test_attack_force_type_error(f):
    with pytest.raises(TypeError, match="force must be a float"):
        attack_force(f)
