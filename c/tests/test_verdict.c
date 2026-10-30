#include "test.h"

#include <rulecat/rulecat.h>

#include <math.h>
#include <stdlib.h>

/* 0.2.0 의미론: verdict 규칙, bool 호환, 예외 게이트·applicability, 공격 강도, Qualifier NaN.
   (Go graph_evaluate_*_test.go, attack_force_test.go, test_qualifier_panic_nan_test.go 대응) */

/* --- 규칙 함수 ---
   정체성 키는 (fn, user_data)이므로 같은 함수를 user_data만 달리해 여러 규칙으로 쓴다. */

/* bool 규칙: user_data = const bool * */
static rc_status vb_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, bool *act, void **ev)
{ (void)c; (void)s; (void)n; *act = *(const bool *)ud; *ev = NULL; return RC_OK; }

/* verdict 규칙: user_data = const double * */
static rc_status vv_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, double *v, void **ev)
{ (void)c; (void)s; (void)n; *v = *(const double *)ud; *ev = NULL; return RC_OK; }

/* 호출되면 표시하는 verdict 규칙: user_data = struct probe * */
struct probe { double v; int calls; };
static rc_status probe_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, double *v, void **ev)
{
    (void)c; (void)s; (void)n;
    struct probe *p = ud;
    p->calls++;
    *v = p->v;
    *ev = NULL;
    return RC_OK;
}

/* evidence를 할당해 내는 verdict 규칙: user_data = const double * */
static rc_status vev_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, double *v, void **ev)
{
    (void)c; (void)s; (void)n;
    *v = *(const double *)ud;
    *ev = malloc(8);
    return RC_OK;
}
static int dtor_calls = 0;
static void count_free(void *p) { dtor_calls++; free(p); }

/* 실패 상태를 내는 verdict 규칙 */
static rc_status vfail_fn(rc_context *c, rc_spec *const *s, size_t n, void *ud, double *v, void **ev)
{ (void)c; (void)s; (void)n; (void)ud; (void)v; *ev = NULL; return RC_ERR_RULE; }

/* --- 공용 헬퍼 --- */

/* 평가 후 첫 결과를 꺼낸다. 결과 수는 *n_out. 호출자가 *res_out을 해제한다. */
static const rc_eval_result *eval_first(rc_graph *g, bool trace, rc_results **res_out, size_t *n_out) {
    rc_context *ctx = rc_map_context_create();
    rc_eval_option opt = { RC_MATRIX, trace, false };
    *res_out = NULL;
    *n_out = 0;
    rc_status st = rc_graph_evaluate(g, ctx, &opt, res_out);
    rc_map_context_destroy(ctx);
    RC_CHECK(st == RC_OK, "evaluate RC_OK");
    if (st != RC_OK) return NULL;
    const rc_eval_result *items = rc_results_items(*res_out, n_out);
    return *n_out > 0 ? items : NULL;
}

/* --- bool 호환: bool 규칙, ±1 verdict 규칙, 옛 공식이 정확히 같다 --- */

static const bool B_T = true;

/* w(q0.8) ← c1(q0.7), c2(q0.4, 비활성), c3(q0.3); c1 ← c4(q0.9) */
static rc_graph *build_compat(bool use_bool, const void *ud[5]) {
    rc_graph *g = rc_graph_create("compat");
    rc_rule *r[5];
    for (int i = 0; i < 5; i++) {
        if (i == 0)
            r[i] = use_bool ? rc_graph_rule(g, vb_fn, (void *)ud[i])
                            : rc_graph_rule_verdict(g, vv_fn, (void *)ud[i]);
        else
            r[i] = use_bool ? rc_graph_counter(g, vb_fn, (void *)ud[i])
                            : rc_graph_counter_verdict(g, vv_fn, (void *)ud[i]);
    }
    rc_rule_qualifier(r[0], 0.8);
    rc_rule_qualifier(r[1], 0.7);
    rc_rule_qualifier(r[2], 0.4);
    rc_rule_qualifier(r[3], 0.3);
    rc_rule_qualifier(r[4], 0.9);
    rc_rule_attacks(r[1], r[0]);
    rc_rule_attacks(r[2], r[0]);
    rc_rule_attacks(r[3], r[0]);
    rc_rule_attacks(r[4], r[1]);
    return g;
}

