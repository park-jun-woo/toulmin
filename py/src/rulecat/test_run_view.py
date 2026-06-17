import unittest

from rulecat.run_view import _RunView, _build_attacker_events
from rulecat.defeat_edge import DefeatEdge
from rulecat.types import NodeEvent, NodeEventType


class RunViewTest(unittest.TestCase):
    def _view(self):
        a = NodeEvent(name="a", type=NodeEventType.ACTIVE)
        b = NodeEvent(name="b", type=NodeEventType.DEFEATED)
        order = [a, b]
        by_name = {"a": a, "b": b}
        attackers = {"a": [b]}
        return _RunView(order, by_name, attackers), a, b

    def test_all_returns_copy(self):
        view, a, b = self._view()
        self.assertEqual(view.all(), [a, b])
        got = view.all()
        got.append("x")
        self.assertEqual(len(view.all()), 2)  # mutation of copy does not leak

    def test_get_found_and_missing(self):
        view, a, b = self._view()
        self.assertIs(view.get("a"), a)
        self.assertIsNone(view.get("nope"))

    def test_attackers_found_missing_and_copy(self):
        view, a, b = self._view()
        self.assertEqual(view.attackers("a"), [b])
        self.assertEqual(view.attackers("nope"), [])
        got = view.attackers("a")
        got.append("x")
        self.assertEqual(len(view.attackers("a")), 1)  # copy


class BuildAttackerEventsTest(unittest.TestCase):
    def test_groups_present_and_skips_missing(self):
        a = NodeEvent(name="a")
        b = NodeEvent(name="b")
        by_name = {"a": a, "b": b}
        defeats = [
            DefeatEdge(from_="b", to="a"),      # ev present -> appended
            DefeatEdge(from_="ghost", to="a"),  # ev is None -> continue
        ]
        result = _build_attacker_events(defeats, by_name)
        self.assertEqual(result, {"a": [b]})

    def test_empty_defeats(self):
        self.assertEqual(_build_attacker_events([], {}), {})


if __name__ == "__main__":
    unittest.main()
