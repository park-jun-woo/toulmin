#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <stdlib.h>
#include <string.h>

/* n바이트 딥카피. n==0이면 NULL. 실패 시 NULL. */
static void *rc__memdup(const void *src, size_t n) {
    if (n == 0) return NULL;
    void *dst = malloc(n);
    if (!dst) return NULL;
    memcpy(dst, src, n);
    return dst;
}

/* opt가 NULL이면 기본값. duration은 trace를 함의(Go/TS/Py와 동일). */
rc_eval_option rc__resolve_option(const rc_eval_option *opt) {
    rc_eval_option o = { RC_MATRIX, false, false };
    if (opt) o = *opt;
    if (o.duration) o.trace = true;
    return o;
}

/* h-Categoriser 핵심 — 재귀 + verdict 캐시(memoization).
   Strict는 공격 합산 제외, 비활성 노드는 -1.0(캐시 미기록). */
static double rc__calc(rc_eval_state *es, size_t i) {
    if (es->err != RC_OK) return -1.0;
    if (es->has_verdict[i]) return es->verdict[i];

    rc_rule_meta *m = &es->g->rules[i];

    if (!es->ran[i]) {
        es->ran[i] = true;
        bool act = false;
        void *ev = NULL;
        rc_status st = m->fn(es->ctx, m->specs, m->n_specs, m->user_data, &act, &ev);
        if (st != RC_OK) {
            rc__set_error(es->g, RC_ERR_RULE,
                          "rule \"%s\": failed (status %d)", m->name, st);
            es->err = RC_ERR_RULE;
            return -1.0;
        }
        es->active[i]   = act;
        es->evidence[i] = ev;
    }

    if (!es->active[i]) return -1.0;

    double total = 0.0;
    if (m->strength != RC_STRICT) {           /* Strict는 공격 합산 제외 */
        for (size_t k = 0; k < es->n_attackers[i]; k++) {
            total += (rc__calc(es, es->attackers[i][k]) + 1.0) / 2.0;
            if (es->err != RC_OK) return -1.0;
        }
    }

    double v = 2.0 * (m->qualifier / (1.0 + total)) - 1.0;
    es->has_verdict[i] = true;
    es->verdict[i] = v;
    return v;
}

/* Evaluate trace의 role은 등록 role이 아니라 defeats 구조로 판정한다
   (Go inferRole / Py infer_role 동형). 등록 role은 Run trace(Phase 005)에서 사용. */
static const char *rc__infer_role(const rc_eval_state *es, size_t i) {
    if (es->g->rules[i].strength == RC_DEFEATER) return "except";
    if (es->is_attacker[i]) return "counter";
    return "rule";
}

/* trace 버퍼에 엔트리를 추가한다. 성장 실패 시 false. */
static bool rc__trace_append(rc_eval_state *es, rc_trace_entry e) {
    if (es->n_trace == es->cap_trace) {
        size_t nc = es->cap_trace ? es->cap_trace * 2 : 4;
        rc_trace_entry *grown = realloc(es->trace, nc * sizeof *grown);
        if (!grown) return false;
        es->trace = grown;
        es->cap_trace = nc;
    }
    es->trace[es->n_trace++] = e;
    return true;
}

/* h-Categoriser + trace 기록. calc과 동일한 재귀지만 방문 순서로 엔트리를 남긴다.
   duration이 켜지면 규칙 함수 실행 구간을 ms로 측정한다. */
static double rc__calc_trace(rc_eval_state *es, size_t i, bool with_duration) {
    if (es->err != RC_OK) return -1.0;
    if (es->has_verdict[i]) return es->verdict[i];

    rc_rule_meta *m = &es->g->rules[i];

    if (!es->ran[i]) {
        es->ran[i] = true;
        double start = with_duration ? rc__now_ms() : 0.0;
        bool act = false;
        void *ev = NULL;
        rc_status st = m->fn(es->ctx, m->specs, m->n_specs, m->user_data, &act, &ev);
        if (st != RC_OK) {
            rc__set_error(es->g, RC_ERR_RULE,
                          "rule \"%s\": failed (status %d)", m->name, st);
            es->err = RC_ERR_RULE;
            return -1.0;
        }
        es->active[i]   = act;
        es->evidence[i] = ev;
        double elapsed = with_duration ? rc__now_ms() - start : 0.0;

        rc_trace_entry entry = {
            .name        = m->short_name,          /* graph 소유 문자열 차용 */
            .role        = rc__infer_role(es, i),  /* 구조 기반 판정 */
            .activated   = act,
            .qualifier   = m->qualifier,
            .verdict     = 0.0,                     /* 기록 시점엔 미확정 */
            .evidence    = ev,
            .ground      = NULL,                    /* Evaluate 경로에선 미기록 */
            .specs       = m->specs,               /* = Backing, 차용 */
            .n_specs     = m->n_specs,
            .duration_ms = elapsed,
        };
        if (!rc__trace_append(es, entry)) {
            rc__set_error(es->g, RC_ERR_NOMEM, "trace: out of memory");
            es->err = RC_ERR_NOMEM;
            return -1.0;
        }
    }

    if (!es->active[i]) return -1.0;

    double total = 0.0;
    if (m->strength != RC_STRICT) {
        for (size_t k = 0; k < es->n_attackers[i]; k++) {
            total += (rc__calc_trace(es, es->attackers[i][k], with_duration) + 1.0) / 2.0;
            if (es->err != RC_OK) return -1.0;
        }
    }

    double v = 2.0 * (m->qualifier / (1.0 + total)) - 1.0;
    es->has_verdict[i] = true;
    es->verdict[i] = v;
    return v;
}

