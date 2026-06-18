import unittest

from rulecat import Graph


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


if __name__ == "__main__":
    unittest.main()
