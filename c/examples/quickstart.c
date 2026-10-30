/* rulecat quick start — access-control 킬러 예제의 C 포트.
 *
 *   cc -std=c11 -Iinclude examples/quickstart.c librulecat.a -lm -o quickstart
 *   ./quickstart
 *
 * 규칙은 함수, 예외는 그래프 간선이다. DSL도, 사이드카도 없다.
 * verdict > 0 → allow, verdict <= 0 → deny (프레임워크가 해석한다).
 * 뒷부분은 등급 규칙(verdict 규칙)과 공격 강도(force) 예시다. */

#include <rulecat/rulecat.h>

#include <stdio.h>
#include <string.h>

/* --- 규칙 함수 3개 --- *
 * 클로저가 없으므로 사실은 ctx에서 읽는다. 각 함수는 별도 심볼이어야 한다
 * (규칙 정체성은 (fn, user_data) 쌍). 판단은 *out_activated로, 증거는 *out_evidence로 낸다. */

/* Rule: 로그인한 사용자인가 (ctx["user"] != NULL) */
static rc_status is_authenticated(rc_context *ctx, rc_spec *const *specs,
                                  size_t n_specs, void *ud,
                                  bool *out_activated, void **out_evidence) {
    (void)specs; (void)n_specs; (void)ud;
    *out_activated = ctx->get(ctx, "user") != NULL;
    *out_evidence  = NULL;
    return RC_OK;
}

/* Counter: IP가 차단 목록에 있는가 (외부 6.6.6.6 또는 내부 10.6.6.6) */
static rc_status is_ip_blocked(rc_context *ctx, rc_spec *const *specs,
                               size_t n_specs, void *ud,
                               bool *out_activated, void **out_evidence) {
    (void)specs; (void)n_specs; (void)ud;
    const char *ip = ctx->get(ctx, "ip");
    *out_activated = ip != NULL &&
                     (strcmp(ip, "6.6.6.6") == 0 || strcmp(ip, "10.6.6.6") == 0);
    *out_evidence  = NULL;
    return RC_OK;
}

/* Except: 내부 IP인가 (ctx["ip"]가 "10."으로 시작) — 차단을 무력화한다 */
static rc_status is_internal_ip(rc_context *ctx, rc_spec *const *specs,
                                size_t n_specs, void *ud,
                                bool *out_activated, void **out_evidence) {
    (void)specs; (void)n_specs; (void)ud;
    const char *ip = ctx->get(ctx, "ip");
    *out_activated = ip != NULL && strncmp(ip, "10.", 3) == 0;
    *out_evidence  = NULL;
    return RC_OK;
}

/* Graded rule(verdict 규칙): 위험 점수 ctx["risk"](0~100 정수)를 v ∈ [-1, +1]로 낸다.
 * 점수 0이면 -1(해당 없음). 주의: 해당 없음은 0이 아니라 -1이다. 0은 '반쯤 해당'이다. */
static rc_status is_risky(rc_context *ctx, rc_spec *const *specs,
                          size_t n_specs, void *ud,
                          double *out_verdict, void **out_evidence) {
    (void)specs; (void)n_specs; (void)ud;
    const int *risk = ctx->get(ctx, "risk");
    *out_verdict  = risk ? (double)*risk / 50.0 - 1.0 : -1.0;   /* 0 → -1, 50 → 0, 100 → +1 */
    *out_evidence = NULL;
    return RC_OK;
}

/* 하나의 요청을 평가해 verdict를 출력한다. */
static void decide(rc_graph *g, const char *label,
                   const char *user, const char *ip) {
    rc_context *ctx = rc_map_context_create();
    if (user) ctx->set(ctx, "user", (void *)user);
    if (ip)   ctx->set(ctx, "ip",   (void *)ip);

    rc_results *res = NULL;
    rc_status st = rc_graph_evaluate(g, ctx, NULL, &res);
    if (st != RC_OK) {
        printf("%-22s ERROR: %s\n", label, rc_graph_last_error(g));
        rc_map_context_destroy(ctx);
        return;
    }

    size_t n;
    const rc_eval_result *items = rc_results_items(res, &n);
    if (n == 0) {
        printf("%-22s no result (rule did not apply) -> deny\n", label);
    } else {
        double v = items[0].verdict;
        printf("%-22s verdict=%+.3f -> %s\n", label, v, v > 0 ? "ALLOW" : "deny");
    }

    rc_results_destroy(res);          /* results는 graph보다 먼저 해제 */
    rc_map_context_destroy(ctx);
}

int main(void) {
    printf("rulecat %s — access-control quick start\n\n", RC_VERSION);

    /* 그래프 구성: 규칙은 노드, 예외는 간선. */
    rc_graph *g = rc_graph_create("api:access");
    rc_rule *auth    = rc_graph_rule(g, is_authenticated, NULL);
    rc_rule *blocked = rc_graph_counter(g, is_ip_blocked, NULL);
    rc_rule *exempt  = rc_graph_except(g, is_internal_ip, NULL);
    rc_rule_attacks(blocked, auth);    /* 차단 IP가 인증을 공격 */
    rc_rule_attacks(exempt, blocked);  /* 내부 IP가 차단을 무력화 (예외의 예외) */

    decide(g, "authed / normal ip",   "alice", "1.2.3.4");   /* warrant only -> +1 */
    decide(g, "authed / blocked ip",  "alice", "6.6.6.6");   /* counter fires  ->  0 */
    decide(g, "authed / internal ip", "alice", "10.6.6.6");  /* except gates the counter -> +1 */
    decide(g, "anonymous",            NULL,    "1.2.3.4");   /* rule inactive  -> none */
    rc_graph_destroy(g);

    /* 등급 규칙 + 공격 강도: 위험 점수는 인증을 반 강도(F = 0.5)로 반박한다.
     *   S = 0.5 · (risk_verdict + 1) / 2,  verdict = 2 / (1 + S) − 1 */
    printf("\ngraded counter (force 0.5)\n");
    rc_graph *g2 = rc_graph_create("api:risk");
    rc_rule *auth2 = rc_graph_rule(g2, is_authenticated, NULL);
    rc_rule *risky = rc_graph_counter_verdict(g2, is_risky, NULL);
    rc_rule_attacks_with_force(risky, auth2, 0.5);

    static const int risks[] = { 0, 50, 100 };
    for (size_t i = 0; i < sizeof risks / sizeof risks[0]; i++) {
        rc_context *ctx = rc_map_context_create();
        ctx->set(ctx, "user", (void *)"alice");
        ctx->set(ctx, "risk", (void *)&risks[i]);
        rc_results *res = NULL;
        if (rc_graph_evaluate(g2, ctx, NULL, &res) == RC_OK) {
            size_t n;
            const rc_eval_result *items = rc_results_items(res, &n);
            if (n > 0)
                printf("risk=%-3d               verdict=%+.3f applicability=%.2f\n",
                       risks[i], items[0].verdict, items[0].applicability);
            rc_results_destroy(res);
        }
        rc_map_context_destroy(ctx);
    }
    rc_graph_destroy(g2);
    return 0;
}