static void test_bool_compat(void) {
    /* 별개 정체성을 위해 규칙마다 다른 주소의 user_data를 쓴다. */
    static const bool bt0 = true, bt1 = true, bf2 = false, bt3 = true, bt4 = true;
    static const double dp0 = 1.0, dp1 = 1.0, dm2 = -1.0, dp3 = 1.0, dp4 = 1.0;
    const void *bud[5] = { &bt0, &bt1, &bf2, &bt3, &bt4 };
    const void *dud[5] = { &dp0, &dp1, &dm2, &dp3, &dp4 };

    /* 옛 공식: raw = q/(1+Σ(v_j+1)/2), v = 2raw−1 */
    double vc4 = 2 * 0.9 - 1;
    double vc1 = 2 * (0.7 / (1.0 + (vc4 + 1.0) / 2.0)) - 1;
    double vc3 = 2 * 0.3 - 1;
    double legacy = 2 * (0.8 / (1.0 + ((vc1 + 1.0) / 2.0 + 0.0 + (vc3 + 1.0) / 2.0))) - 1;

    for (int t = 0; t < 2; t++) {
        rc_graph *gb = build_compat(true, bud);
        rc_graph *gf = build_compat(false, dud);
        rc_results *rb = NULL, *rf = NULL;
        size_t nb, nf;
        const rc_eval_result *ib = eval_first(gb, t == 1, &rb, &nb);
        const rc_eval_result *iff = eval_first(gf, t == 1, &rf, &nf);
        ASSERT_TRUE(nb == 1 && nf == 1);
        if (ib && iff) {
            ASSERT_TRUE(ib->verdict == legacy);    /* 비트 단위 동일 */
            ASSERT_TRUE(iff->verdict == legacy);
            ASSERT_TRUE(ib->applicability == 1.0);
        }
        rc_results_destroy(rb);
        rc_results_destroy(rf);
        rc_graph_destroy(gb);
        rc_graph_destroy(gf);
    }
}

/* --- verdict 범위·NaN 오류 --- */