/* trace 모드에서 워런트마다 상태를 비운다. evidence 배열은 0으로 클리어하되
   해제하지 않는다 — 소유권은 이미 results(owned_evidence)로 이관된 상태. */
static void rc__reset(rc_eval_state *es) {
    size_t n = es->g->n_rules;
    if (n) {
        memset(es->ran,         0, n * sizeof *es->ran);
        memset(es->active,      0, n * sizeof *es->active);
        memset(es->has_verdict, 0, n * sizeof *es->has_verdict);
        memset(es->evidence,    0, n * sizeof *es->evidence);
    }
    es->n_trace = 0;
    es->err = RC_OK;
}

/* Defeater가 아니고 아무도 공격하지 않는 노드가 warrant. */
static bool rc__is_warrant(const rc_eval_state *es, size_t i) {
    if (es->g->rules[i].strength == RC_DEFEATER) return false;
    return !es->is_attacker[i];
}

/* 공용 평가 본체 — rc_graph_evaluate(full=false)와 rc_graph_run(full=true)이 공유한다.
   성공 시 *out_res에 results, *es에 (호출자 소유) 평가 상태를 남긴다 — 호출자가 evidence를
   이관하고 rc__eval_state_free한다. 실패 시 내부에서 es와 부분 results를 모두 정리한다.
   builder_err 체크가 여기에 있어야 rc__run_depth가 실행하는 하위 그래프의 지연 빌더 오류도
   합성 경로에서 RC_ERR_BUILDER로 보고된다. */
rc_status rc__evaluate_internal(rc_graph *g, rc_context *ctx, rc_eval_option o,
                                bool full, rc_results **out_res, rc_eval_state *es) {
    *out_res = NULL;
    if (g->builder_err != RC_OK) return g->builder_err;   /* 지연 빌더 오류 */

    rc_status st = rc__eval_state_init(es, g, ctx);
    if (st != RC_OK) return st;

    size_t cycle_at = 0;
    st = rc__detect_cycle(es, &cycle_at);
    if (st != RC_OK) {
        if (st == RC_ERR_CYCLE)
            rc__set_error(g, RC_ERR_CYCLE, "cycle detected at %s",
                          g->rules[cycle_at].name);
        rc__eval_state_free(es);
        return st;
    }

    rc_results *res = rc__results_create(g);
    if (!res) { rc__eval_state_free(es); return RC_ERR_NOMEM; }

    for (size_t i = 0; i < g->n_rules; i++) {
        if (!rc__is_warrant(es, i)) continue;

        /* trace 모드: 직전 워런트 패스의 evidence를 results로 이관 후 상태 초기화. */
        if (o.trace) {
            rc__results_adopt_evidence(res, es);
            rc__reset(es);
        }

        double v = o.trace ? rc__calc_trace(es, i, o.duration)
                           : rc__calc(es, i);
        if (es->err != RC_OK) {
            st = es->err;
            rc__results_abandon(res, es);
            rc__eval_state_free(es);
            return st;
        }
        if (!es->active[i]) continue;

        rc_eval_result item = {
            .name     = g->rules[i].short_name,   /* graph 소유 문자열 차용 */
            .verdict  = v,
            .evidence = es->evidence[i],
            .trace    = NULL,
            .n_trace  = 0,
        };
        if (o.trace && es->n_trace > 0) {
            /* trace 버퍼를 결과에 딥카피 (엔트리 배열만 — 문자열/spec은 graph 차용). */
            item.trace = rc__memdup(es->trace, es->n_trace * sizeof *es->trace);
            if (!item.trace) {
                rc__results_abandon(res, es);
                rc__eval_state_free(es);
                return RC_ERR_NOMEM;
            }
            item.n_trace = es->n_trace;
        }
        if (!rc__results_append(res, item)) {
            free(item.trace);
            rc__results_abandon(res, es);
            rc__eval_state_free(es);
            return RC_ERR_NOMEM;
        }
    }

    if (full) {   /* run 경로: trace 강제 off이므로 non-trace rc__calc으로 미평가 노드를 채움 */
        for (size_t i = 0; i < g->n_rules; i++)
            if (!es->ran[i]) rc__calc(es, i);
        if (es->err != RC_OK) {
            st = es->err;
            rc__results_abandon(res, es);
            rc__eval_state_free(es);
            return st;
        }
    }

    /* evidence 이관과 es 해제는 호출자 책임 (run은 trace 조립 후 이관해야 하므로). */
    *out_res = res;
    return RC_OK;
}

rc_status rc_graph_evaluate(rc_graph *g, rc_context *ctx,
                            const rc_eval_option *opt, rc_results **out) {
    if (out) *out = NULL;
    if (!g || !ctx || !out) return RC_ERR_INVALID_ARG;

    rc_eval_option o = rc__resolve_option(opt);

    rc_eval_state es;
    rc_results *res = NULL;
    rc_status st = rc__evaluate_internal(g, ctx, o, /*full=*/false, &res, &es);
    if (st != RC_OK) return st;   /* 실패 시 es는 내부에서 이미 정리됨 */

    rc__results_adopt_evidence(res, &es);   /* 마지막 패스 evidence 이관 */
    rc__eval_state_free(&es);
    *out = res;
    return RC_OK;
}
