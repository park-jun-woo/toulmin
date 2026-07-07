#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <stdlib.h>

/* opaque rc_results 실체. */
struct rc_results {
    rc_eval_result *items;
    size_t          len, cap;
    /* dtor 등록 시 evidence를 이관받아 소유. trace 모드에선 워런트 패스마다
       추가되므로 append+grow 구조 + 포인터 중복 제거로 이중 해제를 막는다. */
    void          **owned_evidence;
    size_t          n_owned, cap_owned;
    void          (*evidence_dtor)(void *);
};

/* graph의 evidence_dtor를 연결한 빈 컨테이너를 만든다. */
rc_results *rc__results_create(rc_graph *g) {
    rc_results *r = calloc(1, sizeof *r);
    if (!r) return NULL;
    r->evidence_dtor = g->evidence_dtor;
    return r;
}

/* 결과 항목 하나를 추가한다. 성장 실패 시 false. */
bool rc__results_append(rc_results *r, rc_eval_result item) {
    if (r->len == r->cap) {
        size_t new_cap = r->cap ? r->cap * 2 : 4;
        rc_eval_result *grown = realloc(r->items, new_cap * sizeof *grown);
        if (!grown) return false;
        r->items = grown;
        r->cap = new_cap;
    }
    r->items[r->len++] = item;
    return true;
}

/* 이미 소유한 포인터인지 검사(distinct pointer마다 1회 계약). */
static bool rc__owned_has(const rc_results *r, const void *ev) {
    for (size_t k = 0; k < r->n_owned; k++)
        if (r->owned_evidence[k] == ev) return true;
    return false;
}

/* dtor 등록 시 이 패스의 evidence를 소유로 이관한다.
   trace 모드에선 여러 패스에 걸쳐 누적되며, 같은 포인터는 한 번만 담는다. */
void rc__results_adopt_evidence(rc_results *r, rc_eval_state *es) {
    if (!r || !r->evidence_dtor) return;
    size_t n = es->g->n_rules;
    for (size_t i = 0; i < n; i++) {
        void *ev = es->evidence[i];
        if (!ev || rc__owned_has(r, ev)) continue;   /* 중복 제거 */
        if (r->n_owned == r->cap_owned) {
            size_t nc = r->cap_owned ? r->cap_owned * 2 : 8;
            void **grown = realloc(r->owned_evidence, nc * sizeof *grown);
            if (!grown) {
                /* 이관 실패 — 누수 방지 위해 이 distinct 포인터를 즉시 해제 */
                r->evidence_dtor(ev);
                continue;
            }
            r->owned_evidence = grown;
            r->cap_owned = nc;
        }
        r->owned_evidence[r->n_owned++] = ev;
    }
}

/* 평가 실패 시 생산된 evidence를 정리하고 컨테이너를 해제한다.
   이전 패스에서 이관한 owned + 현재 패스 evidence를 distinct 포인터당 1회 해제. */
void rc__results_abandon(rc_results *r, rc_eval_state *es) {
    if (!r) return;
    if (r->evidence_dtor) {
        for (size_t k = 0; k < r->n_owned; k++)
            if (r->owned_evidence[k]) r->evidence_dtor(r->owned_evidence[k]);
        for (size_t i = 0; i < es->g->n_rules; i++) {
            void *ev = es->evidence[i];
            if (!ev || rc__owned_has(r, ev)) continue;
            bool dup = false;                        /* 현재 패스 내 중복 방지 */
            for (size_t j = 0; j < i; j++)
                if (es->evidence[j] == ev) { dup = true; break; }
            if (!dup) r->evidence_dtor(ev);
        }
    }
    free(r->owned_evidence);
    for (size_t i = 0; i < r->len; i++) free(r->items[i].trace);   /* Phase 004 */
    free(r->items);
    free(r);
}

const rc_eval_result *rc_results_items(const rc_results *r, size_t *out_len) {
    if (out_len) *out_len = r ? r->len : 0;
    return r ? r->items : NULL;
}

void rc_results_destroy(rc_results *r) {
    if (!r) return;
    if (r->evidence_dtor)
        for (size_t i = 0; i < r->n_owned; i++)
            if (r->owned_evidence[i]) r->evidence_dtor(r->owned_evidence[i]);
    free(r->owned_evidence);
    for (size_t i = 0; i < r->len; i++) free(r->items[i].trace);   /* Phase 004 */
    free(r->items);
    free(r);
}
