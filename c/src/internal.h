#ifndef RULECAT_INTERNAL_H
#define RULECAT_INTERNAL_H

#include <rulecat/rulecat.h>

#include <stdarg.h>
#include <stddef.h>

/* --- 함수 정체성 레지스트리 항목 --- */
/* (fn, user_data) 쌍마다 유일한 ID 문자열을 발급한다.
   C에는 클로저가 없으므로 규칙 정체성 키는 (fn, user_data) 쌍이다. */
typedef struct rc_func_entry {
    rc_rule_func fn;
    void        *user_data;
    char        *id;          /* graph 소유(strdup), rc_graph_destroy에서 해제 */
} rc_func_entry;

/* --- 규칙 메타데이터 --- */
/* rules 배열의 각 항목. 사용자에게 반환하는 rc_rule 핸들은 이 배열의 인덱스만 담는다. */
typedef struct rc_rule_meta {
    char        *name;        /* full rule id — graph 소유(strdup) */
    const char  *short_name;  /* name 내부를 가리키는 비소유 포인터(rc__short_name 결과) */
    double       qualifier;   /* 기본 1.0 */
    rc_strength  strength;
    const char  *role;        /* "rule" | "counter" | "except" — 정적 문자열, free 금지 */
    rc_spec    **specs;       /* 포인터 차용 배열 (배열 자체는 graph 소유) */
    size_t       n_specs, cap_specs;
    rc_rule_func fn;
    void        *user_data;
    /* Phase 005 Run 엔진: Active 노드 핸들러 + 실행 간선(하위 그래프, 차용). */
    rc_node_handler run_on;       /* Active 발화 핸들러 (NULL 허용) */
    void           *run_on_data;  /* 핸들러 user_data (클로저 대체) */
    rc_graph       *run_graph;    /* 실행 간선 하위 그래프 — 차용(호출자 소유) */
} rc_rule_meta;

/* --- defeats 간선 (인덱스 기반) --- */
/* Go/TS/Py는 이름 문자열로 저장하지만, C 핸들은 안정적 인덱스를 가지므로
   (from_idx, to_idx)로 저장한다. with의 ID 변경 시 재바인딩이 불필요하다. */
typedef struct rc_defeat_edge { size_t from_idx, to_idx; } rc_defeat_edge;

/* --- 사용자 반환 핸들 (개별 malloc, 인덱스만 담음) --- */
struct rc_rule { rc_graph *graph; size_t idx; };

/* --- 그래프 내부 구조체 --- */
struct rc_graph {
    char           *name;

    rc_rule_meta   *rules;        /* 규칙 메타 배열 */
    size_t          n_rules, cap_rules;

    rc_rule       **handles;      /* 반환한 rc_rule* 소유 (개별 malloc) */
    size_t          n_handles, cap_handles;

    rc_defeat_edge *defeats;      /* defeats 간선 */
    size_t          n_defeats, cap_defeats;

    rc_func_entry  *funcs;        /* (fn, user_data) → id 선형 레지스트리 */
    size_t          n_funcs;
    size_t          cap_funcs;
    unsigned        func_counter; /* 발급 카운터 — 유일성 보장 */

    rc_status       builder_err;  /* 지연 빌더 오류 (첫 오류 보존) */
    char            err_msg[256]; /* rc_graph_last_error 고정 버퍼 */
    void          (*evidence_dtor)(void *); /* 선택적 evidence 소멸자 */
};

/* --- 평가 상태 (인덱스 기반, Go/TS/Py의 name→map 대응) --- */
/* 규칙 인덱스가 안정적이므로 해시맵 대신 인덱스 배열을 쓴다(O(1), 의미 동일). */
typedef struct rc_eval_state {
    rc_graph   *g;
    rc_context *ctx;
    /* 간선(인덱스 기반): attackers[i] = i를 공격하는 노드 인덱스 배열 */
    size_t    **attackers;
    size_t     *n_attackers;
    bool       *is_attacker;    /* attacker set: 누군가를 공격하는 노드 */
    /* 노드별 평가 상태 */
    bool       *ran, *active, *has_verdict;
    double     *verdict;
    void      **evidence;
    /* trace 버퍼 (opt.trace 경로 전용, 워런트별 reset). 소멸은 rc__eval_state_free. */
    rc_trace_entry *trace;
    size_t          n_trace, cap_trace;
    rc_status   err;            /* 첫 오류 래치 */
} rc_eval_state;

/* --- Run 엔진 내부 구조 (Phase 005) --- */

/* rc_trace 실체: 등록 순서 전 노드 엔트리(소유 배열) + ctx(차용). */
struct rc_trace {
    rc_trace_entry *nodes;   /* trace 소유 배열 (문자열/spec/evidence는 차용) */
    size_t          len;
    rc_context     *ctx;     /* 차용 */
};

/* rc_run_result 실체: results + trace 일괄 소유(단일 destroy). */
struct rc_run_result {
    rc_results *results;   /* 소유. evidence 풀(owned_evidence)도 여기 소속 */
    rc_trace   *trace;     /* 소유 */
};

/* 실행 합성 순환 감지용 (graph 포인터 identity 키잉) 동적 색칠 집합. */
typedef struct rc__gcolor     { const rc_graph *g; unsigned char color; } rc__gcolor;
typedef struct rc__gcolor_set { rc__gcolor *items; size_t len, cap; }     rc__gcolor_set;

