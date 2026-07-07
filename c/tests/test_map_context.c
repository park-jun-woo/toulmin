#include "test.h"

#include <rulecat/rulecat.h>

/* 기본 Context 구현(rc_map_context)의 get/set/미존재 키 동작 검증. */
void test_map_context(void) {
    rc_context *ctx = rc_map_context_create();
    ASSERT_NOT_NULL(ctx);
    if (!ctx) return;

    /* 미존재 키는 NULL */
    ASSERT_NULL(ctx->get(ctx, "missing"));

    /* set 후 get은 같은 포인터(값 차용) */
    int a = 42;
    ctx->set(ctx, "age", &a);
    ASSERT_PTR_EQ(ctx->get(ctx, "age"), &a);

    /* 다른 키는 여전히 NULL */
    ASSERT_NULL(ctx->get(ctx, "name"));

    /* 같은 키 재설정은 값 갱신(중복 노드 생성 안 함) */
    int b = 7;
    ctx->set(ctx, "age", &b);
    ASSERT_PTR_EQ(ctx->get(ctx, "age"), &b);

    /* 여러 키 공존 */
    const char *nm = "toulmin";
    ctx->set(ctx, "name", (void *)nm);
    ASSERT_PTR_EQ(ctx->get(ctx, "name"), (void *)nm);
    ASSERT_PTR_EQ(ctx->get(ctx, "age"), &b);

    rc_map_context_destroy(ctx);

    /* NULL 안전 */
    rc_map_context_destroy(NULL);
}
