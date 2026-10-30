#ifndef RULECAT_H
#define RULECAT_H

#include <stdbool.h>
#include <stddef.h>

/* 라이브러리 버전 (SemVer 문자열) */
#define RC_VERSION "0.2.0"

#ifdef __cplusplus
extern "C" {
#endif

/* --- 상태 코드 (예외 대체) --- */
typedef enum rc_status {
    RC_OK = 0,
    RC_ERR_INVALID_ARG = 1,   /* NULL/범위 위반 인자 */
    RC_ERR_DUPLICATE   = 2,   /* 중복 규칙 등록 */
    RC_ERR_CYCLE       = 3,   /* defeats 그래프 순환 */
    RC_ERR_RULE        = 4,   /* 규칙 함수가 실패 반환 */
    RC_ERR_BUILDER     = 5,   /* 빌더 체인 중 지연 기록된 오류 */
    RC_ERR_HANDLER     = 6,   /* RunOn 핸들러 실패 (Phase 005) */
    RC_ERR_RUN_CYCLE   = 7,   /* 실행 합성 순환 (Phase 005) */
    RC_ERR_RUN_DEPTH   = 8,   /* 실행 깊이 초과 (Phase 005) */
    RC_ERR_NOMEM       = 9    /* 메모리 할당 실패 */
} rc_status;

/* --- Context: Go interface → vtable 구조체 --- */
typedef struct rc_context rc_context;
struct rc_context {
    void *(*get)(rc_context *self, const char *key); /* 없으면 NULL */
    void  (*set)(rc_context *self, const char *key, void *value);
    void  *impl; /* 구현체 내부 상태 */
};

/* 기본 구현: 연결 리스트 기반 맵 (규칙 수십 개 규모엔 충분) */
rc_context *rc_map_context_create(void);
void        rc_map_context_destroy(rc_context *ctx);

/* --- Spec: Go interface → vtable 구조체 --- */
typedef struct rc_spec rc_spec;
struct rc_spec {
    const char *(*name)(const rc_spec *self);
    rc_status   (*validate)(const rc_spec *self);
    /* 규칙 ID 조합용 정체성 문자열 (Python json.dumps(__dict__) 대응).
       NULL 허용 — NULL이면 name() 결과를 사용 */
    const char *(*id)(const rc_spec *self);
    void *data;
};

rc_spec *rc_find_spec(rc_spec *const *specs, size_t n_specs, const char *name);

/* --- RuleFunc: 클로저 부재 → 함수 포인터 + user_data --- */
/* bool 규칙 함수(옛 시그니처, 그대로 지원). *out_activated가 true면 verdict +1,
   false면 -1로 바뀐다 — bool 그래프의 판정값은 이전과 같다(Except가 공격하는 경우만 다름). */
typedef rc_status (*rc_rule_func)(
    rc_context *ctx,
    rc_spec *const *specs, size_t n_specs,
    void *user_data,
    bool  *out_activated,
    void **out_evidence);   /* NULL 가능. 소유권은 아래 규약 */

/* verdict 규칙 함수: 규칙이 얼마나 해당하는지 *out_verdict ∈ [-1, +1]로 낸다.
   +1 = 완전히 해당, 0 = 반쯤 해당, -1 = 해당 없음.
   주의: 해당 없음은 0이 아니라 -1이다. 0은 '반쯤 해당'이다(엔진은 이 실수를 검출하지 못한다).
   NaN이나 범위 밖 값은 평가 오류(RC_ERR_RULE)다. v > -1이면 Active, v == -1이면 공격자를
   평가하지 않는다(지연 평가). */
typedef rc_status (*rc_verdict_func)(
    rc_context *ctx,
    rc_spec *const *specs, size_t n_specs,
    void *user_data,
    double *out_verdict,
    void  **out_evidence);  /* NULL 가능. 소유권은 아래 규약 */

/* --- Enums --- */
typedef enum rc_strength {
    RC_DEFEASIBLE = 0,
    RC_STRICT     = 1,
    RC_DEFEATER   = 2
} rc_strength;

typedef enum rc_eval_method { RC_MATRIX = 0 } rc_eval_method;

/* --- Options --- */
typedef struct rc_eval_option {
    rc_eval_method method;
    bool trace;
    bool duration;
} rc_eval_option; /* evaluate/run에 NULL 전달 시 {RC_MATRIX, false, false} */

/* --- Trace --- */
typedef struct rc_trace_entry {
    const char *name;        /* = Claim (short name, graph 소유 문자열) */
    const char *role;        /* "rule" | "counter" | "except" (정적 문자열) */
    bool        activated;
    double      qualifier;
    double      verdict;     /* 노드 우세/패배 판별용 */
    void       *evidence;
    rc_context *ground;      /* ctx 그대로 (모든 entry 동일 포인터) */
    rc_spec *const *specs;   /* = Backing (graph 보관 배열 차용) */
    size_t      n_specs;
    double      duration_ms; /* duration=false면 0.0 */
    double      rule_verdict;/* 규칙 함수 자신의 v ∈ [-1,1] (bool 규칙은 ±1). 실행 안 된 노드는 -1.
                                activated == (rule_verdict > -1). 0.2.0에서 끝에 추가 */
} rc_trace_entry;

/* --- Results --- */
typedef struct rc_eval_result {
    const char     *name;    /* short name (graph 소유 문자열) */
    double          verdict;
    void           *evidence;
    rc_trace_entry *trace;   /* opt.trace=true일 때만, 아니면 NULL */
    size_t          n_trace;
    double          applicability; /* 유효 적용도 a_eff ∈ [0,1]: (v+1)/2에 예외 게이트
                                      Π(1 − F·(v_j+1)/2)를 곱한 값. 0.2.0에서 끝에 추가 */
} rc_eval_result;

typedef struct rc_results rc_results; /* opaque 컨테이너 */
const rc_eval_result *rc_results_items(const rc_results *r, size_t *out_len);
void                  rc_results_destroy(rc_results *r);

/* --- Graph / Rule (Phase 002에서 API 확정) --- */
typedef struct rc_graph rc_graph;
typedef struct rc_rule  rc_rule;

/* graph 생성/해제 */
rc_graph *rc_graph_create(const char *name);
void      rc_graph_destroy(rc_graph *g);

/* 등록: 실패(NULL 인자/중복/메모리) 시 NULL 반환 + last_error 기록 */
rc_rule *rc_graph_rule(rc_graph *g, rc_rule_func fn, void *user_data);
rc_rule *rc_graph_counter(rc_graph *g, rc_rule_func fn, void *user_data);
rc_rule *rc_graph_except(rc_graph *g, rc_rule_func fn, void *user_data);

/* verdict 규칙 등록(0.2.0). 역할·실패 규약은 위 bool 판과 같다. */
rc_rule *rc_graph_rule_verdict(rc_graph *g, rc_verdict_func fn, void *user_data);
rc_rule *rc_graph_counter_verdict(rc_graph *g, rc_verdict_func fn, void *user_data);
rc_rule *rc_graph_except_verdict(rc_graph *g, rc_verdict_func fn, void *user_data);

/* 빌더: attacks는 상태 코드 반환, with/qualifier는 self 반환(오류 지연) */
rc_status   rc_rule_attacks(rc_rule *r, rc_rule *target);           /* force = 1 */
/* 공격 강도 force ∈ [0,1]을 준 간선(0.2.0). force는 반박 분모 항과 예외 게이트 양쪽에 곱해진다.
   NaN·범위 밖이면 RC_ERR_INVALID_ARG를 반환하고, 빌더 오류도 래치해 이후 evaluate/run이
   RC_ERR_BUILDER로 실패한다(Qualifier 관례, Go는 panic). 간선은 추가되지 않는다. */
rc_status   rc_rule_attacks_with_force(rc_rule *r, rc_rule *target, double force);
rc_rule    *rc_rule_with(rc_rule *r, rc_spec *spec);
rc_rule    *rc_rule_qualifier(rc_rule *r, double q);   /* NaN·[0,1] 밖은 지연 빌더 오류 */
const char *rc_rule_id(const rc_rule *r);   /* 현재 규칙 ID 조회 */

/* 마지막 오류 메시지 (graph 소유 버퍼, 다음 오류 시 덮어씀) */
const char *rc_graph_last_error(const rc_graph *g);

/* evidence 소멸자 등록 (선택). 등록 시 결과 객체 해제 때 evidence마다 1회 호출 */
void rc_graph_set_evidence_destructor(rc_graph *g, void (*dtor)(void *));

/* 그래프를 평가한다(h-Categoriser). opt=NULL이면 기본 {RC_MATRIX,false,false}.
     a = (v+1)/2,  S = Σ_반박 F·(v_j+1)/2,  a_eff = a·Π_예외 (1 − F·(v_j+1)/2)
     raw = q·a_eff/(1+S),  verdict = 2·raw − 1    (Strict는 반박·예외를 모두 무시)
   예외(Except) 여부는 공격자 노드의 strength(RC_DEFEATER) 또는 등록 role("except")로 판정한다.
   규칙 함수가 실패(RC_OK 외)하거나 범위 밖·NaN verdict를 내면 RC_ERR_RULE.
   성공: RC_OK + *out에 rc_results (호출자가 rc_results_destroy, graph보다 먼저 해제).
   실패: rc_status 반환 + rc_graph_last_error(g)에 메시지, *out은 NULL. */
rc_status rc_graph_evaluate(rc_graph *g, rc_context *ctx,
                            const rc_eval_option *opt, rc_results **out);

/* --- Run 엔진 (Phase 005) --- */

/* Trace: 한 Run의 읽기 전용 뷰 — 등록 순서 전 노드 엔트리 + 이 Run의 ctx.
   (Go []TraceEntry + Trace, Py Trace 클래스 대응) */
typedef struct rc_trace rc_trace;   /* opaque */
const rc_trace_entry *rc_trace_all(const rc_trace *t, size_t *out_len); /* 등록 순서 전 노드 */
const rc_trace_entry *rc_trace_get(const rc_trace *t, const char *name); /* short name, 없으면 NULL */
rc_context           *rc_trace_ctx(const rc_trace *t);

/* NodeHandler: 클로저 부재 → 함수 포인터 + user_data. Active 노드에서 발화한다.
   self는 발화 노드의 entry, trace는 전 노드 읽기 전용 뷰. ctx는 별도 인자가 아니라
   rc_trace_ctx(trace)로 접근한다(Go func(self TraceEntry, t Trace) error와 동형).
   RC_OK 외를 반환하면 Run이 RC_ERR_HANDLER로 중단된다. */
typedef rc_status (*rc_node_handler)(
    const rc_trace_entry *self,
    const rc_trace       *trace,
    void                 *user_data);

/* 빌더 (체이닝): Active 발화 핸들러 등록 / 실행 간선(하위 그래프) 선언.
   sub는 차용 — 소유·수명은 호출자 책임이며 부모 graph destroy가 sub를 해제하지 않는다.
   sub == NULL이면 지연 빌더 오류. */
rc_rule *rc_rule_run_on(rc_rule *r, rc_node_handler h, void *user_data);
rc_rule *rc_rule_run(rc_rule *r, rc_graph *sub);

/* RunResult: 한 Run의 results + trace를 일괄 소유한다(단일 destroy로 함께 해제). */
typedef struct rc_run_result rc_run_result;   /* opaque */
const rc_results *rc_run_result_results(const rc_run_result *rr);
const rc_trace   *rc_run_result_trace(const rc_run_result *rr);
void              rc_run_result_destroy(rc_run_result *rr);

/* 그래프를 Run한다: full pass 평가 → 전 노드 평면 trace 조립 → 등록 순서로 Active
   (activated(규칙 v > -1) && verdict > 0) 노드만 (a) 핸들러 → (b) 하위 그래프 순으로 디스패치.
   opt.trace/opt.duration은 내부에서 강제 off. 하위 그래프 verdict는 격리(즉시 폐기)된다.
   성공: RC_OK + *out. 핸들러/하위 Run 오류 시에도 *out은 부분 결과로 유효할 수 있으니
   *out이 NULL이 아니면 항상 rc_run_result_destroy로 해제한다. 순환/깊이/빌더 오류 등
   디스패치 이전 실패 시 *out은 NULL. */
rc_status rc_graph_run(rc_graph *g, rc_context *ctx,
                       const rc_eval_option *opt, rc_run_result **out);

/* --- 테스트 헬퍼: 기대값 (클로저 대체 값 구조체) --- */
typedef enum rc_expect_kind {
    RC_EXPECT_ABOVE,     /* verdict > lo */
    RC_EXPECT_AT_MOST,   /* verdict <= lo */
    RC_EXPECT_BETWEEN,   /* lo < verdict <= hi */
    RC_EXPECT_NO_RESULT  /* 결과 0건 */
} rc_expect_kind;

typedef struct rc_expectation {
    rc_expect_kind kind;
    double lo, hi;
} rc_expectation;

rc_expectation rc_verdict_above(double v);
rc_expectation rc_verdict_at_most(double v);
rc_expectation rc_verdict_between(double lo, double hi);
rc_expectation rc_no_result(void);

/* 첫 결과의 verdict를 기대값에 대조한다.
   통과: RC_OK. 실패: RC_ERR_INVALID_ARG + errbuf에 사유(errbuf NULL 허용). */
rc_status rc_expect_check(const rc_results *results, rc_expectation e,
                          char *errbuf, size_t errlen);

/* --- 테스트 헬퍼: 테이블 주도 케이스 --- */
typedef struct rc_test_case {
    const char           *name;
    rc_context           *ctx;    /* NULL이면 케이스마다 임시 map_context 생성/해제 */
    const rc_eval_option *opt;    /* NULL이면 기본 옵션 */
    rc_expectation        expect;
} rc_test_case;

/* 케이스를 순서대로 실행, 첫 실패에서 중단.
   반환 코드 + errbuf에 "[케이스명] 사유"(errbuf NULL 허용). */
rc_status rc_run_cases(rc_graph *g, const rc_test_case *cases, size_t n,
                       char *errbuf, size_t errlen);

#ifdef __cplusplus
}
#endif
#endif /* RULECAT_H */
