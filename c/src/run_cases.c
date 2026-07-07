#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

/* 테이블 주도 케이스를 순서대로 실행한다. 첫 실패에서 중단하며
   errbuf에 "[케이스명] 사유"를 기록한다(errbuf NULL 허용). */
rc_status rc_run_cases(rc_graph *g, const rc_test_case *cases, size_t n,
                       char *errbuf, size_t errlen) {
    if (!g || (!cases && n > 0)) return RC_ERR_INVALID_ARG;

    for (size_t i = 0; i < n; i++) {
        const rc_test_case *tc = &cases[i];
        rc_context *ctx = tc->ctx ? tc->ctx : rc_map_context_create();
        if (!ctx) {
            rc__fmt(errbuf, errlen, "[%s] context create failed", tc->name);
            return RC_ERR_NOMEM;
        }

        rc_results *res = NULL;
        rc_status st = rc_graph_evaluate(g, ctx, tc->opt, &res);
        if (st != RC_OK) {
            rc__fmt(errbuf, errlen, "[%s] evaluate: %s",
                    tc->name, rc_graph_last_error(g));
        } else {
            char why[160];
            st = rc_expect_check(res, tc->expect, why, sizeof why);
            if (st != RC_OK) rc__fmt(errbuf, errlen, "[%s] %s", tc->name, why);
        }

        rc_results_destroy(res);
        if (!tc->ctx) rc_map_context_destroy(ctx);
        if (st != RC_OK) return st;
    }
    return RC_OK;
}