static void test_verdict_range(void) {
    static double bad[4];
    bad[0] = 1.5; bad[1] = -1.0001; bad[2] = NAN; bad[3] = INFINITY;
    for (int k = 0; k < 4; k++) {
        for (int t = 0; t < 2; t++) {
            rc_graph *g = rc_graph_create("range");
            rc_context *ctx = rc_map_context_create();
            ASSERT_NOT_NULL(rc_graph_rule_verdict(g, vv_fn, &bad[k]));
            rc_eval_option opt = { RC_MATRIX, t == 1, false };
            rc_results *res = NULL;
            ASSERT_TRUE(rc_graph_evaluate(g, ctx, &opt, &res) == RC_ERR_RULE);
            ASSERT_NULL(res);
            ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "out of range [-1, 1]"));
            ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "rule verdict"));

            rc_run_result *rr = NULL;
            ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_ERR_RULE);
            ASSERT_NULL(rr);
            rc_map_context_destroy(ctx);
            rc_graph_destroy(g);
        }
    }

    /* 공격자의 범위 오류도 전파된다. */
    {
        static const double one = 1.0, over = 2.0;
        rc_graph *g = rc_graph_create("range-attacker");
        rc_context *ctx = rc_map_context_create();
        rc_rule *w = rc_graph_rule_verdict(g, vv_fn, (void *)&one);
        rc_rule *c = rc_graph_counter_verdict(g, vv_fn, (void *)&over);
        rc_rule_attacks(c, w);
        rc_results *res = NULL;
        ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_ERR_RULE);
        ASSERT_NULL(res);
        rc_map_context_destroy(ctx);
        rc_graph_destroy(g);
    }

    /* 경계값 ±1은 합법. */
    {
        static const double lo = -1.0, hi = 1.0;
        rc_graph *g = rc_graph_create("bounds");
        rc_graph_rule_verdict(g, vv_fn, (void *)&hi);
        rc_graph_rule_verdict(g, vv_fn, (void *)&lo);
        rc_results *res = NULL; size_t n;
        const rc_eval_result *it = eval_first(g, false, &res, &n);
        ASSERT_TRUE(n == 1);   /* -1은 해당 없음 → 결과 없음 */
        if (it) ASSERT_TRUE(it->verdict == 1.0);
        rc_results_destroy(res);
        rc_graph_destroy(g);
    }

    /* 범위 오류 규칙이 낸 evidence는 dtor로 회수된다(누수 없음). */
    {
        static const double over = 3.0;
        rc_graph *g = rc_graph_create("range-ev");
        rc_context *ctx = rc_map_context_create();
        rc_graph_set_evidence_destructor(g, count_free);
        rc_graph_rule_verdict(g, vev_fn, (void *)&over);
        dtor_calls = 0;
        rc_results *res = NULL;
        ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_ERR_RULE);
        ASSERT_TRUE(dtor_calls == 1);
        rc_map_context_destroy(ctx);
        rc_graph_destroy(g);
    }

    /* verdict 규칙의 실패 상태는 기존과 같은 "failed" 경로. */
    {
        rc_graph *g = rc_graph_create("vfail");
        rc_context *ctx = rc_map_context_create();
        rc_graph_rule_verdict(g, vfail_fn, NULL);
        rc_results *res = NULL;
        ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_ERR_RULE);
        ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "failed"));
        rc_map_context_destroy(ctx);
        rc_graph_destroy(g);
    }
}

/* --- 부분 verdict, applicability, 지연 평가 --- */

static void test_partial_verdict(void) {
    /* v = 0은 '반쯤 해당'(해당 없음이 아님): a = 0.5 → verdict 0, 결과는 남는다. */
    {
        static const double zero = 0.0;
        rc_graph *g = rc_graph_create("half");
        rc_graph_rule_verdict(g, vv_fn, (void *)&zero);
        rc_results *res = NULL; size_t n;
        const rc_eval_result *it = eval_first(g, false, &res, &n);
        ASSERT_TRUE(n == 1);
        if (it) {
            ASSERT_TRUE(it->verdict == 0.0);
            ASSERT_TRUE(it->applicability == 0.5);
        }
        rc_results_destroy(res);
        rc_graph_destroy(g);
    }
    /* v = 0.5 → a = 0.75; q = 0.8 → raw 0.6 → verdict 0.2 */
    {
        static const double half = 0.5;
        rc_graph *g = rc_graph_create("q");
        rc_rule_qualifier(rc_graph_rule_verdict(g, vv_fn, (void *)&half), 0.8);
        rc_results *res = NULL; size_t n;
        const rc_eval_result *it = eval_first(g, false, &res, &n);
        ASSERT_TRUE(n == 1);
        if (it) {
            ASSERT_NEAR(it->verdict, 0.2, 1e-12);
            ASSERT_TRUE(it->applicability == 0.75);
        }
        rc_results_destroy(res);
        rc_graph_destroy(g);
    }
    /* v = -1이면 결과 없음 + 공격자 미평가(지연). */
    {
        static const double m1 = -1.0;
        struct probe p = { 1.0, 0 };
        rc_graph *g = rc_graph_create("lazy");
        rc_rule *w = rc_graph_rule_verdict(g, vv_fn, (void *)&m1);
        rc_rule *c = rc_graph_counter_verdict(g, probe_fn, &p);
        rc_rule_attacks(c, w);
        rc_results *res = NULL; size_t n;
        eval_first(g, false, &res, &n);
        ASSERT_TRUE(n == 0);
        ASSERT_TRUE(p.calls == 0);
        rc_results_destroy(res);
        rc_graph_destroy(g);
    }
    /* v = -0.5는 Active: 공격자를 평가한다. a = 0.25 → verdict -0.5. */
    {
        static const double mh = -0.5;
        struct probe p = { -1.0, 0 };
        rc_graph *g = rc_graph_create("barely");
        rc_rule *w = rc_graph_rule_verdict(g, vv_fn, (void *)&mh);
        rc_rule *c = rc_graph_counter_verdict(g, probe_fn, &p);
        rc_rule_attacks(c, w);
        rc_results *res = NULL; size_t n;
        const rc_eval_result *it = eval_first(g, false, &res, &n);
        ASSERT_TRUE(n == 1);
        ASSERT_TRUE(p.calls == 1);
        if (it) {
            ASSERT_TRUE(it->verdict == -0.5);
            ASSERT_TRUE(it->applicability == 0.25);
        }
        rc_results_destroy(res);
        rc_graph_destroy(g);
    }
}

