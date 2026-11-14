#define _POSIX_C_SOURCE 200809L   /* strdup — 전 .c 공통 규약 */

#include "internal.h"

#include <math.h>
#include <stdlib.h>

/* r → target defeats 간선을 강도 force로 추가한다(인덱스 기반, 검증은 호출자). */
static rc_status rc__attacks(rc_rule *r, rc_rule *target, double force) {
    rc_graph *g = r->graph;
    if (!rc__defeats_reserve(g, g->n_defeats + 1)) return RC_ERR_NOMEM;
    rc_defeat_edge e;
    e.from_idx = r->idx;
    e.to_idx = target->idx;
    e.force = force;
    g->defeats[g->n_defeats++] = e;
    return RC_OK;
}

/* r → target defeats 간선을 추가한다(force = 1, 기존 동작). */
rc_status rc_rule_attacks(rc_rule *r, rc_rule *target) {
    if (!r || !target || r->graph != target->graph) return RC_ERR_INVALID_ARG;
    return rc__attacks(r, target, 1.0);
}

/* 강도 force ∈ [0,1]의 defeats 간선을 추가한다. NaN·범위 밖이면 간선 없이
   RC_ERR_INVALID_ARG를 반환하고 빌더 오류를 래치한다(Qualifier 관례 — Go는 panic). */
rc_status rc_rule_attacks_with_force(rc_rule *r, rc_rule *target, double force) {
    if (!r || !target || r->graph != target->graph) return RC_ERR_INVALID_ARG;
    if (isnan(force) || force < 0.0 || force > 1.0) {
        rc__builder_fail(r->graph, RC_ERR_INVALID_ARG,
                         "force must be between 0.0 and 1.0, got %g", force);
        return RC_ERR_INVALID_ARG;
    }
    return rc__attacks(r, target, force);
}

/* r → target supersede 간선을 추가한다(인덱스 기반 — with의 ID 변경에도 재바인딩 불필요).
   r의 a_eff > 0일 때 target이 배제된다(판정은 rc__resolve_supersedes). */
rc_status rc_rule_supersedes(rc_rule *r, rc_rule *target) {
    if (!r || !target || r->graph != target->graph) return RC_ERR_INVALID_ARG;
    rc_graph *g = r->graph;
    if (!rc__supersedes_reserve(g, g->n_supersedes + 1)) return RC_ERR_NOMEM;
    g->supersedes[g->n_supersedes].from_idx = r->idx;
    g->supersedes[g->n_supersedes].to_idx = target->idx;
    g->supersedes[g->n_supersedes].force = 1.0;   /* 미사용 */
    g->n_supersedes++;
    return RC_OK;
}

/* spec을 규칙에 부착하고 ID를 재계산한다. 오류는 지연 처리(self 반환). */
rc_rule *rc_rule_with(rc_rule *r, rc_spec *spec) {
    if (!r) return r;
    rc_graph *g = r->graph;
    if (!spec) {
        rc__builder_fail(g, RC_ERR_INVALID_ARG, "with: spec must not be NULL");
        return r;
    }
    rc_status st = spec->validate(spec);
    if (st != RC_OK) {
        rc__builder_fail(g, st, "with: spec \"%s\" invalid", spec->name(spec));
        return r;
    }

    rc_rule_meta *m = &g->rules[r->idx];
    if (!rc__specs_append(g, m, spec)) return r;   /* NOMEM → builder_fail 내부 기록 */

    /* ID 재계산 (간선은 인덱스 기반이라 재바인딩 불필요) */
    char *new_id = rc__rule_id(g, rc__meta_fn(m), m->user_data, m->specs, m->n_specs);
    if (!new_id) { rc__builder_fail(g, RC_ERR_NOMEM, "out of memory"); return r; }
    free(m->name);
    m->name = new_id;
    /* short_name은 옛 name 내부를 가리켰던 비소유 포인터 — free 금지, 재지정만 */
    m->short_name = rc__short_name(new_id);
    return r;
}

/* qualifier [0,1] 설정. NaN·범위 위반은 지연 처리(self 반환). */
rc_rule *rc_rule_qualifier(rc_rule *r, double q) {
    if (!r) return r;
    if (isnan(q) || q < 0.0 || q > 1.0) {
        rc__builder_fail(r->graph, RC_ERR_INVALID_ARG,
                         "qualifier must be between 0.0 and 1.0, got %g", q);
        return r;
    }
    r->graph->rules[r->idx].qualifier = q;
    return r;
}

/* 현재 규칙의 full ID를 조회한다. r이 NULL이면 NULL. */
const char *rc_rule_id(const rc_rule *r) {
    if (!r) return NULL;
    return r->graph->rules[r->idx].name;
}

/* Active 발화 핸들러를 등록한다(이벤트는 Active 하나만). 오류 없이 self 반환. */
rc_rule *rc_rule_run_on(rc_rule *r, rc_node_handler h, void *user_data) {
    if (!r) return r;
    r->graph->rules[r->idx].run_on      = h;
    r->graph->rules[r->idx].run_on_data = user_data;
    return r;
}

/* 실행 간선(하위 그래프)을 선언한다. sub는 차용. NULL이면 지연 빌더 오류(self 반환). */
rc_rule *rc_rule_run(rc_rule *r, rc_graph *sub) {
    if (!r) return r;
    if (!sub) {
        rc__builder_fail(r->graph, RC_ERR_INVALID_ARG,
                         "run requires a non-NULL sub-graph");
        return r;
    }
    r->graph->rules[r->idx].run_graph = sub;
    return r;
}
