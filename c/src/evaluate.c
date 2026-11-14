#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <math.h>
#include <stdlib.h>
#include <string.h>

/* n바이트 딥카피. n==0이면 NULL. 실패 시 NULL. */
static void *rc__memdup(const void *src, size_t n) {
    if (n == 0) return NULL;
    void *dst = malloc(n);
    if (!dst) return NULL;
    memcpy(dst, src, n);
    return dst;
}

/* opt가 NULL이면 기본값. duration은 trace를 함의(Go/TS/Py와 동일). */
rc_eval_option rc__resolve_option(const rc_eval_option *opt) {
    rc_eval_option o = { RC_MATRIX, false, false };
    if (opt) o = *opt;
    if (o.duration) o.trace = true;
    return o;
}

/* Evaluate trace의 role은 등록 role이 아니라 defeats 구조로 판정한다
   (Go inferRole / Py infer_role 동형). 등록 role은 Run trace(Phase 005)에서 사용. */
static const char *rc__infer_role(const rc_eval_state *es, size_t i) {
    if (es->g->rules[i].strength == RC_DEFEATER) return "except";
    if (es->is_attacker[i]) return "counter";
    return "rule";
}

/* trace 버퍼에 엔트리를 추가한다. 성장 실패 시 false. */
static bool rc__trace_append(rc_eval_state *es, rc_trace_entry e) {
    if (es->n_trace == es->cap_trace) {
        size_t nc = es->cap_trace ? es->cap_trace * 2 : 4;
        rc_trace_entry *grown = realloc(es->trace, nc * sizeof *grown);
        if (!grown) return false;
        es->trace = grown;
        es->cap_trace = nc;
    }
    es->trace[es->n_trace++] = e;
    return true;
}

/* 규칙 함수를 원래 타입으로 호출해 자기 verdict v ∈ [-1,1]과 evidence를 얻는다.
   bool 규칙은 true → +1, false → -1로 바꾼다(Go boolRuleFunc 대응).
   실패 반환 또는 범위 밖·NaN verdict면 es->err를 RC_ERR_RULE로 래치한다(Go safeCall 대응).
   범위 오류일 때 함수가 낸 evidence는 es->evidence[i]에 남겨 정리 경로
   (rc__results_abandon)가 dtor로 회수하게 한다. */
static rc_status rc__call_rule(rc_eval_state *es, size_t i, double *out_v, void **out_ev) {
    rc_rule_meta *m = &es->g->rules[i];
    void *ev = NULL;
    double v = -1.0;
    rc_status st;
    if (m->fn) {
        bool act = false;
        st = m->fn(es->ctx, m->specs, m->n_specs, m->user_data, &act, &ev);
        v = act ? 1.0 : -1.0;
    } else {
        st = m->vfn(es->ctx, m->specs, m->n_specs, m->user_data, &v, &ev);
    }
    if (st != RC_OK) {
        rc__set_error(es->g, RC_ERR_RULE,
                      "rule \"%s\": failed (status %d)", m->name, st);
        es->err = RC_ERR_RULE;
        return RC_ERR_RULE;
    }
    if (isnan(v) || v < -1.0 || v > 1.0) {
        es->evidence[i] = ev;
        rc__set_error(es->g, RC_ERR_RULE,
                      "rule \"%s\": rule verdict %g out of range [-1, 1]", m->name, v);
        es->err = RC_ERR_RULE;
        return RC_ERR_RULE;
    }
    *out_v = v;
    *out_ev = ev;
    return RC_OK;
}

/* 공격자가 예외(Except, defeater)로 게이트하는지: strength 또는 등록 role로 판정
   (Go isExceptAttacker 대응). 나머지 공격자는 분모에 들어가는 반박이다. */
static bool rc__is_except_attacker(const rc_eval_state *es, size_t j) {
    const rc_rule_meta *m = &es->g->rules[j];
    return m->strength == RC_DEFEATER || (m->role && strcmp(m->role, "except") == 0);
}

