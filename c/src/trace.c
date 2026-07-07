#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <stdlib.h>
#include <string.h>

/* n개 노드 슬롯(0 초기화)을 가진 trace를 만든다. ctx는 차용. n==0이면 nodes는 NULL. */
rc_trace *rc__trace_create(size_t n, rc_context *ctx) {
    rc_trace *t = malloc(sizeof *t);
    if (!t) return NULL;
    t->nodes = NULL;
    if (n) {
        t->nodes = calloc(n, sizeof *t->nodes);
        if (!t->nodes) { free(t); return NULL; }
    }
    t->len = n;
    t->ctx = ctx;
    return t;
}

/* nodes 배열과 컨테이너만 해제한다 — 엔트리의 문자열/spec/evidence는 차용이라 손대지 않는다
   (evidence 소유권은 rc_results가 단독으로 가지며 rc_results_destroy에서 1회 해제). */
void rc__trace_destroy_internal(rc_trace *t) {
    if (!t) return;
    free(t->nodes);
    free(t);
}

const rc_trace_entry *rc_trace_all(const rc_trace *t, size_t *out_len) {
    if (out_len) *out_len = t ? t->len : 0;
    return t ? t->nodes : NULL;
}

const rc_trace_entry *rc_trace_get(const rc_trace *t, const char *name) {
    if (!t || !name) return NULL;
    for (size_t i = 0; i < t->len; i++)
        if (t->nodes[i].name && strcmp(t->nodes[i].name, name) == 0)
            return &t->nodes[i];
    return NULL;
}

rc_context *rc_trace_ctx(const rc_trace *t) {
    return t ? t->ctx : NULL;
}
