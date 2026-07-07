#define _POSIX_C_SOURCE 200809L   /* strdup — 전 .c 공통 규약 */

#include <rulecat/rulecat.h>

#include <stdlib.h>
#include <string.h>

typedef struct rc__map_node {
    char *key;                 /* strdup 소유 */
    void *value;               /* 차용 — Context는 value를 해제하지 않는다 */
    struct rc__map_node *next;
} rc__map_node;

static void *rc__map_get(rc_context *self, const char *key) {
    for (rc__map_node *n = self->impl; n; n = n->next)
        if (strcmp(n->key, key) == 0) return n->value;
    return NULL;
}

static void rc__map_set(rc_context *self, const char *key, void *value) {
    for (rc__map_node *n = self->impl; n; n = n->next)
        if (strcmp(n->key, key) == 0) { n->value = value; return; }
    rc__map_node *n = malloc(sizeof *n);
    if (!n) return;                     /* set은 void — 실패 시 무시(문서화) */
    n->key = strdup(key);
    n->value = value;
    n->next = self->impl;
    self->impl = n;
}

rc_context *rc_map_context_create(void) {
    rc_context *ctx = calloc(1, sizeof *ctx);
    if (!ctx) return NULL;
    ctx->get = rc__map_get;
    ctx->set = rc__map_set;
    return ctx;
}

void rc_map_context_destroy(rc_context *ctx) {
    if (!ctx) return;
    rc__map_node *n = ctx->impl;
    while (n) { rc__map_node *next = n->next; free(n->key); free(n); n = next; }
    free(ctx);
}
