#include "test.h"

#include <rulecat/rulecat.h>

#include <stdbool.h>
#include <stdlib.h>

static rc_status pw_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }
static rc_status sub_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }
static rc_status g1_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }
static rc_status g2_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }

static rc_status mark_ran(const rc_trace_entry *self, const rc_trace *t, void *ud)
{ (void)self; (void)t; *(bool *)ud = true; return RC_OK; }

/* ctx가 아래로 흐르는지 확인: 기대 ctx와 rc_trace_ctx(t)가 같은 포인터여야 한다. */
static rc_context *g_expected_ctx = NULL;
static bool        g_ctx_ok = false;
static rc_status check_ctx(const rc_trace_entry *self, const rc_trace *t, void *ud)
{ (void)self; (void)ud; g_ctx_ok = (rc_trace_ctx(t) == g_expected_ctx); return RC_OK; }

/* 발화 순서 기록: 핸들러 먼저(1) → 하위 그래프 나중(2). */
static int g_seq = 0;
static int g_parent_seq = 0, g_sub_seq = 0;
static rc_status rec_parent(const rc_trace_entry *self, const rc_trace *t, void *ud)
{ (void)self; (void)t; (void)ud; g_parent_seq = ++g_seq; return RC_OK; }
static rc_status rec_sub(const rc_trace_entry *self, const rc_trace *t, void *ud)
{ (void)self; (void)t; (void)ud; g_sub_seq = ++g_seq; return RC_OK; }

/* 부모 w Active → rc_rule_run(w, sub) → sub의 핸들러 발화 + ctx 아래로 전달. */
static void test_run_compose_basic(void) {
    rc_context *ctx = rc_map_context_create();
    g_expected_ctx = ctx;
    g_ctx_ok = false;

    rc_graph *sub = rc_graph_create("sub");
    bool sub_ran = false;
    rc_rule *sr = rc_graph_rule(sub, sub_fn, NULL);
    rc_rule_run_on(sr, mark_ran, &sub_ran);
    /* 두 번째 핸들러로 ctx 흐름 확인 (같은 노드에 하나만 → 별 규칙 필요) */

    rc_graph *g = rc_graph_create("g");
    rc_rule *w = rc_graph_rule(g, pw_fn, NULL);
    rc_rule_run_on(w, check_ctx, NULL);   /* 부모 노드에서 ctx 확인 */
    rc_rule_run(w, sub);

    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_OK);
    ASSERT_TRUE(sub_ran);                 /* 하위 그래프 핸들러 발화 */
    ASSERT_TRUE(g_ctx_ok);                /* ctx가 그대로 흐름 */

    rc_run_result_destroy(rr);
    rc_graph_destroy(g);
    rc_graph_destroy(sub);                /* 하위는 차용 — 개별 해제 */
    rc_map_context_destroy(ctx);
}

/* 하위 그래프가 Run되어도 부모 results의 verdict는 1.0 그대로(격리). */
static void test_run_compose_verdict_isolation(void) {
    rc_context *ctx = rc_map_context_create();
    rc_graph *sub = rc_graph_create("sub");
    bool sub_ran = false;
    rc_rule_run_on(rc_graph_rule(sub, sub_fn, NULL), mark_ran, &sub_ran);

    rc_graph *g = rc_graph_create("g");
    rc_rule *w = rc_graph_rule(g, pw_fn, NULL);
    rc_rule_run(w, sub);

    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_OK);
    size_t n; const rc_eval_result *items = rc_results_items(rc_run_result_results(rr), &n);
    ASSERT_TRUE(n == 1);                   /* 부모 워런트 1건 — 하위 결과는 안 샘 */
    if (n == 1) ASSERT_NEAR(items[0].verdict, 1.0, 1e-9);
    ASSERT_TRUE(sub_ran);

    rc_run_result_destroy(rr);
    rc_graph_destroy(g);
    rc_graph_destroy(sub);
    rc_map_context_destroy(ctx);
}

/* 핸들러 먼저, 하위 그래프 나중 (노드당 (a) run_on → (b) run_graph). */
static void test_run_compose_handler_then_sub(void) {
    rc_context *ctx = rc_map_context_create();
    g_seq = 0; g_parent_seq = 0; g_sub_seq = 0;

    rc_graph *sub = rc_graph_create("sub");
    rc_rule_run_on(rc_graph_rule(sub, sub_fn, NULL), rec_sub, NULL);

    rc_graph *g = rc_graph_create("g");
    rc_rule *w = rc_graph_rule(g, pw_fn, NULL);
    rc_rule_run_on(w, rec_parent, NULL);
    rc_rule_run(w, sub);

    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_OK);
    ASSERT_TRUE(g_parent_seq == 1);       /* 핸들러 먼저 */
    ASSERT_TRUE(g_sub_seq == 2);          /* 하위 그래프 나중 */

    rc_run_result_destroy(rr);
    rc_graph_destroy(g);
    rc_graph_destroy(sub);
    rc_map_context_destroy(ctx);
}

/* g1 노드 →run→ g2, g2 노드 →run→ g1 → RC_ERR_RUN_CYCLE. */
static void test_run_cycle_detection(void) {
    rc_context *ctx = rc_map_context_create();
    rc_graph *g1 = rc_graph_create("g1");
    rc_graph *g2 = rc_graph_create("g2");
    rc_rule_run(rc_graph_rule(g1, g1_fn, NULL), g2);
    rc_rule_run(rc_graph_rule(g2, g2_fn, NULL), g1);

    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g1, ctx, NULL, &rr) == RC_ERR_RUN_CYCLE);
    ASSERT_NULL(rr);                      /* 디스패치 이전 실패 → *out NULL */
    ASSERT_TRUE(strstr(rc_graph_last_error(g1), "run cycle") != NULL);

    rc_graph_destroy(g1);
    rc_graph_destroy(g2);
    rc_map_context_destroy(ctx);
}

/* rc_rule_run(w, NULL) → 지연 빌더 오류 → rc_graph_run은 RC_ERR_BUILDER. */
static void test_run_nil_subgraph(void) {
    rc_context *ctx = rc_map_context_create();
    rc_graph *g = rc_graph_create("g");
    rc_rule *w = rc_graph_rule(g, pw_fn, NULL);
    rc_rule_run(w, NULL);                 /* 지연 오류 */

    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_ERR_BUILDER);
    ASSERT_NULL(rr);

    rc_graph_destroy(g);
    rc_map_context_destroy(ctx);
}

void test_run_compose(void) {
    test_run_compose_basic();
    test_run_compose_verdict_isolation();
    test_run_compose_handler_then_sub();
    test_run_cycle_detection();
    test_run_nil_subgraph();
}
