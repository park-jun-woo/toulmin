#include "test.h"

#include <rulecat/rulecat.h>

/* 0.3.0 배제(supersede): 배제자의 적용도 a_eff > 0일 때만 대상을 크리스프하게 배제한다.
   네 포트 공통 VIP 시나리오 — vip가 general을 배제한다. 예외 게이트는 반영하고 반박은
   반영하지 않는다. (Go run_supersedes_except_gate_test.go / run_supersedes_chain_test.go 대응) */

/* 항상 해당하는 bool 규칙. 정체성 키가 (fn, user_data)라 user_data로 규칙을 구분한다. */
static rc_status sup_true_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }

static int TAG_GENERAL, TAG_VIP, TAG_SUSPENDED, TAG_DISPUTED, TAG_X, TAG_Y, TAG_Z;

enum vip_wire { WIRE_NONE, WIRE_EXCEPT_FULL, WIRE_EXCEPT_HALF, WIRE_COUNTER };

/* general(0) ← supersede ← vip(1), 그다음 wire대로 공격자(2)를 붙인다. */
static rc_graph *build_vip(enum vip_wire wire) {
    rc_graph *g = rc_graph_create("vip");
    rc_rule *general = rc_graph_rule(g, sup_true_fn, &TAG_GENERAL);
    rc_rule *vip = rc_graph_rule(g, sup_true_fn, &TAG_VIP);
    ASSERT_TRUE(rc_rule_supersedes(vip, general) == RC_OK);
    if (wire == WIRE_EXCEPT_FULL)
        rc_rule_attacks(rc_graph_except(g, sup_true_fn, &TAG_SUSPENDED), vip);
    else if (wire == WIRE_EXCEPT_HALF)
        rc_rule_attacks_with_force(rc_graph_except(g, sup_true_fn, &TAG_SUSPENDED), vip, 0.5);
    else if (wire == WIRE_COUNTER)
        rc_rule_attacks(rc_graph_counter(g, sup_true_fn, &TAG_DISPUTED), vip);
    return g;
}

static void check_vip(enum vip_wire wire, double vip_verdict, double vip_a_eff, bool want_excluded) {
    rc_graph *g = build_vip(wire);
    rc_context *ctx = rc_map_context_create();

    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_OK);
    size_t n = 0;
    const rc_trace_entry *e = rc_trace_all(rc_run_result_trace(rr), &n);
    ASSERT_TRUE(n >= 2);
    ASSERT_TRUE(e[0].excluded == want_excluded);   /* general */
    ASSERT_NEAR(e[0].verdict, 1.0, 1e-12);         /* 보고 전용 */
    ASSERT_FALSE(e[1].excluded);                   /* vip 자신 */
    ASSERT_TRUE(e[1].activated);
    ASSERT_NEAR(e[1].verdict, vip_verdict, 1e-12);
    size_t nr = 0;
    const rc_eval_result *items = rc_results_items(rc_run_result_results(rr), &nr);
    ASSERT_TRUE(nr == 2);                           /* general, vip (등록 순서) */
    if (nr == 2) ASSERT_NEAR(items[1].applicability, vip_a_eff, 1e-12);
    rc_run_result_destroy(rr);

    rc_results *res = NULL;
    rc_eval_option opt = { RC_MATRIX, true, false };
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, &opt, &res) == RC_OK);
    items = rc_results_items(res, &nr);
    ASSERT_TRUE(nr == 2 && items[0].n_trace == 1);
    if (nr == 2 && items[0].n_trace == 1)
        ASSERT_TRUE(items[0].trace[0].excluded == want_excluded);
    rc_results_destroy(res);

    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* x ⊐ y ⊐ z 사슬. 배제된 y도 z를 배제한다. x가 예외로 완전히 막히면 x의 배제만 풀린다. */
static void check_chain(bool gate_x, bool want_y, bool want_z) {
    rc_graph *g = rc_graph_create("chain");
    rc_rule *z = rc_graph_rule(g, sup_true_fn, &TAG_Z);
    rc_rule *y = rc_graph_rule(g, sup_true_fn, &TAG_Y);
    rc_rule *x = rc_graph_rule(g, sup_true_fn, &TAG_X);
    rc_rule_supersedes(x, y);
    rc_rule_supersedes(y, z);
    if (gate_x) rc_rule_attacks(rc_graph_except(g, sup_true_fn, &TAG_SUSPENDED), x);
    rc_context *ctx = rc_map_context_create();
    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_OK);
    size_t n = 0;
    const rc_trace_entry *e = rc_trace_all(rc_run_result_trace(rr), &n);
    ASSERT_TRUE(n >= 3);
    if (n >= 3) {
        ASSERT_TRUE(e[0].excluded == want_z);
        ASSERT_TRUE(e[1].excluded == want_y);
        ASSERT_FALSE(e[2].excluded);
    }
    rc_run_result_destroy(rr);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

static void check_cycle_and_args(void) {
    rc_graph *g = rc_graph_create("cycle");
    rc_rule *a = rc_graph_rule(g, sup_true_fn, &TAG_X);
    rc_rule *b = rc_graph_rule(g, sup_true_fn, &TAG_Y);
    ASSERT_TRUE(rc_rule_supersedes(NULL, b) == RC_ERR_INVALID_ARG);
    ASSERT_TRUE(rc_rule_supersedes(a, NULL) == RC_ERR_INVALID_ARG);
    rc_rule_supersedes(a, b);
    rc_rule_supersedes(b, a);
    rc_context *ctx = rc_map_context_create();
    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_ERR_CYCLE);
    ASSERT_NULL(res);
    ASSERT_TRUE(strstr(rc_graph_last_error(g), "supersede cycle detected") != NULL);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

void test_supersedes(void) {
    check_vip(WIRE_NONE,         1.0, 1.0, true);    /* 예외 없음: 기존과 같음 */
    check_vip(WIRE_EXCEPT_FULL, -1.0, 0.0, false);   /* 예외가 완전히 막음: 배제 안 됨 */
    check_vip(WIRE_EXCEPT_HALF,  0.0, 0.5, true);    /* 예외 force 0.5: a_eff 0.5 > 0 */
    check_vip(WIRE_COUNTER,      0.0, 1.0, true);    /* 반박: 배제 유지 */
    check_chain(false, true, true);
    check_chain(true, false, true);
    check_cycle_and_args();
}
