#include "test.h"

#include <rulecat/rulecat.h>

/* 규칙 함수는 각각 별도 심볼 — 정체성 (fn,ud) 규약. */
static rc_status w_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true;  *ev = NULL; return RC_OK; }
static rc_status r_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true;  *ev = NULL; return RC_OK; }
static rc_status d_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true;  *ev = NULL; return RC_OK; }
static rc_status off_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = false; *ev = NULL; return RC_OK; }
static rc_status fail_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; (void)act;    *ev = NULL; return RC_ERR_RULE; }

/* 단일 warrant → verdict 1.0 */
static void test_single_warrant(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_graph_rule(g, w_fn, NULL);

    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_OK);
    size_t n; const rc_eval_result *items = rc_results_items(res, &n);
    ASSERT_TRUE(n == 1);
    ASSERT_NEAR(items[0].verdict, 1.0, 1e-9);

    rc_results_destroy(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* warrant + counter → verdict 0.0 */
static void test_warrant_counter(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w = rc_graph_rule(g, w_fn, NULL);
    rc_rule *r = rc_graph_counter(g, r_fn, NULL);
    ASSERT_TRUE(rc_rule_attacks(r, w) == RC_OK);

    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_OK);
    size_t n; const rc_eval_result *items = rc_results_items(res, &n);
    ASSERT_TRUE(n == 1);
    ASSERT_NEAR(items[0].verdict, 0.0, 1e-9);

    rc_results_destroy(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* 보상 사슬 w←r←d → 1/3 */
static void test_compensation_chain(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w = rc_graph_rule(g, w_fn, NULL);
    rc_rule *r = rc_graph_counter(g, r_fn, NULL);
    rc_rule *d = rc_graph_except(g, d_fn, NULL);
    ASSERT_TRUE(rc_rule_attacks(r, w) == RC_OK);
    ASSERT_TRUE(rc_rule_attacks(d, r) == RC_OK);

    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_OK);
    size_t n; const rc_eval_result *items = rc_results_items(res, &n);
    ASSERT_TRUE(n == 1);
    ASSERT_NEAR(items[0].verdict, 1.0 / 3.0, 1e-5);

    rc_results_destroy(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* 순환 → RC_ERR_CYCLE + last_error "cycle" */
static void test_cycle_detection(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_rule *a = rc_graph_rule(g, w_fn, NULL);
    rc_rule *b = rc_graph_counter(g, r_fn, NULL);
    ASSERT_TRUE(rc_rule_attacks(a, b) == RC_OK);
    ASSERT_TRUE(rc_rule_attacks(b, a) == RC_OK);

    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_ERR_CYCLE);
    ASSERT_NULL(res);
    ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "cycle"));

    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* qualifier 0.0 → verdict -1.0 */
static void test_qualifier_zero(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w = rc_graph_rule(g, w_fn, NULL);
    rc_rule_qualifier(w, 0.0);

    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_OK);
    size_t n; const rc_eval_result *items = rc_results_items(res, &n);
    ASSERT_TRUE(n == 1);
    ASSERT_NEAR(items[0].verdict, -1.0, 1e-9);

    rc_results_destroy(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* 비활성 규칙 → 결과 없음 */
static void test_inactive_rule(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_graph_rule(g, off_fn, NULL);

    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_OK);
    size_t n; rc_results_items(res, &n);
    ASSERT_TRUE(n == 0);

    rc_results_destroy(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* 지연 빌더 오류 → evaluate 진입 시 RC_ERR_BUILDER */
static void test_builder_error_deferred(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w = rc_graph_rule(g, w_fn, NULL);
    rc_rule_qualifier(w, 1.5);   /* 범위 위반 → builder_err 래치 */

    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_ERR_BUILDER);
    ASSERT_NULL(res);

    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* 규칙 함수 실패 → RC_ERR_RULE + "failed" 메시지 */
static void test_rule_failure(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_graph_rule(g, fail_fn, NULL);

    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_ERR_RULE);
    ASSERT_NULL(res);
    ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "failed"));

    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* NULL 인자 → RC_ERR_INVALID_ARG */
static void test_null_args(void) {
    rc_graph *g = rc_graph_create("test");
    rc_context *ctx = rc_map_context_create();
    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(NULL, ctx, NULL, &res) == RC_ERR_INVALID_ARG);
    ASSERT_TRUE(rc_graph_evaluate(g, NULL, NULL, &res) == RC_ERR_INVALID_ARG);
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, NULL) == RC_ERR_INVALID_ARG);
    ASSERT_NULL(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

void test_evaluate(void) {
    test_single_warrant();
    test_warrant_counter();
    test_compensation_chain();
    test_cycle_detection();
    test_qualifier_zero();
    test_inactive_rule();
    test_builder_error_deferred();
    test_rule_failure();
    test_null_args();
}