/* raw = q·a_eff/(1+sum) → verdict 2·raw − 1을 캐시하고 a_eff를 기록한다(Go settle). */
static double rc__settle(rc_eval_state *es, size_t i, double a_eff, double sum) {
    es->applicability[i] = a_eff;
    double raw = es->g->rules[i].qualifier * a_eff / (1.0 + sum);
    double v = 2.0 * raw - 1.0;
    es->has_verdict[i] = true;
    es->verdict[i] = v;
    return v;
}

/* h-Categoriser 핵심 — 재귀 + verdict 캐시(memoization).
   trace=true면 첫 방문 노드마다 엔트리를 방문 순서로 남기고, with_duration이면 규칙 함수
   실행 구간을 ms로 잰다. 규칙 v == -1(해당 없음)이면 공격자를 평가하지 않고 -1.0(캐시 미기록).
     a = (v+1)/2,  S = Σ_반박 F·(v_j+1)/2,  a_eff = a·Π_예외 (1 − F·(v_j+1)/2)
     raw = q·a_eff/(1+S),  verdict = 2·raw − 1
   Strict는 반박·예외를 모두 무시한다. bool 규칙·예외 없음·F=1이면 옛 식
   raw = q/(1+Σ raw(attackers))와 비트 단위로 같다. */
static double rc__calc_node(rc_eval_state *es, size_t i, bool trace, bool with_duration) {
    if (es->err != RC_OK) return -1.0;
    if (es->has_verdict[i]) return es->verdict[i];

    rc_rule_meta *m = &es->g->rules[i];

    if (!es->ran[i]) {
        es->ran[i] = true;
        double start = (trace && with_duration) ? rc__now_ms() : 0.0;
        double rv = -1.0;
        void *ev = NULL;
        if (rc__call_rule(es, i, &rv, &ev) != RC_OK) return -1.0;
        es->rule_verdict[i] = rv;
        es->active[i]       = rv > -1.0;
        es->evidence[i]     = ev;

        if (trace) {
            double elapsed = with_duration ? rc__now_ms() - start : 0.0;
            rc_trace_entry entry = {
                .name         = m->short_name,          /* graph 소유 문자열 차용 */
                .role         = rc__infer_role(es, i),  /* 구조 기반 판정 */
                .activated    = es->active[i],
                .qualifier    = m->qualifier,
                .verdict      = 0.0,                     /* 기록 시점엔 미확정 */
                .evidence     = ev,
                .ground       = NULL,                    /* Evaluate 경로에선 미기록 */
                .specs        = m->specs,               /* = Backing, 차용 */
                .n_specs      = m->n_specs,
                .duration_ms  = elapsed,
                .rule_verdict = rv,                      /* 규칙 자신의 v */
                .excluded     = false,                   /* 워런트 평가 후 배제 패스가 표시 */
            };
            if (!rc__trace_append(es, entry)) {
                rc__set_error(es->g, RC_ERR_NOMEM, "trace: out of memory");
                es->err = RC_ERR_NOMEM;
                return -1.0;
            }
        }
    }

    if (!es->active[i]) return -1.0;

    double a = (es->rule_verdict[i] + 1.0) / 2.0;
    if (m->strength == RC_STRICT) return rc__settle(es, i, a, 0.0);

    double sum = 0.0, gate = 1.0;
    for (size_t k = 0; k < es->n_attackers[i]; k++) {
        size_t j = es->attackers[i][k];
        double vj = rc__calc_node(es, j, trace, with_duration);
        if (es->err != RC_OK) return -1.0;
        double x = es->forces[i][k] * (vj + 1.0) / 2.0;
        if (rc__is_except_attacker(es, j)) gate *= 1.0 - x;   /* 예외: 곱셈 게이트 */
        else                               sum  += x;         /* 반박: 분모 */
    }
    return rc__settle(es, i, a * gate, sum);
}

