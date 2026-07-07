#include "test.h"

#include <rulecat/rulecat.h>

#include <stdbool.h>
#include <stdlib.h>

/* 각 규칙 함수는 별도 심볼 — 정체성 (fn,ud) 규약. */
static rc_status ra_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }
static rc_status rb_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }
static rc_status rc_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true; *ev = NULL; return RC_OK; }

/* 핸들러 클로저 대체: user_data로 bool 플래그 주소 전달. */
static rc_status mark_ran(const rc_trace_entry *self, const rc_trace *t, void *ud)
{ (void)self; (void)t; *(bool *)ud = true; return RC_OK; }

/* 항상 실패하는 핸들러. */
static rc_status fail_handler(const rc_trace_entry *self, const rc_trace *t, void *ud)
{ (void)self; (void)t; (void)ud; return RC_ERR_RULE; }

/* w(Active) + w2(counter에 패배: verdict 0) — w의 핸들러만 발화한다. */
static void test_run_on_active_only(void) {
    rc_graph *g = rc_graph_create("g");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w  = rc_graph_rule(g, ra_fn, NULL);
    rc_rule *w2 = rc_graph_rule(g, rb_fn, NULL);
    rc_rule *c  = rc_graph_counter(g, rc_fn, NULL);
    ASSERT_TRUE(rc_rule_attacks(c, w2) == RC_OK);   /* w2 verdict → 0 */

    bool w_ran = false, w2_ran = false;
    rc_rule_run_on(w, mark_ran, &w_ran);
    rc_rule_run_on(w2, mark_ran, &w2_ran);

    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_OK);
    ASSERT_TRUE(w_ran);
    ASSERT_FALSE(w2_ran);   /* w2 defeated → not Active → 미발화 */

    rc_run_result_destroy(rr);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* trace: 전 노드 등록 순서(inactive 포함), ground == ctx(동일 포인터), verdict == 평가값. */
static void test_run_trace_fields(void) {
    rc_graph *g = rc_graph_create("g");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w  = rc_graph_rule(g, ra_fn, NULL);
    rc_rule *w2 = rc_graph_rule(g, rb_fn, NULL);
    rc_rule *c  = rc_graph_counter(g, rc_fn, NULL);
    ASSERT_TRUE(rc_rule_attacks(c, w2) == RC_OK);

    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_OK);
    const rc_trace *tr = rc_run_result_trace(rr);

    size_t n; const rc_trace_entry *all = rc_trace_all(tr, &n);
    ASSERT_TRUE(n == 3);                     /* 전 노드 — inactive w2 포함 */
    ASSERT_PTR_EQ(rc_trace_ctx(tr), ctx);
    for (size_t i = 0; i < n; i++) {
        ASSERT_PTR_EQ(all[i].ground, ctx);   /* 전 entry 동일 ctx */
        ASSERT_TRUE(all[i].duration_ms == 0.0);
    }

    /* rc_trace_get: no-spec 규칙은 full name == short name → rc_rule_id로 조회 */
    const rc_trace_entry *we = rc_trace_get(tr, rc_rule_id(w));
    ASSERT_NOT_NULL(we);
    if (we) ASSERT_NEAR(we->verdict, 1.0, 1e-9);

    const rc_trace_entry *w2e = rc_trace_get(tr, rc_rule_id(w2));
    ASSERT_NOT_NULL(w2e);
    if (w2e) {
        ASSERT_TRUE(w2e->activated);         /* 실행은 되었으나 */
        ASSERT_NEAR(w2e->verdict, 0.0, 1e-9);/* 패배 → verdict 0 → not Active */
    }
    ASSERT_NULL(rc_trace_get(tr, "no-such-node"));

    rc_run_result_destroy(rr);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* 실패 핸들러 → RC_ERR_HANDLER, *out은 유효(trace 검시 가능), 메시지에 노드 이름. */
static void test_run_handler_error(void) {
    rc_graph *g = rc_graph_create("g");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w = rc_graph_rule(g, ra_fn, NULL);
    rc_rule_run_on(w, fail_handler, NULL);

    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_ERR_HANDLER);
    ASSERT_NOT_NULL(rr);                     /* 부분 결과 반환 — destroy 필요 */
    ASSERT_TRUE(strstr(rc_graph_last_error(g), "handler failed") != NULL);
    if (rr) {
        const rc_trace *tr = rc_run_result_trace(rr);
        size_t n; rc_trace_all(tr, &n);
        ASSERT_TRUE(n == 1);                 /* trace 검시 가능 */
    }

    rc_run_result_destroy(rr);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* run_on이 등록된 그래프에서 rc_graph_evaluate → 핸들러 미발화(Evaluate는 순수). */
static void test_evaluate_pure(void) {
    rc_graph *g = rc_graph_create("g");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w = rc_graph_rule(g, ra_fn, NULL);
    bool ran = false;
    rc_rule_run_on(w, mark_ran, &ran);

    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_OK);
    ASSERT_FALSE(ran);                       /* Evaluate는 핸들러를 발화하지 않음 */
    size_t n; rc_results_items(res, &n);
    ASSERT_TRUE(n == 1);

    rc_results_destroy(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* 66개 그래프 선형 사슬(순환 없음) → 깊이 백스톱 RC_ERR_RUN_DEPTH. */
static void test_run_depth_guard(void) {
    enum { N = 66 };
    rc_graph *gs[N];
    for (int i = 0; i < N; i++) gs[i] = rc_graph_create("chain");
    for (int i = 0; i < N; i++) {
        rc_rule *r = rc_graph_rule(gs[i], ra_fn, NULL);   /* 항상 Active */
        if (i + 1 < N) rc_rule_run(r, gs[i + 1]);         /* 다음 그래프로 실행 간선 */
    }

    rc_context *ctx = rc_map_context_create();
    rc_run_result *rr = NULL;
    ASSERT_TRUE(rc_graph_run(gs[0], ctx, NULL, &rr) == RC_ERR_RUN_DEPTH);
    if (rr) rc_run_result_destroy(rr);       /* 최상위 rr은 유효 — 정리 */

    rc_map_context_destroy(ctx);
    for (int i = 0; i < N; i++) rc_graph_destroy(gs[i]);  /* 하위는 차용 — 개별 해제 */
}

void test_run(void) {
    test_run_on_active_only();
    test_run_trace_fields();
    test_run_handler_error();
    test_evaluate_pure();
    test_run_depth_guard();
}
