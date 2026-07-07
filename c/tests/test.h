#ifndef RC_TEST_H
#define RC_TEST_H

/* 외부 프레임워크 없는 assert 기반 자체 테스트 매크로.
   러너의 등록 방식은 Phase 003에서 완성한다. 지금은 카운터 기반 뼈대. */

#include <math.h>
#include <stdio.h>
#include <string.h>

/* main.c에서 정의되는 전역 카운터 */
extern int rc_test_checks;
extern int rc_test_failures;

#define RC_CHECK(cond, msg) do {                                        \
        rc_test_checks++;                                               \
        if (!(cond)) {                                                  \
            rc_test_failures++;                                         \
            fprintf(stderr, "FAIL %s:%d: %s\n", __FILE__, __LINE__, msg); \
        }                                                              \
    } while (0)

#define ASSERT_TRUE(cond)      RC_CHECK((cond), #cond)
#define ASSERT_FALSE(cond)     RC_CHECK(!(cond), "!(" #cond ")")
#define ASSERT_NULL(p)         RC_CHECK((p) == NULL, #p " == NULL")
#define ASSERT_NOT_NULL(p)     RC_CHECK((p) != NULL, #p " != NULL")
#define ASSERT_PTR_EQ(a, b)    RC_CHECK((a) == (b), #a " == " #b)
#define ASSERT_STR_EQ(a, b)    RC_CHECK(strcmp((a), (b)) == 0, #a " == " #b)
#define ASSERT_NEAR(a, b, eps) RC_CHECK(fabs((a) - (b)) < (eps), #a " ~= " #b)

#endif /* RC_TEST_H */