/* --- 예외 게이트와 applicability --- */

static void test_except_gate(void) {
    struct gate_case {
        const char *name;
        double v;
        bool except;
        double verdict, applic;
    } cases[] = {
        { "three counters",         1.0, false, -0.5,                         1.0 },
        { "three full excepts",     1.0, true,  -1.0,                         0.0 },
        { "three partial excepts",  0.5, true,  2 * pow(0.25, 3) - 1,         pow(0.25, 3) },
        { "inactive excepts",      -1.0, true,   1.0,                         1.0 },
    };
    for (size_t k = 0; k < sizeof cases / sizeof cases[0]; k++) {
        for (int t = 0; t < 2; t++) {
            double vs[3] = { cases[k].v, cases[k].v, cases[k].v };
            static const double one = 1.0;
            rc_graph *g = rc_graph_create("gate");
            rc_rule *w = rc_graph_rule_verdict(g, vv_fn, (void *)&one);
            for (int i = 0; i < 3; i++) {
                rc_rule *a = cases[k].except ? rc_graph_except_verdict(g, vv_fn, &vs[i])
                                             : rc_graph_counter_verdict(g, vv_fn, &vs[i]);
                rc_rule_attacks(a, w);
            }
            rc_results *res = NULL; size_t n;
            const rc_eval_result *it = eval_first(g, t == 1, &res, &n);
            ASSERT_TRUE(n == 1);
            if (it) {
                RC_CHECK(fabs(it->verdict - cases[k].verdict) < 1e-12, cases[k].name);
                RC_CHECK(fabs(it->applicability - cases[k].applic) < 1e-12, cases[k].name);
            }
            rc_results_destroy(res);
            rc_graph_destroy(g);
        }
    }

    /* bool 예외도 같은 게이트: w ← c ← d(Except, 강도 0.5)
       → c 적용도 0.5 → c verdict 0 → w = 1/3. 강도 1이면 c가 -1로 게이트되어 w = 1. */
    for (int full = 0; full < 2; full++) {
        static const bool tw = true, tc = true, td = true;
        rc_graph *g = rc_graph_create("bool-except");
        rc_rule *w = rc_graph_rule(g, vb_fn, (void *)&tw);
        rc_rule *c = rc_graph_counter(g, vb_fn, (void *)&tc);
        rc_rule *d = rc_graph_except(g, vb_fn, (void *)&td);
        rc_rule_attacks(c, w);
        if (full) rc_rule_attacks(d, c);
        else      rc_rule_attacks_with_force(d, c, 0.5);
        rc_results *res = NULL; size_t n;
        const rc_eval_result *it = eval_first(g, false, &res, &n);
        ASSERT_TRUE(n == 1);
        if (it) {
            ASSERT_NEAR(it->verdict, full ? 1.0 : 1.0 / 3.0, 1e-12);
            ASSERT_TRUE(it->applicability == 1.0);   /* w 자신은 예외에 맞지 않음 */
        }
        rc_results_destroy(res);
        rc_graph_destroy(g);
    }
}

