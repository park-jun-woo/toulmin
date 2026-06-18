import unittest

from rulecat import Graph, MapContext
from rulecat.graph import run_max_depth
from rulecat.types import EvalOption


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
        a.run_on(lambda node, t: seen.append(node.verdict))
        # rule without any handler -> run_on is None -> skipped
        g.rule(lambda ctx, specs: (True, None))
        res = g.run(MapContext())
        self.assertEqual(len(seen), 1)
        self.assertGreater(seen[0], 0)
        self.assertEqual(len(res.trace.all()), 2)

    def test_handler_exception_wrapped_with_cause(self):
        g = Graph("t")
        cause = RuntimeError("boom")

        def bad(node, t):
            raise cause

        g.rule(lambda ctx, specs: (True, None)).run_on(bad)
        with self.assertRaises(RuntimeError) as cm:
            g.run(MapContext())
        self.assertIs(cm.exception.__cause__, cause)
        self.assertIn("boom", str(cm.exception))
        self.assertIn("run_on", str(cm.exception))

    def test_only_active_node_fires(self):
        # warrant attacked by an active counter -> warrant verdict 0.0 (defeated),
        # so only the counter (active, verdict>0) fires.
        g = Graph("t")
        fired = []
        w = g.rule(lambda ctx, specs: (True, None))
        c = g.counter(lambda ctx, specs: (True, None))
        c.attacks(w)
        w.run_on(lambda node, t: fired.append(("w", node.verdict)))
        c.run_on(lambda node, t: fired.append(("c", node.verdict)))
        g.run(MapContext())
        self.assertEqual([name for name, _ in fired], ["c"])

    def test_trace_exposes_verdict_for_other_nodes(self):
        # gradient/verdict query from a handler via the trace list.
        g = Graph("t")
        from rulecat.short_name import short_name
        branch = []
        w = g.rule(lambda ctx, specs: (True, None))
        c = g.counter(lambda ctx, specs: (True, None))
        c.attacks(w)
        w_name = short_name(w.id)

        def h(node, t):
            other = next((e for e in t.all() if e.name == w_name), None)
            branch.append("hard" if (other and other.verdict >= 0.5) else "soft")

        c.run_on(h)
        g.run(MapContext())
        # warrant attacked -> verdict 0.0 -> soft
        self.assertEqual(branch, ["soft"])


class RunCycleTest(unittest.TestCase):
    def test_run_raises_on_run_cycle(self):
        def a_rule(ctx, specs):
            return (True, None)

        g = Graph("cyc")
        g.rule(a_rule).run(g)  # self run-cycle
        with self.assertRaises(RuntimeError) as cm:
            g.run(MapContext())
        self.assertIn("run cycle", str(cm.exception))


class RunDepthTest(unittest.TestCase):
    def test_depth_guard_raises(self):
        g = Graph("deep")
        g.rule(lambda ctx, specs: (True, None))
        opt = EvalOption(method=0, trace=False, duration=False)
        with self.assertRaises(RuntimeError) as cm:
            g._run_depth(MapContext(), opt, run_max_depth + 1)
        self.assertIn("depth exceeded", str(cm.exception))

    def test_active_node_runs_sub_graph(self):
        def parent_rule(ctx, specs):
            return (True, None)

        def sub_rule(ctx, specs):
            return (True, None)

        sub = Graph("sub-run")
        sub.rule(sub_rule).run_on(lambda node, t: t.ctx().set("ran", True))

        parent = Graph("parent-run")
        parent.rule(parent_rule).run(sub)

        ctx = MapContext()
        parent.run(ctx)
        self.assertTrue(ctx.get("ran"))

    def test_inactive_node_with_run_graph_skips_sub(self):
        def parent_rule(ctx, specs):
            return (False, None)

        def sub_rule(ctx, specs):
            return (True, None)

        sub = Graph("sub-skip")
        sub.rule(sub_rule).run_on(lambda node, t: t.ctx().set("ran", True))

        parent = Graph("parent-skip")
        parent.rule(parent_rule).run(sub)

        ctx = MapContext()
        parent.run(ctx)
        self.assertIsNone(ctx.get("ran"))

    def test_sub_run_error_is_wrapped(self):
        cause = RuntimeError("boom")

        def parent_rule(ctx, specs):
            return (True, None)

        def sub_rule(ctx, specs):
            return (True, None)

        def bad(node, t):
            raise cause

        sub = Graph("sub-err")
        sub.rule(sub_rule).run_on(bad)

        parent = Graph("parent-err")
        parent.rule(parent_rule).run(sub)

        with self.assertRaises(RuntimeError) as cm:
            parent.run(MapContext())
        msg = str(cm.exception)
        self.assertTrue(msg.startswith('run "'))
        self.assertIn("sub-err", msg)


class EvaluateInternalTest(unittest.TestCase):
    def test_none_ctx_raises(self):
        g = Graph("t")
        g.rule(lambda ctx, specs: (True, None))
        with self.assertRaises(ValueError):
            g.evaluate(None)

    def test_duplicate_target_and_attackers(self):
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