/* trace 모드에서 워런트마다 상태를 비운다. evidence 배열은 0으로 클리어하되
   해제하지 않는다 — 소유권은 이미 results(owned_evidence)로 이관된 상태. */
static void rc__reset(rc_eval_state *es) {
    size_t n = es->g->n_rules;
    if (n) {
        memset(es->ran,         0, n * sizeof *es->ran);
        memset(es->active,      0, n * sizeof *es->active);
        memset(es->has_verdict, 0, n * sizeof *es->has_verdict);
        memset(es->rule_verdict,  0, n * sizeof *es->rule_verdict);
        memset(es->applicability, 0, n * sizeof *es->applicability);
        memset(es->evidence,    0, n * sizeof *es->evidence);
    }
    es->n_trace = 0;
    es->err = RC_OK;
}

/* 보고 전용 배제 패스: 배제자 중 하나라도 a_eff > 0이면 target을 es->excluded에 표시한다.
   a_eff는 예외 게이트를 반영하고 반박은 반영하지 않는다(Strict는 a_eff = a). 비활성 배제자는
   settle에 이르지 않아 a_eff 0이다. 배제자는 모두 평가(멱등)하고 단락하지 않는다. 배제된
   배제자도 배제권을 유지한다. 평가 중 오류는 es->err에 래치된다. */
static void rc__resolve_supersedes(rc_eval_state *es) {
    rc_graph *g = es->g;
    if (g->n_rules) memset(es->excluded, 0, g->n_rules * sizeof *es->excluded);
    for (size_t e = 0; e < g->n_supersedes; e++) {
        size_t from = g->supersedes[e].from_idx;
        rc__calc_node(es, from, false, false);
        if (es->err != RC_OK) return;
        if (es->applicability[from] > 0.0) es->excluded[g->supersedes[e].to_idx] = true;
    }
}

/* trace 버퍼 엔트리에 배제 표시를 옮긴다. 엔트리 name은 graph 소유 short_name을 차용하므로
   포인터 동일성으로 노드를 찾는다. */
static void rc__mark_trace_excluded(rc_eval_state *es) {
    for (size_t k = 0; k < es->n_trace; k++) {
        es->trace[k].excluded = false;
        for (size_t t = 0; t < es->g->n_rules; t++)
            if (es->excluded[t] && es->trace[k].name == es->g->rules[t].short_name)
                es->trace[k].excluded = true;
    }
}

/* Defeater가 아니고 아무도 공격하지 않는 노드가 warrant. */
static bool rc__is_warrant(const rc_eval_state *es, size_t i) {
    if (es->g->rules[i].strength == RC_DEFEATER) return false;
    return !es->is_attacker[i];
}

/* 공용 평가 본체 — rc_graph_evaluate(full=false)와 rc_graph_run(full=true)이 공유한다.
   성공 시 *out_res에 results, *es에 (호출자 소유) 평가 상태를 남긴다 — 호출자가 evidence를
   이관하고 rc__eval_state_free한다. 실패 시 내부에서 es와 부분 results를 모두 정리한다.
   builder_err 체크가 여기에 있어야 rc__run_depth가 실행하는 하위 그래프의 지연 빌더 오류도
   합성 경로에서 RC_ERR_BUILDER로 보고된다. */
