from rulecat import find_spec


class TestSpec:
    def __init__(self, value: str):
        self.value = value
    def spec_name(self) -> str:
        return "TestSpec"
    def validate(self) -> None:
        pass


def test_find_by_name():
    specs = [TestSpec("a")]
    found = find_spec(specs, "TestSpec")
    assert found is not None
    assert found.value == "a"


def test_not_found():
    assert find_spec([], "nope") is None