/* --- 공격 강도 --- */

static void test_attack_force(void) {
    struct force_case {
        const char *name;
        bool except;
        double force, verdict, applic;
    } cases[] = {
        { "counter F=0",    false, 0.0,  1.0,       1.0  },
        { "counter F=0.5",  false, 0.5,  1.0 / 3.0, 1.0  },   /* S = 0.5 → raw 2/3 */
        { "counter F=1",    false, 1.0,  0.0,       1.0  },
        { "except F=0",     true,  0.0,  1.0,       1.0  },
        { "except F=0.5",   true,  0.5,  0.0,       0.5  },   /* a_eff = 0.5 */
        { "except F=0.25",  true,  0.25, 0.5,       0.75 },
    };
    static const double one = 1.0, one_b = 1.0;
    for (size_t k = 0; k < sizeof cases / sizeof cases[0]; k++) {
        for (int t = 0; t < 2; t++) {
            rc_graph *g = rc_graph_create("force");
            rc_rule *w = rc_graph_rule_verdict(g, vv_fn, (void *)&one);
            rc_rule *a = cases[k].except ? rc_graph_except_verdict(g, vv_fn, (void *)&one_b)
                                         : rc_graph_counter_verdict(g, vv_fn, (void *)&one_b);
            ASSERT_TRUE(rc_rule_attacks_with_force(a, w, cases[k].force) == RC_OK);
            rc_results *res = NULL; size_t n;
            const rc_eval_result *it = eval_first(g, t == 1, &res, &n);
            ASSERT_TRUE(n == 1);
            if (it) {
                RC_CHECK(fabs(it->verdict - cases[k].verdict) < 1e-12, cases[k].name);
                RC_CHECK(fabs(it->applicability - cases[k].applic) < 1e-12, cases[k].name);
            }
            rc_results_destroy(res);
            rc_graph_destroy(g);
        }
    }

    /* 기본값: rc_rule_attacks == rc_rule_attacks_with_force(…, 1.0) (비트 단위). */
    {
        static const double half = 0.5;
        double got[2];
        for (int variant = 0; variant < 2; variant++) {
            rc_graph *g = rc_graph_create("default");
            rc_rule *w = rc_graph_rule_verdict(g, vv_fn, (void *)&one);
            rc_rule *c = rc_graph_counter_verdict(g, vv_fn, (void *)&half);
            rc_rule *d = rc_graph_except_verdict(g, vv_fn, (void *)&one_b);
            rc_rule_qualifier(d, 0.3);
            if (variant == 0) { rc_rule_attacks(c, w); rc_rule_attacks(d, w); }
            else { rc_rule_attacks_with_force(c, w, 1.0); rc_rule_attacks_with_force(d, w, 1.0); }
            rc_results *res = NULL; size_t n;
            const rc_eval_result *it = eval_first(g, false, &res, &n);
            got[variant] = it ? it->verdict : NAN;
            rc_results_destroy(res);
            rc_graph_destroy(g);
        }
        ASSERT_TRUE(got[0] == got[1]);
        /* d: q 0.3 → v_d = -0.4 → x = 0.3; c: v = 0.5, q 1 → x = 0.75.
           a_eff = 0.7, raw = 0.7/1.75 = 0.4 → verdict -0.2 */
        ASSERT_NEAR(got[0], -0.2, 1e-12);
    }

    /* 검증: NaN·범위 밖 → RC_ERR_INVALID_ARG + 빌더 오류 래치, 간선 미추가. 경계 0·1은 합법. */
    {
        double bad[3] = { NAN, -0.1, 1.1 };
        for (int k = 0; k < 3; k++) {
            rc_graph *g = rc_graph_create("force-bad");
            rc_context *ctx = rc_map_context_create();
            rc_rule *w = rc_graph_rule_verdict(g, vv_fn, (void *)&one);
            rc_rule *c = rc_graph_counter_verdict(g, vv_fn, (void *)&one_b);
            ASSERT_TRUE(rc_rule_attacks_with_force(c, w, bad[k]) == RC_ERR_INVALID_ARG);
            ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "force must be between 0.0 and 1.0"));
            rc_results *res = NULL;
            ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_ERR_BUILDER);
            ASSERT_NULL(res);
            rc_run_result *rr = NULL;
            ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_ERR_BUILDER);
            ASSERT_NULL(rr);
            rc_map_context_destroy(ctx);
            rc_graph_destroy(g);
        }
        rc_graph *g = rc_graph_create("force-ok");
        rc_graph *g2 = rc_graph_create("other");
        rc_rule *w = rc_graph_rule_verdict(g, vv_fn, (void *)&one);
        rc_rule *c = rc_graph_counter_verdict(g, vv_fn, (void *)&one_b);
        rc_rule *x = rc_graph_rule_verdict(g2, vv_fn, (void *)&one);
        ASSERT_TRUE(rc_rule_attacks_with_force(c, w, 0.0) == RC_OK);
        ASSERT_TRUE(rc_rule_attacks_with_force(c, w, 1.0) == RC_OK);
        ASSERT_TRUE(rc_rule_attacks_with_force(NULL, w, 0.5) == RC_ERR_INVALID_ARG);
        ASSERT_TRUE(rc_rule_attacks_with_force(c, NULL, 0.5) == RC_ERR_INVALID_ARG);
        ASSERT_TRUE(rc_rule_attacks_with_force(x, w, 0.5) == RC_ERR_INVALID_ARG);  /* 다른 graph */
        rc_results *res = NULL; size_t n;
        eval_first(g, false, &res, &n);   /* NULL 인자 오류는 래치하지 않는다 */
        ASSERT_TRUE(n == 1);
        rc_results_destroy(res);
        rc_graph_destroy(g);
        rc_graph_destroy(g2);
    }
}

