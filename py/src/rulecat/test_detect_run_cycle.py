import unittest

from rulecat import Graph
from rulecat.detect_run_cycle import detect_run_cycle


def _active(ctx, specs):
    return (True, None)


class DetectRunCycleTest(unittest.TestCase):
    def test_single_graph_no_run_edges_is_acyclic(self):
        # Rule has run_graph None -> `sub is None` continue branch.
        g = Graph("solo")
        g.rule(_active)
        self.assertIsNone(detect_run_cycle(g))

    def test_linear_chain_is_acyclic(self):
        # WHITE sub -> dfs recurse returns None.
        def child_rule(ctx, specs):
            return (True, None)

        sub = Graph("child")
        sub.rule(child_rule)

        root = Graph("root")
        root.rule(_active).run(sub)

        self.assertIsNone(detect_run_cycle(root))

    def test_direct_cycle_detected(self):
        # A -> B -> A: B sees A as GRAY -> cycle message.
        def a_rule(ctx, specs):
            return (True, None)

        def b_rule(ctx, specs):
            return (True, None)

        a = Graph("A")
        b = Graph("B")
        a.rule(a_rule).run(b)
        b.rule(b_rule).run(a)

        err = detect_run_cycle(a)
        self.assertIsNotNone(err)
        self.assertIn("run cycle", err)
        self.assertIn('"A"', err)

    def test_cycle_propagated_from_nested_dfs(self):
        # root -> mid -> mid (self loop): error bubbles up through the
        # `if err: return err` propagation branch in the outer frame.
        def root_rule(ctx, specs):
            return (True, None)

        def mid_rule_a(ctx, specs):
            return (True, None)

        def mid_rule_b(ctx, specs):
            return (True, None)

        mid = Graph("mid")
        mid.rule(mid_rule_a).run(mid)  # self-cycle

        root = Graph("root2")
        root.rule(root_rule).run(mid)

        err = detect_run_cycle(root)
        self.assertIsNotNone(err)
        self.assertIn('"mid"', err)

    def test_diamond_dag_is_acyclic(self):
        # left and right both point at shared (BLACK on 2nd visit) -> no cycle.
        def left_rule(ctx, specs):
            return (True, None)

        def right_rule(ctx, specs):
            return (True, None)

        def shared_rule(ctx, specs):
            return (True, None)

        shared = Graph("shared")
        shared.rule(shared_rule)

        root = Graph("root3")
        root.rule(left_rule).run(shared)
        root.rule(right_rule).run(shared)

        self.assertIsNone(detect_run_cycle(root))


if __name__ == "__main__":
    unittest.main()
