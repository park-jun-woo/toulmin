#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <stdlib.h>

enum { RC__WHITE = 0, RC__GRAY = 1, RC__BLACK = 2 };

/* GRAY 노드를 다시 만나면 순환. 재귀 깊이는 최대 규칙 수(수십 규모)라 재귀로 충분. */
static rc_status rc__cycle_dfs(rc_eval_state *es, size_t node,
                               unsigned char *color, size_t *cycle_at) {
    color[node] = RC__GRAY;
    for (size_t k = 0; k < es->n_attackers[node]; k++) {
        size_t nb = es->attackers[node][k];
        if (color[nb] == RC__GRAY) { *cycle_at = nb; return RC_ERR_CYCLE; }
        if (color[nb] == RC__WHITE) {
            rc_status st = rc__cycle_dfs(es, nb, color, cycle_at);
            if (st != RC_OK) return st;
        }
    }
    color[node] = RC__BLACK;
    return RC_OK;
}

/* defeats 그래프 전체를 순회하며 순환을 감지한다. */
rc_status rc__detect_cycle(rc_eval_state *es, size_t *cycle_at) {
    size_t n = es->g->n_rules;
    if (n == 0) return RC_OK;
    unsigned char *color = calloc(n, 1);
    if (!color) return RC_ERR_NOMEM;
    rc_status st = RC_OK;
    for (size_t i = 0; i < n && st == RC_OK; i++)
        if (color[i] == RC__WHITE)
            st = rc__cycle_dfs(es, i, color, cycle_at);
    free(color);
    return st;
}

/* supersede 간선 from → to를 따라 DFS. 간선 목록을 선형 탐색한다(규칙 수십 규모). */
static rc_status rc__supersede_dfs(const rc_graph *g, size_t node,
                                   unsigned char *color, size_t *cycle_at) {
    color[node] = RC__GRAY;
    for (size_t e = 0; e < g->n_supersedes; e++) {
        if (g->supersedes[e].from_idx != node) continue;
        size_t nb = g->supersedes[e].to_idx;
        if (color[nb] == RC__GRAY) { *cycle_at = nb; return RC_ERR_CYCLE; }
        if (color[nb] == RC__WHITE) {
            rc_status st = rc__supersede_dfs(g, nb, color, cycle_at);
            if (st != RC_OK) return st;
        }
    }
    color[node] = RC__BLACK;
    return RC_OK;
}

/* supersede 간선 순환을 defeats와 독립으로 감지한다(상호 배제는 작성 오류). */
rc_status rc__detect_supersede_cycle(const rc_graph *g, size_t *cycle_at) {
    size_t n = g->n_rules;
    if (n == 0 || g->n_supersedes == 0) return RC_OK;
    unsigned char *color = calloc(n, 1);
    if (!color) return RC_ERR_NOMEM;
    rc_status st = RC_OK;
    for (size_t i = 0; i < n && st == RC_OK; i++)
        if (color[i] == RC__WHITE)
            st = rc__supersede_dfs(g, i, color, cycle_at);
    free(color);
    return st;
}
