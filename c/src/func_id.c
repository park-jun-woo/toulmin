#define _POSIX_C_SOURCE 200809L   /* strdup — 전 .c 공통 규약 */

#include "internal.h"

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>

/* 레지스트리에 (fn, user_data, id)를 추가한다. id 문자열은 strdup으로 복사해
   graph가 소유한다. 반환값은 저장된 id(graph 소유). 할당 실패 시 NULL. */
static const char *rc__func_registry_add(rc_graph *g, rc_rule_func fn,
                                         void *user_data, const char *id) {
    if (g->n_funcs == g->cap_funcs) {
        size_t new_cap = g->cap_funcs ? g->cap_funcs * 2 : 8;
        rc_func_entry *grown = realloc(g->funcs, new_cap * sizeof *grown);
        if (!grown) return NULL;
        g->funcs = grown;
        g->cap_funcs = new_cap;
    }
    char *copy = strdup(id);
    if (!copy) return NULL;
    g->funcs[g->n_funcs].fn = fn;
    g->funcs[g->n_funcs].user_data = user_data;
    g->funcs[g->n_funcs].id = copy;
    g->n_funcs++;
    return copy;
}

const char *rc__func_id(rc_graph *g, rc_rule_func fn, void *user_data) {
    for (size_t i = 0; i < g->n_funcs; i++) {
        if (g->funcs[i].fn == fn && g->funcs[i].user_data == user_data)
            return g->funcs[i].id;
    }
    char buf[64];
    /* 함수 포인터 → 정수 캐스팅은 ISO C 밖이지만 POSIX 보장. 주소+카운터로 유일성 확보.
       카운터가 유일성을 보장(주소 재사용 걱정 불필요), 주소는 디버깅 가독성용. */
    snprintf(buf, sizeof buf, "fn_%zx_%u",
             (size_t)(uintptr_t)fn, ++g->func_counter);
    return rc__func_registry_add(g, fn, user_data, buf);
}
