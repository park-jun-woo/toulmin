#define _POSIX_C_SOURCE 200809L   /* strdup — 전 .c 공통 규약 */

#include "internal.h"

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* err_msg에 메시지만 포맷 기록(builder_err 미변경). */
void rc__set_error(rc_graph *g, rc_status st, const char *fmt, ...) {
    (void)st;
    if (!g) return;
    va_list ap;
    va_start(ap, fmt);
    vsnprintf(g->err_msg, sizeof g->err_msg, fmt, ap);
    va_end(ap);
}

/* 지연 빌더 오류: 세부 원인을 err_msg에 포맷하고, 첫 오류만 RC_ERR_BUILDER로 래치. */
void rc__builder_fail(rc_graph *g, rc_status st, const char *fmt, ...) {
    (void)st;
    if (!g) return;
    va_list ap;
    va_start(ap, fmt);
    vsnprintf(g->err_msg, sizeof g->err_msg, fmt, ap);
    va_end(ap);
    if (g->builder_err == RC_OK) g->builder_err = RC_ERR_BUILDER;
}

/* rules 배열이 need개를 담도록 성장한다. */
bool rc__rules_reserve(rc_graph *g, size_t need) {
    if (need <= g->cap_rules) return true;
    size_t new_cap = g->cap_rules ? g->cap_rules * 2 : 8;
    while (new_cap < need) new_cap *= 2;
    rc_rule_meta *grown = realloc(g->rules, new_cap * sizeof *grown);
    if (!grown) { rc__set_error(g, RC_ERR_NOMEM, "out of memory"); return false; }
    g->rules = grown;
    g->cap_rules = new_cap;
    return true;
}

/* handles 배열에 핸들 포인터를 추가한다(graph가 소유). */
bool rc__handles_append(rc_graph *g, rc_rule *r) {
    if (g->n_handles == g->cap_handles) {
        size_t new_cap = g->cap_handles ? g->cap_handles * 2 : 8;
        rc_rule **grown = realloc(g->handles, new_cap * sizeof *grown);
        if (!grown) { rc__set_error(g, RC_ERR_NOMEM, "out of memory"); return false; }
        g->handles = grown;
        g->cap_handles = new_cap;
    }
    g->handles[g->n_handles++] = r;
    return true;
}

/* 규칙 메타의 specs 배열에 spec 포인터를 차용 추가한다. */
bool rc__specs_append(rc_graph *g, rc_rule_meta *m, rc_spec *spec) {
    if (m->n_specs == m->cap_specs) {
        size_t new_cap = m->cap_specs ? m->cap_specs * 2 : 4;
        rc_spec **grown = realloc(m->specs, new_cap * sizeof *grown);
        if (!grown) { rc__builder_fail(g, RC_ERR_NOMEM, "out of memory"); return false; }
        m->specs = grown;
        m->cap_specs = new_cap;
    }
    m->specs[m->n_specs++] = spec;
    return true;
}

/* defeats 배열이 need개를 담도록 성장한다. */
bool rc__defeats_reserve(rc_graph *g, size_t need) {
    if (need <= g->cap_defeats) return true;
    size_t new_cap = g->cap_defeats ? g->cap_defeats * 2 : 8;
    while (new_cap < need) new_cap *= 2;
    rc_defeat_edge *grown = realloc(g->defeats, new_cap * sizeof *grown);
    if (!grown) { rc__set_error(g, RC_ERR_NOMEM, "out of memory"); return false; }
    g->defeats = grown;
    g->cap_defeats = new_cap;
    return true;
}

rc_graph *rc_graph_create(const char *name) {
    rc_graph *g = calloc(1, sizeof *g);
    if (!g) return NULL;
    if (name) {
        g->name = strdup(name);
        if (!g->name) { free(g); return NULL; }
    }
    g->builder_err = RC_OK;
    return g;
}

void rc_graph_destroy(rc_graph *g) {
    if (!g) return;
    for (size_t i = 0; i < g->n_rules; i++) {
        free(g->rules[i].name);
        /* short_name은 name 내부를 가리키는 비소유 포인터 — free 금지 */
        free(g->rules[i].specs);   /* 배열만 — spec 자체는 호출자 소유 */
    }
    free(g->rules);
    for (size_t i = 0; i < g->n_handles; i++) free(g->handles[i]);
    free(g->handles);
    free(g->defeats);
    for (size_t i = 0; i < g->n_funcs; i++) free(g->funcs[i].id);
    free(g->funcs);
    free(g->name);
    free(g);
}