/* --- Qualifier NaN --- */

static void test_qualifier_nan(void) {
    rc_graph *g = rc_graph_create("qnan");
    rc_context *ctx = rc_map_context_create();
    rc_rule *w = rc_graph_rule(g, vb_fn, (void *)&B_T);
    ASSERT_PTR_EQ(rc_rule_qualifier(w, NAN), w);   /* 체이닝 위해 self 반환 */
    ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "qualifier must be between 0.0 and 1.0"));
    rc_results *res = NULL;
    ASSERT_TRUE(rc_graph_evaluate(g, ctx, NULL, &res) == RC_ERR_BUILDER);
    ASSERT_NULL(res);
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
}

/* --- trace의 rule_verdict, Run 경로 --- */

static int fired = 0;
static rc_status count_fire(const rc_trace_entry *self, const rc_trace *t, void *ud)
{ (void)self; (void)t; (void)ud; fired++; return RC_OK; }

static void test_trace_rule_verdict(void) {
    static const double half = 0.5, zero = 0.0, m1 = -1.0;
    /* Evaluate trace: w(v 0.5) ← c(v 0) — 엔트리마다 규칙 자신의 v. */
    {
        rc_graph *g = rc_graph_create("trace-v");
        rc_rule *w = rc_graph_rule_verdict(g, vv_fn, (void *)&half);
        rc_rule *c = rc_graph_counter_verdict(g, vv_fn, (void *)&zero);
        rc_rule_attacks(c, w);
        rc_results *res = NULL; size_t n;
        const rc_eval_result *it = eval_first(g, true, &res, &n);
        ASSERT_TRUE(n == 1);
        if (it) {
            ASSERT_TRUE(it->n_trace == 2);
            if (it->n_trace == 2) {
                ASSERT_TRUE(it->trace[0].rule_verdict == 0.5);
                ASSERT_TRUE(it->trace[0].activated);
                ASSERT_TRUE(it->trace[1].rule_verdict == 0.0);
                ASSERT_TRUE(it->trace[1].activated);   /* 0은 Active */
            }
            /* a = 0.75, S = 0.5·(c verdict 0 + 1) = 0.5 → raw 0.5 → verdict 0 */
            ASSERT_NEAR(it->verdict, 0.0, 1e-12);
        }
        rc_results_destroy(res);
        rc_graph_destroy(g);
    }
    /* Run trace: 전 노드 rule_verdict, Active(v > -1 && verdict > 0)만 발화. */
    {
        rc_graph *g = rc_graph_create("run-v");
        rc_rule *a = rc_graph_rule_verdict(g, vv_fn, (void *)&half);   /* verdict 0.5 → 발화 */
        rc_rule *b = rc_graph_rule_verdict(g, vv_fn, (void *)&zero);   /* verdict 0 → 미발화 */
        rc_rule *c = rc_graph_except_verdict(g, vv_fn, (void *)&m1);   /* 해당 없음 */
        rc_rule_run_on(a, count_fire, NULL);
        rc_rule_run_on(b, count_fire, NULL);
        rc_rule_run_on(c, count_fire, NULL);
        rc_context *ctx = rc_map_context_create();
        rc_run_result *rr = NULL;
        fired = 0;
        ASSERT_TRUE(rc_graph_run(g, ctx, NULL, &rr) == RC_OK);
        ASSERT_TRUE(fired == 1);
        size_t n = 0;
        const rc_trace_entry *all = rc_trace_all(rc_run_result_trace(rr), &n);
        ASSERT_TRUE(n == 3);
        if (n == 3) {
            ASSERT_TRUE(all[0].rule_verdict == 0.5);
            ASSERT_TRUE(all[1].rule_verdict == 0.0);
            ASSERT_TRUE(all[1].activated);
            ASSERT_TRUE(all[2].rule_verdict == -1.0);
            ASSERT_FALSE(all[2].activated);
            ASSERT_STR_EQ(all[2].role, "except");
        }
        size_t nr = 0;
        const rc_eval_result *items = rc_results_items(rc_run_result_results(rr), &nr);
        ASSERT_TRUE(nr == 2);
        if (nr == 2) {
            ASSERT_TRUE(items[0].applicability == 0.75);
            ASSERT_TRUE(items[1].applicability == 0.5);
        }
        rc_run_result_destroy(rr);
        rc_map_context_destroy(ctx);
        rc_graph_destroy(g);
    }
}

