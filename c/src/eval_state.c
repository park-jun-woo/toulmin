#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <stdlib.h>
#include <string.h>

/* rc_eval_state의 모든 배열을 해제한다(evidence 자체는 소유권 규약상 손대지 않음). */
void rc__eval_state_free(rc_eval_state *es) {
    if (!es) return;
    if (es->attackers && es->g) {
        for (size_t i = 0; i < es->g->n_rules; i++) free(es->attackers[i]);
    }
    free(es->attackers);
    free(es->n_attackers);
    free(es->is_attacker);
    free(es->ran);
    free(es->active);
    free(es->has_verdict);
    free(es->verdict);
    free(es->evidence);
    free(es->trace);
    memset(es, 0, sizeof *es);
}

/* 인덱스 배열과 attacker 인접 리스트를 구성한다.
   0 초기화가 곧 기본값(ran=false, verdict 캐시 없음). */
rc_status rc__eval_state_init(rc_eval_state *es, rc_graph *g, rc_context *ctx) {
    memset(es, 0, sizeof *es);
    es->g = g;
    es->ctx = ctx;
    es->err = RC_OK;

    size_t n = g->n_rules;
    if (n == 0) return RC_OK;   /* 노드 없음 — 배열은 NULL 유지 */

    es->attackers   = calloc(n, sizeof *es->attackers);
    es->n_attackers = calloc(n, sizeof *es->n_attackers);
    es->is_attacker = calloc(n, sizeof *es->is_attacker);
    es->ran         = calloc(n, sizeof *es->ran);
    es->active      = calloc(n, sizeof *es->active);
    es->has_verdict = calloc(n, sizeof *es->has_verdict);
    es->verdict     = calloc(n, sizeof *es->verdict);
    es->evidence    = calloc(n, sizeof *es->evidence);
    if (!es->attackers || !es->n_attackers || !es->is_attacker || !es->ran ||
        !es->active || !es->has_verdict || !es->verdict || !es->evidence) {
        rc__eval_state_free(es);
        return RC_ERR_NOMEM;
    }

    /* 1) 노드별 attacker 수 집계 */
    for (size_t e = 0; e < g->n_defeats; e++)
        es->n_attackers[g->defeats[e].to_idx]++;

    /* 2) 리스트 할당 후 n_attackers를 채움 카운터로 재사용 */
    for (size_t i = 0; i < n; i++) {
        if (es->n_attackers[i]) {
            es->attackers[i] = malloc(es->n_attackers[i] * sizeof **es->attackers);
            if (!es->attackers[i]) { rc__eval_state_free(es); return RC_ERR_NOMEM; }
        }
        es->n_attackers[i] = 0;
    }

    /* 3) 간선 채우기 + attacker set 표시 */
    for (size_t e = 0; e < g->n_defeats; e++) {
        size_t from = g->defeats[e].from_idx;
        size_t to   = g->defeats[e].to_idx;
        es->attackers[to][es->n_attackers[to]++] = from;
        es->is_attacker[from] = true;
    }

    return RC_OK;
}
