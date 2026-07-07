#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <stdlib.h>

/* 실행 합성 순환 감지: RunGraph 간선으로 이어지는 그래프-of-그래프를 3색 DFS로 순회한다
   (Python detect_run_cycle의 그래프 객체 identity 키잉 동형). GRAY 재진입은 실행 순환,
   BLACK로 끝난 공유 하위 그래프는 두 경로에서 도달해도 합법(다이아몬드 DAG). */
enum { RC__G_WHITE = 0, RC__G_GRAY = 1, RC__G_BLACK = 2 };

/* 미등록 그래프는 WHITE로 간주한다(선형 탐색 — 합성 그래프 수는 작다). */
static unsigned char rc__gcolor_set_get(const rc__gcolor_set *set, const rc_graph *g) {
    for (size_t i = 0; i < set->len; i++)
        if (set->items[i].g == g) return set->items[i].color;
    return RC__G_WHITE;
}

/* 그래프의 색을 갱신하거나 새로 등록한다. 성장 실패 시 RC_ERR_NOMEM. */
static rc_status rc__gcolor_set_put(rc__gcolor_set *set, const rc_graph *g,
                                    unsigned char color) {
    for (size_t i = 0; i < set->len; i++)
        if (set->items[i].g == g) { set->items[i].color = color; return RC_OK; }
    if (set->len == set->cap) {
        size_t nc = set->cap ? set->cap * 2 : 8;
        rc__gcolor *grown = realloc(set->items, nc * sizeof *grown);
        if (!grown) return RC_ERR_NOMEM;
        set->items = grown;
        set->cap = nc;
    }
    set->items[set->len].g = g;
    set->items[set->len].color = color;
    set->len++;
    return RC_OK;
}

void rc__gcolor_set_free(rc__gcolor_set *set) {
    if (!set) return;
    free(set->items);
    set->items = NULL;
    set->len = set->cap = 0;
}

rc_status rc__run_cycle_dfs(const rc_graph *g, rc__gcolor_set *set,
                            const rc_graph **cycle_at) {
    rc_status st = rc__gcolor_set_put(set, g, RC__G_GRAY);
    if (st != RC_OK) return st;

    for (size_t i = 0; i < g->n_rules; i++) {
        const rc_graph *sub = g->rules[i].run_graph;
        if (!sub) continue;
        unsigned char c = rc__gcolor_set_get(set, sub);   /* 미등록 → WHITE */
        if (c == RC__G_GRAY) { *cycle_at = sub; return RC_ERR_RUN_CYCLE; }
        if (c == RC__G_WHITE) {
            st = rc__run_cycle_dfs(sub, set, cycle_at);
            if (st != RC_OK) return st;
        }
    }

    return rc__gcolor_set_put(set, g, RC__G_BLACK);
}
