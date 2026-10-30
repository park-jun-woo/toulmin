#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <stdlib.h>

/* 실행 합성 재귀 깊이 백스톱 — 정적 순환 감지가 못 보는 폭주 합성 대비(런타임 방어). */
enum { RC_RUN_MAX_DEPTH = 64 };

const rc_results *rc_run_result_results(const rc_run_result *rr) {
    return rr ? rr->results : NULL;
}

const rc_trace *rc_run_result_trace(const rc_run_result *rr) {
    return rr ? rr->trace : NULL;
}

/* trace(nodes 배열)와 results(evidence 풀)를 함께 해제한다 — evidence는 여기서 정확히 1회.
   trace 엔트리는 같은 evidence 포인터를 차용만 하므로 이중 해제가 없다. */
void rc_run_result_destroy(rc_run_result *rr) {
    if (!rr) return;
    rc__trace_destroy_internal(rr->trace);   /* nodes 배열만 free — 차용 항목은 유지 */
    rc_results_destroy(rr->results);         /* evidence dtor 1회 호출 지점 */
    free(rr);
}

/* full pass 평가 → 전 노드 평면 trace(등록 순서 조립, 인덱스 i ↔ g->rules[i]) → Active
   노드만 (a) 핸들러 → (b) 하위 그래프 Run 순으로 디스패치. ctx는 아래로 흐르고(공유 가변),
   각 레벨은 자기 trace를 만든다. 하위 verdict는 격리되고 오류만 위로 전파된다. */
static rc_status rc__run_depth(rc_graph *g, rc_context *ctx, rc_eval_option o,
                               int depth, rc_run_result **out) {
    if (depth > RC_RUN_MAX_DEPTH) {
        rc__set_error(g, RC_ERR_RUN_DEPTH,
                      "run depth exceeded %d (possible runaway composition)",
                      RC_RUN_MAX_DEPTH);
        return RC_ERR_RUN_DEPTH;
    }

    rc_eval_state es;
    rc_results *results = NULL;
    rc_status st = rc__evaluate_internal(g, ctx, o, /*full=*/true, &results, &es);
    if (st != RC_OK) return st;   /* 실패 시 es/results는 내부에서 이미 정리됨 */

    /* 전 노드 평면 trace — 등록 순서. inactive 노드는 verdict 0.0(캐시 미존재). */
    rc_trace *tr = rc__trace_create(g->n_rules, ctx);
    if (!tr) {
        rc__results_adopt_evidence(results, &es);   /* evidence 회수 후 폐기 */
        rc__eval_state_free(&es);
        rc_results_destroy(results);
        return RC_ERR_NOMEM;
    }
    for (size_t i = 0; i < g->n_rules; i++) {
        rc_rule_meta *m = &g->rules[i];
        tr->nodes[i] = (rc_trace_entry){
            .name        = m->short_name,   /* = Claim */
            .role        = m->role,         /* Run trace는 등록 role 사용 */
            .activated   = es.active[i],
            .qualifier   = m->qualifier,
            .verdict     = es.has_verdict[i] ? es.verdict[i] : 0.0,
            .evidence    = es.evidence[i],
            .ground      = ctx,             /* = Ground: ctx 그대로, 전 entry 동일 */
            .specs       = m->specs,        /* = Backing */
            .n_specs     = m->n_specs,
            .duration_ms = 0.0,             /* full pass는 non-trace — 측정 없음 */
            .rule_verdict = es.ran[i] ? es.rule_verdict[i] : -1.0,   /* 규칙 자신의 v */
        };
    }

    rc__results_adopt_evidence(results, &es);   /* evidence 풀 이관(노드당 1회) */
    rc__eval_state_free(&es);

    rc_run_result *rr = malloc(sizeof *rr);
    if (!rr) {
        rc__trace_destroy_internal(tr);
        rc_results_destroy(results);
        return RC_ERR_NOMEM;
    }
    rr->results = results;
    rr->trace = tr;
    *out = rr;   /* 이후 오류여도 rr은 유효 — 호출자 destroy */

    /* 디스패치: Active(activated(규칙 v > -1) && verdict > 0)만, 등록 순서(결정적). */
    for (size_t i = 0; i < g->n_rules; i++) {
        rc_rule_meta *m = &g->rules[i];
        const rc_trace_entry *self = &tr->nodes[i];
        if (!(self->activated && self->verdict > 0.0)) continue;

        if (m->run_on) {   /* (a) 핸들러 먼저 */
            rc_status hst = m->run_on(self, tr, m->run_on_data);
            if (hst != RC_OK) {
                rc__set_error(g, RC_ERR_HANDLER,
                              "run_on \"%s\": handler failed (status %d)",
                              self->name, hst);
                return RC_ERR_HANDLER;   /* *out 유지 — trace 검시 가능 */
            }
        }
        if (m->run_graph) {   /* (b) 실행 간선 — 하위 그래프 Run */
            rc_run_result *sub_rr = NULL;
            rc_status sst = rc__run_depth(m->run_graph, ctx, o, depth + 1, &sub_rr);
            rc_run_result_destroy(sub_rr);   /* verdict 격리 — 하위 결과는 위로 안 샘 */
            if (sst != RC_OK) {
                rc__set_error(g, sst, "run \"%s\" -> \"%s\": %s",
                              self->name, m->run_graph->name,
                              rc_graph_last_error(m->run_graph));
                return sst;
            }
        }
    }
    return RC_OK;
}

rc_status rc_graph_run(rc_graph *g, rc_context *ctx,
                       const rc_eval_option *opt, rc_run_result **out) {
    if (out) *out = NULL;
    if (!g || !ctx || !out) return RC_ERR_INVALID_ARG;
    if (g->builder_err != RC_OK) return g->builder_err;   /* 순환 검사 전 fast-fail */

    /* 실행 합성 순환 — 최상위 1회. */
    const rc_graph *cyc = NULL;
    rc__gcolor_set set = {0};
    rc_status st = rc__run_cycle_dfs(g, &set, &cyc);
    rc__gcolor_set_free(&set);
    if (st == RC_ERR_RUN_CYCLE) {
        rc__set_error(g, st, "run cycle detected at graph \"%s\"",
                      cyc && cyc->name ? cyc->name : "");
        return st;
    }
    if (st != RC_OK) return st;   /* 색칠 집합 할당 실패(NOMEM) */

    rc_eval_option o = rc__resolve_option(opt);
    o.trace = false;      /* 강제 off — full pass는 공유·비초기화 상태에서만 정확 */
    o.duration = false;
    return rc__run_depth(g, ctx, o, 0, out);
}
