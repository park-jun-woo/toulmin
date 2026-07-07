#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <stdarg.h>
#include <stdio.h>

/* buf/len 방어 후 vsnprintf. errbuf NULL 허용 계약을 일괄 처리. */
void rc__fmt(char *buf, size_t len, const char *fmt, ...) {
    if (!buf || len == 0) return;
    va_list ap;
    va_start(ap, fmt);
    vsnprintf(buf, len, fmt, ap);
    va_end(ap);
}

rc_expectation rc_verdict_above(double v) {
    return (rc_expectation){ RC_EXPECT_ABOVE, v, 0 };
}

rc_expectation rc_verdict_at_most(double v) {
    return (rc_expectation){ RC_EXPECT_AT_MOST, v, 0 };
}

rc_expectation rc_verdict_between(double lo, double hi) {
    return (rc_expectation){ RC_EXPECT_BETWEEN, lo, hi };
}

rc_expectation rc_no_result(void) {
    return (rc_expectation){ RC_EXPECT_NO_RESULT, 0, 0 };
}

/* 통과: RC_OK. 실패: RC_ERR_INVALID_ARG + errbuf에 사유(NULL 허용).
   경계 의미는 py expectations.py와 일치: >, <=, lo<v<=hi, 결과 0건. */
rc_status rc_expect_check(const rc_results *results, rc_expectation e,
                          char *errbuf, size_t errlen) {
    size_t n;
    const rc_eval_result *items = rc_results_items(results, &n);

    if (e.kind == RC_EXPECT_NO_RESULT) {
        if (n == 0) return RC_OK;
        rc__fmt(errbuf, errlen, "expected no results, got %zu", n);
        return RC_ERR_INVALID_ARG;
    }
    if (n == 0) {
        rc__fmt(errbuf, errlen, "expected at least one result");
        return RC_ERR_INVALID_ARG;
    }

    double v = items[0].verdict;
    switch (e.kind) {
    case RC_EXPECT_ABOVE:
        if (v > e.lo) return RC_OK;
        rc__fmt(errbuf, errlen, "expected verdict > %g, got %g", e.lo, v);
        break;
    case RC_EXPECT_AT_MOST:
        if (v <= e.lo) return RC_OK;
        rc__fmt(errbuf, errlen, "expected verdict <= %g, got %g", e.lo, v);
        break;
    case RC_EXPECT_BETWEEN:
        if (v > e.lo && v <= e.hi) return RC_OK;
        rc__fmt(errbuf, errlen, "expected %g < verdict <= %g, got %g",
                e.lo, e.hi, v);
        break;
    default:
        rc__fmt(errbuf, errlen, "unknown expectation kind %d", (int)e.kind);
        break;
    }
    return RC_ERR_INVALID_ARG;
}