rc_status rc__evaluate_internal(rc_graph *g, rc_context *ctx, rc_eval_option o,
                                bool full, rc_results **out_res, rc_eval_state *es) {
    *out_res = NULL;
    if (g->builder_err != RC_OK) return g->builder_err;   /* 지연 빌더 오류 */

    rc_status st = rc__eval_state_init(es, g, ctx);
    if (st != RC_OK) return st;

    size_t cycle_at = 0;
    st = rc__detect_cycle(es, &cycle_at);
    if (st != RC_OK) {
        if (st == RC_ERR_CYCLE)
            rc__set_error(g, RC_ERR_CYCLE, "cycle detected at %s",
                          g->rules[cycle_at].name);
        rc__eval_state_free(es);
        return st;
    }
    st = rc__detect_supersede_cycle(g, &cycle_at);
    if (st != RC_OK) {
        if (st == RC_ERR_CYCLE)
            rc__set_error(g, RC_ERR_CYCLE, "supersede cycle detected at %s",
                          g->rules[cycle_at].name);
        rc__eval_state_free(es);
        return st;
    }

    rc_results *res = rc__results_create(g);
    if (!res) { rc__eval_state_free(es); return RC_ERR_NOMEM; }

    for (size_t i = 0; i < g->n_rules; i++) {
        if (!rc__is_warrant(es, i)) continue;

        /* trace 모드: 직전 워런트 패스의 evidence를 results로 이관 후 상태 초기화. */
        if (o.trace) {
            rc__results_adopt_evidence(res, es);
            rc__reset(es);
        }

        double v = rc__calc_node(es, i, o.trace, o.duration);
        if (es->err != RC_OK) {
            st = es->err;
            rc__results_abandon(res, es);
            rc__eval_state_free(es);
            return st;
        }
        if (!es->active[i]) continue;

        rc_eval_result item = {
            .name     = g->rules[i].short_name,   /* graph 소유 문자열 차용 */
            .verdict  = v,
            .evidence = es->evidence[i],
            .trace    = NULL,
            .n_trace  = 0,
            .applicability = es->applicability[i],   /* a_eff */
        };
        if (o.trace) {
            rc__resolve_supersedes(es);   /* 워런트별(리셋된) 상태에서 배제 판정 */
            if (es->err != RC_OK) {
                st = es->err;
                rc__results_abandon(res, es);
                rc__eval_state_free(es);
                return st;
            }
            rc__mark_trace_excluded(es);
        }
        if (o.trace && es->n_trace > 0) {
            /* trace 버퍼를 결과에 딥카피 (엔트리 배열만 — 문자열/spec은 graph 차용). */
            item.trace = rc__memdup(es->trace, es->n_trace * sizeof *es->trace);
            if (!item.trace) {
                rc__results_abandon(res, es);
                rc__eval_state_free(es);
                return RC_ERR_NOMEM;
            }
            item.n_trace = es->n_trace;
        }
        if (!rc__results_append(res, item)) {
            free(item.trace);
            rc__results_abandon(res, es);
            rc__eval_state_free(es);
            return RC_ERR_NOMEM;
        }
    }

    if (full) {   /* run 경로: trace 강제 off이므로 non-trace 평가로 미평가 노드를 채움 */
        for (size_t i = 0; i < g->n_rules; i++)
            if (!es->ran[i]) rc__calc_node(es, i, false, false);
        if (es->err == RC_OK) rc__resolve_supersedes(es);   /* run trace의 excluded */
        if (es->err != RC_OK) {
            st = es->err;
            rc__results_abandon(res, es);
            rc__eval_state_free(es);
            return st;
        }
    }

    /* evidence 이관과 es 해제는 호출자 책임 (run은 trace 조립 후 이관해야 하므로). */
    *out_res = res;
    return RC_OK;
}

rc_status rc_graph_evaluate(rc_graph *g, rc_context *ctx,
                            const rc_eval_option *opt, rc_results **out) {
    if (out) *out = NULL;
    if (!g || !ctx || !out) return RC_ERR_INVALID_ARG;

    rc_eval_option o = rc__resolve_option(opt);

    rc_eval_state es;
    rc_results *res = NULL;
    rc_status st = rc__evaluate_internal(g, ctx, o, /*full=*/false, &res, &es);
    if (st != RC_OK) return st;   /* 실패 시 es는 내부에서 이미 정리됨 */

    rc__results_adopt_evidence(res, &es);   /* 마지막 패스 evidence 이관 */
    rc__eval_state_free(&es);
    *out = res;
    return RC_OK;
}
