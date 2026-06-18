import pytest
from rulecat import Graph, MapContext


def _active(ctx, specs):
    return (True, "active")


# --- 1. basic Active -> sub Run (ctx side effect) ---


def test_basic_active_runs_sub_graph():
    def parent_rule(ctx, specs):
        return (True, "parent")

    def sub_rule(ctx, specs):
        return (True, "sub")

    sub = Graph("sub")
    sub.rule(sub_rule).run_on(lambda node, t: t.ctx().set("sub_ran", True))

    parent = Graph("parent")
    parent.rule(parent_rule).run(sub)

    ctx = MapContext()
    parent.run(ctx)
    assert ctx.get("sub_ran") is True


# --- 2. Active-only: Defeated/Inactive node does not run sub ---


def test_defeated_node_does_not_run_sub():
    def warrant(ctx, specs):
        return (True, "w")

    def attacker(ctx, specs):
        return (True, "c")

    def sub_rule(ctx, specs):
        return (True, "sub")

    sub = Graph("sub-defeated")
    sub.rule(sub_rule).run_on(lambda node, t: t.ctx().set("sub_ran", True))

    parent = Graph("parent-defeated")
    w = parent.rule(warrant)
    c = parent.counter(attacker)
    c.attacks(w)
    w.run(sub)  # w is DEFEATED -> sub must NOT run

    ctx = MapContext()
    parent.run(ctx)
    assert ctx.get("sub_ran") is None


def test_inactive_node_does_not_run_sub():
    def inactive_rule(ctx, specs):
        return (False, None)

    def sub_rule(ctx, specs):
        return (True, "sub")

    sub = Graph("sub-inactive")
    sub.rule(sub_rule).run_on(lambda node, t: t.ctx().set("sub_ran", True))

    parent = Graph("parent-inactive")
    parent.rule(inactive_rule).run(sub)

    ctx = MapContext()
    parent.run(ctx)
    assert ctx.get("sub_ran") is None


# --- 3. handler-then-sub ordering ---


def test_handler_fires_before_sub_run():
    order = []

    def parent_rule(ctx, specs):
        return (True, "p")

    def sub_rule(ctx, specs):
        return (True, "s")

    sub = Graph("sub-order")
    sub.rule(sub_rule).run_on(lambda node, t: order.append("sub"))

    parent = Graph("parent-order")
    parent.rule(parent_rule).run_on(
        lambda node, t: order.append("handler")
    ).run(sub)

    parent.run(MapContext())
    assert order == ["handler", "sub"]


# --- 4. ctx flows down ---


def test_ctx_flows_down_to_sub():
    captured = {}

    def parent_rule(ctx, specs):
        return (True, "p")

    def sub_rule(ctx, specs):
        return (True, "s")

    sub = Graph("sub-ctx")
    sub.rule(sub_rule).run_on(
        lambda node, t: captured.__setitem__("seen", t.ctx().get("from_parent"))
    )

    parent = Graph("parent-ctx")
    parent.rule(parent_rule).run_on(
        lambda node, t: t.ctx().set("from_parent", 42)
    ).run(sub)

    parent.run(MapContext())
    assert captured.get("seen") == 42


# --- 5. cycle A -> B -> A -> run() raises ---


def test_run_cycle_raises():
    def a_rule(ctx, specs):
        return (True, "a")

    def b_rule(ctx, specs):
        return (True, "b")

    a = Graph("A")
    b = Graph("B")
    a.rule(a_rule).run(b)
    b.rule(b_rule).run(a)

    with pytest.raises(RuntimeError) as exc:
        a.run(MapContext())
    assert "run cycle" in str(exc.value)


# --- 6. depth guard ---


def test_depth_guard_raises_on_long_chain():
    from rulecat.graph import run_max_depth

    def chain_rule(ctx, specs):
        return (True, "chain")

    graphs = [Graph(f"chain-{i}") for i in range(run_max_depth + 3)]
    for i, g in enumerate(graphs):
        r = g.rule(chain_rule)
        if i + 1 < len(graphs):
            r.run(graphs[i + 1])

    with pytest.raises(RuntimeError) as exc:
        graphs[0].run(MapContext())
    assert "depth exceeded" in str(exc.value)


# --- 7. verdict isolation ---


def test_sub_verdict_does_not_leak_into_parent_results():
    def parent_rule(ctx, specs):
        return (True, "p")

    def sub_rule(ctx, specs):
        return (True, "s")

    sub = Graph("sub-verdict")
    sub.rule(sub_rule)

    parent = Graph("parent-verdict")
    parent.rule(parent_rule).run(sub)

    sub_result = sub.run(MapContext())
    result = parent.run(MapContext())
    # parent results contain only the parent's own warrant — sub verdict does
    # not leak upward, even though the sub graph also produced a result.
    assert len(result.results) == 1
    assert len(sub_result.results) == 1


# --- 8. diamond DAG legal ---


def test_diamond_dag_runs_shared_sub_twice():
    counter = {"n": 0}

    def left_rule(ctx, specs):
        return (True, "l")

    def right_rule(ctx, specs):
        return (True, "r")

    def shared_rule(ctx, specs):
        return (True, "shared")

    shared = Graph("shared")
    shared.rule(shared_rule).run_on(
        lambda node, t: counter.__setitem__("n", counter["n"] + 1)
    )

    root = Graph("root-diamond")
    root.rule(left_rule).run(shared)
    root.rule(right_rule).run(shared)

    root.run(MapContext())  # must not raise (legal DAG)
    assert counter["n"] == 2


# --- 9. wrapped run "..." -> "..." propagation ---


def test_sub_handler_error_wrapped_and_propagated():
    cause = RuntimeError("boom")

    def parent_rule(ctx, specs):
        return (True, "p")

    def sub_rule(ctx, specs):
        return (True, "s")

    def bad(node, t):
        raise cause

    sub = Graph("sub-error")
    sub.rule(sub_rule).run_on(bad)

    parent = Graph("parent-error")
    parent.rule(parent_rule).run(sub)

    with pytest.raises(RuntimeError) as exc:
        parent.run(MapContext())

    msg = str(exc.value)
    assert msg.startswith('run "')
    assert "→" in msg
    assert "sub-error" in msg
    # outer wraps the sub's handler-wrap, which wraps the original cause
    assert isinstance(exc.value.__cause__, RuntimeError)
    assert exc.value.__cause__.__cause__ is cause