const char *rc_graph_last_error(const rc_graph *g) {
    if (!g) return NULL;
    return g->err_msg;
}

void rc_graph_set_evidence_destructor(rc_graph *g, void (*dtor)(void *)) {
    if (!g) return;
    g->evidence_dtor = dtor;
}

/* 규칙 메타의 함수 포인터(bool 또는 verdict)를 정체성 키용 범용 포인터로. */
rc__any_fn rc__meta_fn(const rc_rule_meta *m) {
    return m->fn ? (rc__any_fn)m->fn : (rc__any_fn)m->vfn;
}

/* 규칙을 등록한다. fn(bool)과 vfn(verdict) 중 정확히 하나만 non-NULL로 넘긴다.
   실패(NULL 인자/중복/메모리) 시 NULL + last_error 기록. */
static rc_rule *rc__register(rc_graph *g, rc_rule_func fn, rc_verdict_func vfn,
                             void *user_data, rc_strength strength, const char *role) {
    if (!g) return NULL;
    if (!fn && !vfn) { rc__set_error(g, RC_ERR_INVALID_ARG, "fn must not be NULL"); return NULL; }

    rc__any_fn key = fn ? (rc__any_fn)fn : (rc__any_fn)vfn;
    char *rid = rc__rule_id(g, key, user_data, NULL, 0);
    if (!rid) { rc__set_error(g, RC_ERR_NOMEM, "out of memory"); return NULL; }

    /* 중복 검출: 같은 (fn, user_data) 쌍 → 같은 funcID → 이름 충돌 */
    for (size_t i = 0; i < g->n_rules; i++) {
        if (strcmp(g->rules[i].name, rid) == 0) {
            rc__set_error(g, RC_ERR_DUPLICATE, "duplicate rule registration: %s", rid);
            free(rid);
            return NULL;
        }
    }

    if (!rc__rules_reserve(g, g->n_rules + 1)) { free(rid); return NULL; }
    size_t idx = g->n_rules++;
    rc_rule_meta *m = &g->rules[idx];
    memset(m, 0, sizeof *m);
    m->name = rid;
    m->short_name = rc__short_name(rid);   /* name 내부 포인터(비소유) */
    m->qualifier = 1.0;
    m->strength = strength;
    m->role = role;                         /* 정적 문자열 */
    m->fn = fn;
    m->vfn = vfn;
    m->user_data = user_data;

    rc_rule *r = malloc(sizeof *r);
    if (!r) {
        rc__set_error(g, RC_ERR_NOMEM, "out of memory");
        g->n_rules--;         /* 방금 예약한 메타 롤백 */
        free(rid);
        return NULL;
    }
    r->graph = g;
    r->idx = idx;
    if (!rc__handles_append(g, r)) {   /* NOMEM → set_error 내부 기록 */
        free(r);
        g->n_rules--;
        free(rid);
        return NULL;
    }
    return r;
}

/* bool 규칙: true → +1, false → -1 (평가 시 변환). */
rc_rule *rc_graph_rule(rc_graph *g, rc_rule_func fn, void *user_data) {
    return rc__register(g, fn, NULL, user_data, RC_DEFEASIBLE, "rule");
}

rc_rule *rc_graph_counter(rc_graph *g, rc_rule_func fn, void *user_data) {
    return rc__register(g, fn, NULL, user_data, RC_DEFEASIBLE, "counter");
}

rc_rule *rc_graph_except(rc_graph *g, rc_rule_func fn, void *user_data) {
    return rc__register(g, fn, NULL, user_data, RC_DEFEATER, "except");
}

/* verdict 규칙: 규칙 함수가 v ∈ [-1,1]을 직접 낸다. */
rc_rule *rc_graph_rule_verdict(rc_graph *g, rc_verdict_func fn, void *user_data) {
    return rc__register(g, NULL, fn, user_data, RC_DEFEASIBLE, "rule");
}

rc_rule *rc_graph_counter_verdict(rc_graph *g, rc_verdict_func fn, void *user_data) {
    return rc__register(g, NULL, fn, user_data, RC_DEFEASIBLE, "counter");
}

rc_rule *rc_graph_except_verdict(rc_graph *g, rc_verdict_func fn, void *user_data) {
    return rc__register(g, NULL, fn, user_data, RC_DEFEATER, "except");
}
