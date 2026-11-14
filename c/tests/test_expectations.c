#include "test.h"

#include <rulecat/rulecat.h>

static rc_status ew_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true;  *ev = NULL; return RC_OK; }
static rc_status er_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true;  *ev = NULL; return RC_OK; }
static rc_status ed_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = true;  *ev = NULL; return RC_OK; }
static rc_status eoff_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; *act = false; *ev = NULL; return RC_OK; }

/* verdict 1.0 결과 그래프. */
static rc_graph *make_v1(void) {
    rc_graph *g = rc_graph_create("v1");
    rc_graph_rule(g, ew_fn, NULL);
    return g;
}
/* verdict 1/3 결과 그래프 (보상 사슬 w←r←d, 반 강도 Except).
   반박 r의 적용도 0.5 → r verdict 0 → w = 1/3. 강도 1이면 r이 완전히 게이트되어 w = 1.0. */
static rc_graph *make_v13(void) {
    rc_graph *g = rc_graph_create("v13");
    rc_rule *w = rc_graph_rule(g, ew_fn, NULL);
    rc_rule *r = rc_graph_counter(g, er_fn, NULL);
    rc_rule *d = rc_graph_except(g, ed_fn, NULL);
    rc_rule_attacks(r, w);
    rc_rule_attacks_with_force(d, r, 0.5);
    return g;
}
/* 결과 0건 그래프 (비활성 규칙). */
static rc_graph *make_empty(void) {
    rc_graph *g = rc_graph_create("empty");
    rc_graph_rule(g, eoff_fn, NULL);
    return g;
}

static rc_results *eval_of(rc_graph *g, rc_context *ctx) {
    rc_results *res = NULL;
    rc_graph_evaluate(g, ctx, NULL, &res);
    return res;
}

static void test_expectations(void) {
    rc_context *ctx = rc_map_context_create();
    char why[160];

    rc_graph *g1 = make_v1();
    rc_results *res_v1 = eval_of(g1, ctx);
    ASSERT_TRUE(rc_expect_check(res_v1, rc_verdict_above(0.5), why, sizeof why) == RC_OK);
    ASSERT_TRUE(rc_expect_check(res_v1, rc_verdict_at_most(0.5), why, sizeof why) != RC_OK);
    ASSERT_TRUE(rc_expect_check(res_v1, rc_verdict_between(0.9, 1.1), why, sizeof why) == RC_OK);
    ASSERT_TRUE(rc_expect_check(res_v1, rc_no_result(), why, sizeof why) != RC_OK);
    /* errbuf NULL 허용 */
    ASSERT_TRUE(rc_expect_check(res_v1, rc_verdict_at_most(0.5), NULL, 0) != RC_OK);

    rc_graph *g13 = make_v13();
    rc_results *res_v13 = eval_of(g13, ctx);
    ASSERT_TRUE(rc_expect_check(res_v13, rc_verdict_between(0.0, 0.5), why, sizeof why) == RC_OK);
    ASSERT_TRUE(rc_expect_check(res_v13, rc_verdict_above(0.5), why, sizeof why) != RC_OK);

    rc_graph *ge = make_empty();
    rc_results *res_empty = eval_of(ge, ctx);
    ASSERT_TRUE(rc_expect_check(res_empty, rc_no_result(), why, sizeof why) == RC_OK);
    ASSERT_TRUE(rc_expect_check(res_empty, rc_verdict_above(0.0), why, sizeof why) != RC_OK);

    rc_results_destroy(res_v1);
    rc_results_destroy(res_v13);
    rc_results_destroy(res_empty);
    rc_graph_destroy(g1);
    rc_graph_destroy(g13);
    rc_graph_destroy(ge);
    rc_map_context_destroy(ctx);
}

/* 테이블 주도: 통과 케이스 묶음 + 실패 케이스 검출. */
static void test_run_cases(void) {
    rc_graph *g = make_v1();   /* 항상 verdict 1.0 */
    char why[256];

    rc_test_case ok_cases[] = {
        { .name = "above",   .ctx = NULL, .opt = NULL, .expect = rc_verdict_above(0.9) },
        { .name = "between", .ctx = NULL, .opt = NULL, .expect = rc_verdict_between(0.5, 1.5) },
    };
    ASSERT_TRUE(rc_run_cases(g, ok_cases, 2, why, sizeof why) == RC_OK);

    rc_test_case bad_cases[] = {
        { .name = "wrong", .ctx = NULL, .opt = NULL, .expect = rc_verdict_at_most(0.0) },
    };
    why[0] = '\0';
    ASSERT_TRUE(rc_run_cases(g, bad_cases, 1, why, sizeof why) != RC_OK);
    ASSERT_NOT_NULL(strstr(why, "wrong"));   /* "[wrong] ..." 케이스명 포함 */

    rc_graph_destroy(g);
}

void test_expect_and_cases(void) {
    test_expectations();
    test_run_cases();
}