/* --- verdict 등록 API --- */

static void test_verdict_register(void) {
    static const double one = 1.0;
    rc_graph *g = rc_graph_create("reg");
    ASSERT_NULL(rc_graph_rule_verdict(g, NULL, NULL));
    ASSERT_NULL(rc_graph_counter_verdict(g, NULL, NULL));
    ASSERT_NULL(rc_graph_except_verdict(g, NULL, NULL));
    ASSERT_NULL(rc_graph_rule_verdict(NULL, vv_fn, NULL));
    rc_rule *a = rc_graph_rule_verdict(g, vv_fn, (void *)&one);
    ASSERT_NOT_NULL(a);
    /* 같은 (fn, user_data) → 중복 */
    ASSERT_NULL(rc_graph_counter_verdict(g, vv_fn, (void *)&one));
    ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "duplicate"));
    /* bool 규칙과 verdict 규칙은 같은 user_data여도 별개 정체성 */
    ASSERT_NOT_NULL(rc_graph_rule(g, vb_fn, (void *)&one));
    rc_graph_destroy(g);
}

void test_verdict(void) {
    test_bool_compat();
    test_verdict_range();
    test_partial_verdict();
    test_except_gate();
    test_attack_force();
    test_qualifier_nan();
    test_trace_rule_verdict();
    test_verdict_register();
}