/* RunGraph 간선을 DFS로 순회해 실행 합성 순환을 감지한다(최상위 1회, DAG 검증).
   순환 시 RC_ERR_RUN_CYCLE + *cycle_at에 순환 그래프. 할당 실패 시 RC_ERR_NOMEM. */
rc_status rc__run_cycle_dfs(const rc_graph *g, rc__gcolor_set *set,
                            const rc_graph **cycle_at);
/* 색칠 집합의 items 배열을 해제한다. */
void rc__gcolor_set_free(rc__gcolor_set *set);

/* rc_graph_evaluate/rc_graph_run 공용 평가 본체. full=true면 미평가 노드까지 rc__calc으로
   채운다(Go fillAll 대응). 성공 시 *out_res에 results, *es에 (호출자 소유) 평가 상태를
   남긴다 — 호출자가 evidence 이관 후 rc__eval_state_free. 실패 시 내부에서 es를 정리한다. */
rc_status rc__evaluate_internal(rc_graph *g, rc_context *ctx, rc_eval_option o,
                                bool full, rc_results **out_res, rc_eval_state *es);

/* n개 노드 슬롯을 가진 trace를 만든다(nodes calloc, ctx 차용). 실패 시 NULL. */
rc_trace *rc__trace_create(size_t n, rc_context *ctx);
/* trace의 nodes 배열과 컨테이너만 해제한다(문자열/spec/evidence는 차용이라 손대지 않음). */
void      rc__trace_destroy_internal(rc_trace *t);

/* rc_eval_state 배열을 구성한다(간선·attacker set 포함). 실패 시 RC_ERR_NOMEM. */
rc_status rc__eval_state_init(rc_eval_state *es, rc_graph *g, rc_context *ctx);

/* rc_eval_state의 배열만 해제한다(evidence 소유권은 results가 담당). */
void rc__eval_state_free(rc_eval_state *es);

/* WHITE/GRAY/BLACK DFS로 defeats 그래프 순환을 감지한다.
   순환 시 RC_ERR_CYCLE + *cycle_at에 순환 노드 인덱스. 할당 실패 시 RC_ERR_NOMEM. */
rc_status rc__detect_cycle(rc_eval_state *es, size_t *cycle_at);

/* opt가 NULL이면 기본값 {RC_MATRIX,false,false}, duration=true면 trace를 함의한다. */
rc_eval_option rc__resolve_option(const rc_eval_option *opt);

/* --- rc_results 컨테이너 헬퍼 (src/results.c) --- */
/* graph의 evidence_dtor를 연결한 빈 결과 컨테이너를 만든다. 실패 시 NULL. */
rc_results *rc__results_create(rc_graph *g);

/* 결과 항목을 추가한다. 성장 실패 시 false. */
bool rc__results_append(rc_results *r, rc_eval_result item);

/* dtor 등록 시 이 패스의 evidence를 노드당 1개 결과로 이관한다(소유권 획득). */
void rc__results_adopt_evidence(rc_results *r, rc_eval_state *es);

/* 평가 실패 시 부분 결과와 이미 생산된 evidence를 정리하고 컨테이너를 해제한다. */
void rc__results_abandon(rc_results *r, rc_eval_state *es);

/* --- 내부 헬퍼 (비공개, rc__ 접두사) --- */

/* (fn, user_data) 쌍의 정체성 ID를 조회하거나 새로 발급한다.
   반환 문자열은 graph 소유 — 호출자는 free 금지. 할당 실패 시 NULL. */
const char *rc__func_id(rc_graph *g, rc_rule_func fn, void *user_data);

/* funcID + specs 직렬화로 규칙 ID를 만든다.
   반환 문자열은 새로 할당된 것으로 호출자가 소유(free 책임). 실패 시 NULL. */
char *rc__rule_id(rc_graph *g, rc_rule_func fn, void *user_data,
                  rc_spec *const *specs, size_t n_specs);

/* full 경로 문자열에서 짧은 이름을 추출한다.
   반환값은 full 내부를 가리키는 포인터(할당 없음) — full보다 오래 살 수 없다. */
const char *rc__short_name(const char *full);

/* err_msg에 오류 메시지를 포맷 기록한다. builder_err는 건드리지 않는다. */
void rc__set_error(rc_graph *g, rc_status st, const char *fmt, ...);

/* CLOCK_MONOTONIC 기반 ms 단위 단조 시계 (TS performance.now() 대응). */
double rc__now_ms(void);

/* buf가 NULL이거나 len==0이면 아무것도 안 함. 아니면 vsnprintf로 포맷. */
void rc__fmt(char *buf, size_t len, const char *fmt, ...);

/* 지연 빌더 오류: err_msg에 세부 원인을 포맷하고, builder_err가 RC_OK일 때만
   RC_ERR_BUILDER로 래치한다(첫 오류 보존). st는 메시지용 세부 코드일 뿐. */
void rc__builder_fail(rc_graph *g, rc_status st, const char *fmt, ...);

/* 배열 성장 헬퍼. 실패 시 false 반환 + rc__set_error(NOMEM) 기록. */
bool rc__rules_reserve(rc_graph *g, size_t need);
bool rc__handles_append(rc_graph *g, rc_rule *r);
bool rc__specs_append(rc_graph *g, rc_rule_meta *m, rc_spec *spec);
bool rc__defeats_reserve(rc_graph *g, size_t need);

#endif /* RULECAT_INTERNAL_H */
