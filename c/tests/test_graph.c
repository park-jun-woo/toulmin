#include "test.h"

#include <rulecat/rulecat.h>

/* 최소 규칙 함수 스텁 — 항상 활성. */
static rc_status always_true(rc_context *ctx, rc_spec *const *specs, size_t n,
                             void *ud, bool *out_activated, void **out_evidence) {
    (void)ctx; (void)specs; (void)n; (void)ud;
    *out_activated = true;
    *out_evidence = NULL;
    return RC_OK;
}

static rc_status other_fn(rc_context *ctx, rc_spec *const *specs, size_t n,
                          void *ud, bool *out_activated, void **out_evidence) {
    (void)ctx; (void)specs; (void)n; (void)ud;
    *out_activated = false;
    *out_evidence = NULL;
    return RC_OK;
}

/* --- with 테스트용 최소 rc_spec 스텁 --- */
static const char *stub_spec_name(const rc_spec *self)     { (void)self; return "StubSpec"; }
static rc_status   stub_spec_validate(const rc_spec *self) { (void)self; return RC_OK; }
static const char *stub_spec_id(const rc_spec *self)       { (void)self; return "StubSpec{v=1}"; }

/* validate가 실패하는 스펙(with 지연 오류 검증용) */
static rc_status   bad_spec_validate(const rc_spec *self)  { (void)self; return RC_ERR_INVALID_ARG; }

/* 등록 + 중복 검출: 같은 (fn,ud)는 중복, ud가 다르면 별개. */
static void test_register_and_duplicate(void) {
    rc_graph *g = rc_graph_create("test");
    ASSERT_NOT_NULL(g);

    rc_rule *r1 = rc_graph_rule(g, always_true, NULL);
    ASSERT_NOT_NULL(r1);

    /* 같은 (fn, user_data) 쌍 재등록 → 중복 → NULL + "duplicate" */
    rc_rule *dup = rc_graph_rule(g, always_true, NULL);
    ASSERT_NULL(dup);
    ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "duplicate"));

    /* user_data가 다르면 별개 규칙 */
    static int tag;
    rc_rule *r2 = rc_graph_counter(g, always_true, &tag);
    ASSERT_NOT_NULL(r2);

    /* 다른 함수도 별개 규칙 */
    rc_rule *r3 = rc_graph_except(g, other_fn, NULL);
    ASSERT_NOT_NULL(r3);

    /* fn == NULL → 등록 실패 */
    ASSERT_NULL(rc_graph_rule(g, NULL, NULL));

    rc_graph_destroy(g);
}

/* qualifier 범위 위반은 지연 처리 — self를 그대로 반환(체인 유지), 크래시 없음. */
static void test_qualifier_range_deferred(void) {
    rc_graph *g = rc_graph_create("test");
    rc_rule *r = rc_graph_rule(g, always_true, NULL);
    ASSERT_NOT_NULL(r);

    ASSERT_PTR_EQ(rc_rule_qualifier(r, 1.5), r);   /* 범위 초과 — 체인 유지 */
    ASSERT_PTR_EQ(rc_rule_qualifier(r, -0.1), r);  /* 범위 미만 — 체인 유지 */
    ASSERT_PTR_EQ(rc_rule_qualifier(r, 0.8), r);   /* 정상 — 체인 유지 */

    /* NULL 안전 */
    ASSERT_NULL(rc_rule_qualifier(NULL, 0.5));

    rc_graph_destroy(g);
}

/* with는 규칙 ID를 바꾼다(spec 정체성 반영). */
static void test_with_changes_id(void) {
    rc_graph *g = rc_graph_create("test");
    rc_rule *r = rc_graph_rule(g, always_true, NULL);
    ASSERT_NOT_NULL(r);

    const char *before = rc_rule_id(r);
    ASSERT_NOT_NULL(before);
    /* strdup 없이 비교 위해 스냅샷 저장(before는 with 후 free됨) */
    char snapshot[256];
    strncpy(snapshot, before, sizeof snapshot - 1);
    snapshot[sizeof snapshot - 1] = '\0';

    rc_spec spec = { stub_spec_name, stub_spec_validate, stub_spec_id, NULL };
    rc_rule *chained = rc_rule_with(r, &spec);
    ASSERT_PTR_EQ(chained, r);   /* 체이닝 — self 반환 */

    const char *after = rc_rule_id(r);
    ASSERT_NOT_NULL(after);
    ASSERT_FALSE(strcmp(snapshot, after) == 0);       /* ID가 바뀜 */
    ASSERT_NOT_NULL(strstr(after, "StubSpec{v=1}"));  /* spec 정체성 포함 */

    rc_graph_destroy(g);
}

/* with 오류(NULL spec / validate 실패)는 지연 처리 — self 반환, builder 오류 래치. */
static void test_with_deferred_errors(void) {
    rc_graph *g = rc_graph_create("test");
    rc_rule *r = rc_graph_rule(g, always_true, NULL);
    ASSERT_NOT_NULL(r);

    ASSERT_PTR_EQ(rc_rule_with(r, NULL), r);   /* NULL spec — self 반환 */

    rc_spec bad = { stub_spec_name, bad_spec_validate, NULL, NULL };
    ASSERT_PTR_EQ(rc_rule_with(r, &bad), r);   /* validate 실패 — self 반환 */
    ASSERT_NOT_NULL(strstr(rc_graph_last_error(g), "invalid"));

    /* NULL 안전 */
    ASSERT_NULL(rc_rule_with(NULL, &bad));

    rc_graph_destroy(g);
}

/* attacks: 유효 간선은 RC_OK, NULL/교차 그래프는 RC_ERR_INVALID_ARG. */
static void test_attacks(void) {
    rc_graph *g = rc_graph_create("test");
    rc_rule *a = rc_graph_rule(g, always_true, NULL);
    rc_rule *b = rc_graph_counter(g, other_fn, NULL);
    ASSERT_NOT_NULL(a);
    ASSERT_NOT_NULL(b);

    ASSERT_TRUE(rc_rule_attacks(b, a) == RC_OK);

    /* NULL 인자 */
    ASSERT_TRUE(rc_rule_attacks(NULL, a) == RC_ERR_INVALID_ARG);
    ASSERT_TRUE(rc_rule_attacks(a, NULL) == RC_ERR_INVALID_ARG);

    /* 교차 그래프 */
    rc_graph *g2 = rc_graph_create("other");
    rc_rule *c = rc_graph_rule(g2, always_true, NULL);
    ASSERT_NOT_NULL(c);
    ASSERT_TRUE(rc_rule_attacks(a, c) == RC_ERR_INVALID_ARG);

    rc_graph_destroy(g);
    rc_graph_destroy(g2);
}

void test_graph(void) {
    test_register_and_duplicate();
    test_qualifier_range_deferred();
    test_with_changes_id();
    test_with_deferred_errors();
    test_attacks();
}
