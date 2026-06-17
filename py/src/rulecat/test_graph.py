import unittest

from rulecat import Graph, MapContext, NodeEventType
from rulecat.graph import _classify_event, _select_handler
from rulecat.rule_meta import RuleMeta


class ClassifyEventTest(unittest.TestCase):
    def test_inactive(self):
        self.assertEqual(_classify_event(False, 1.0), NodeEventType.INACTIVE)
        self.assertEqual(_classify_event(False, -1.0), NodeEventType.INACTIVE)

    def test_active_when_verdict_positive(self):
        self.assertEqual(_classify_event(True, 0.5), NodeEventType.ACTIVE)

    def test_defeated_when_verdict_not_positive(self):
        self.assertEqual(_classify_event(True, 0.0), NodeEventType.DEFEATED)
        self.assertEqual(_classify_event(True, -0.5), NodeEventType.DEFEATED)


class SelectHandlerTest(unittest.TestCase):
    def _meta(self):
        return RuleMeta(
            name="r", qualifier=1.0, strength=0,
            on_active="A", on_defeated="D", on_inactive="I",
        )

    def test_active(self):
        self.assertEqual(_select_handler(self._meta(), NodeEventType.ACTIVE), "A")

    def test_defeated(self):
        self.assertEqual(_select_handler(self._meta(), NodeEventType.DEFEATED), "D")

    def test_inactive(self):
        self.assertEqual(_select_handler(self._meta(), NodeEventType.INACTIVE), "I")


class RunTest(unittest.TestCase):
    def test_none_ctx_raises(self):
        g = Graph("t")
        g.rule(lambda ctx, specs: (True, None))
        with self.assertRaises(ValueError):
            g.run(None)

    def test_dispatch_calls_handler_and_skips_handlerless(self):
        g = Graph("t")
        seen = []
        a = g.rule(lambda ctx, specs: (True, None))
        a.on_active(lambda ctx, ev, view: seen.append(ev.type))
        # rule without any handler -> _select_handler returns None -> continue
        g.rule(lambda ctx, specs: (True, None))
        res = g.run(MapContext())
        self.assertEqual(seen, [NodeEventType.ACTIVE])
        self.assertEqual(len(res.view.all()), 2)

    def test_handler_exception_wrapped_with_cause(self):
        g = Graph("t")
        cause = RuntimeError("boom")

        def bad(ctx, ev, view):
            raise cause

        g.rule(lambda ctx, specs: (True, None)).on_active(bad)
        with self.assertRaises(RuntimeError) as cm:
            g.run(MapContext())
        self.assertIs(cm.exception.__cause__, cause)
        self.assertIn("boom", str(cm.exception))
        self.assertIn("ACTIVE", str(cm.exception))


class EvaluateInternalTest(unittest.TestCase):
    def test_none_ctx_raises(self):
        g = Graph("t")
        g.rule(lambda ctx, specs: (True, None))
        with self.assertRaises(ValueError):
            g.evaluate(None)

    def test_duplicate_target_and_attackers(self):
        # Two counters attack the same warrant: the edges dict sees the target
        # missing (True branch) then present (False branch), and the attacker
        # set is built from a non-empty attacker list.
        g = Graph("t")
        w = g.rule(lambda ctx, specs: (True, None))
        c1 = g.counter(lambda ctx, specs: (True, None))
        c2 = g.counter(lambda ctx, specs: (True, None))
        c1.attacks(w)
        c2.attacks(w)
        results = g.evaluate(MapContext())
        self.assertTrue(results)

    def test_cycle_raises(self):
        g = Graph("t")
        a = g.rule(lambda ctx, specs: (True, None))
        b = g.counter(lambda ctx, specs: (True, None))
        a.attacks(b)
        b.attacks(a)
        with self.assertRaises(RuntimeError):
            g.evaluate(MapContext())


if __name__ == "__main__":
    unittest.main()
