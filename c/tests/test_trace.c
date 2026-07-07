#include "test.h"

#include <rulecat/rulecat.h>

#include <stdlib.h>

/* 각 규칙 함수는 별도 심볼 — 정체성 (fn,ud) 규약. */
static rc_status tw_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }
static rc_status tr_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }
static rc_status tw2_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }

/* w←r, trace=true → 결과 1건, trace 2엔트리(role rule/counter), verdict는 non-trace와 동일. */
static void test_trace_entries(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w = rc_graph_rule(g, tw_fn, NULL);
    rc_rule *r = rc_graph_counter(g, tr_fn, NULL);
    ASSERT_TRUE(rc_rule_attacks(r, w) == RC_OK);

    rc_eval_option opt = { RC_MATRIX, true, false };
    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, &opt, &res) == RC_OK);
    size_t n; const rc_eval_result *items = rc_results_items(res, &n);
    ASSERT_TRUE(n == 1);
    ASSERT_NEAR(items[0].verdict, 0.0, 1e-9);   /* trace 여부와 무관하게 동일 */
    ASSERT_TRUE(items[0].n_trace == 2);
    ASSERT_NOT_NULL(items[0].trace);

    /* role은 등록이 아닌 구조 기반: warrant="rule", 공격자="counter". */
    const rc_trace_entry *rule_e = NULL, *cnt_e = NULL;
    for (size_t k = 0; k < items[0].n_trace; k++) {
        const rc_trace_entry *e = &items[0].trace[k];
        if (strcmp(e->role, "rule") == 0) rule_e = e;
        else if (strcmp(e->role, "counter") == 0) cnt_e = e;
    }
    ASSERT_NOT_NULL(rule_e);
    ASSERT_NOT_NULL(cnt_e);
    if (rule_e && cnt_e) {
        ASSERT_TRUE(rule_e->activated);
        ASSERT_TRUE(cnt_e->activated);
        ASSERT_TRUE(rule_e->n_specs == 0);
        ASSERT_NEAR(rule_e->qualifier, 1.0, 1e-9);
        ASSERT_NULL(rule_e->ground);        /* Evaluate 경로는 ground 미기록 */
        ASSERT_NEAR(rule_e->duration_ms, 0.0, 1e-9);  /* duration off */
    }
    /* 방문 순서: warrant가 먼저 기록됨. */
    ASSERT_STR_EQ(items[0].trace[0].role, "rule");
    ASSERT_STR_EQ(items[0].trace[1].role, "counter");

    rc_results_destroy(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* duration=true → trace 자동 활성 + duration_ms >= 0. */
static void test_trace_duration(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_graph_rule(g, tw_fn, NULL);

    rc_eval_option opt = { RC_MATRIX, false, true };   /* duration → trace 함의 */
    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, &opt, &res) == RC_OK);
    size_t n; const rc_eval_result *items = rc_results_items(res, &n);
    ASSERT_TRUE(n == 1);
    ASSERT_NOT_NULL(items[0].trace);            /* duration이 trace를 켬 */
    ASSERT_TRUE(items[0].n_trace == 1);
    ASSERT_TRUE(items[0].trace[0].duration_ms >= 0.0);

    rc_results_destroy(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* trace 모드 + evidence dtor: 공유 규칙이 여러 워런트 패스에 같은 포인터를
   반환해도 dtor는 distinct 포인터당 정확히 1회 (이중 해제/누수 없음). */
static void *g_shared_ev = NULL;
static int   g_dtor_calls = 0;
static rc_status shared_ev_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud,
                              bool *act, void **ev) {
    (void)c; (void)s; (void)n; (void)ud;
    if (!g_shared_ev) g_shared_ev = malloc(16);
    *act = true; *ev = g_shared_ev;
    return RC_OK;
}
static void count_dtor(void *p) { free(p); g_dtor_calls++; }

static void test_trace_shared_evidence(void) {
    g_shared_ev = NULL;
    g_dtor_calls = 0;

    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_graph_set_evidence_destructor(g, count_dtor);
    rc_rule *w1 = rc_graph_rule(g, tw_fn, NULL);
    rc_rule *w2 = rc_graph_rule(g, tw2_fn, NULL);
    rc_rule *r  = rc_graph_counter(g, shared_ev_fn, NULL);   /* 두 워런트를 모두 공격 */
    ASSERT_TRUE(rc_rule_attacks(r, w1) == RC_OK);
    ASSERT_TRUE(rc_rule_attacks(r, w2) == RC_OK);

    rc_eval_option opt = { RC_MATRIX, true, false };
    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, &opt, &res) == RC_OK);
    size_t n; rc_results_items(res, &n);
    ASSERT_TRUE(n == 2);   /* 두 워런트 → 공유 규칙은 각 패스에서 재실행 */

    rc_results_destroy(res);
    ASSERT_TRUE(g_dtor_calls == 1);   /* 공유 포인터는 정확히 1회 해제 */

    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

void test_trace(void) {
    test_trace_entries();
    test_trace_duration();
    test_trace_shared_evidence();
}
