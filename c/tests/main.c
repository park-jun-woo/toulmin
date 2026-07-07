#include "test.h"

#include <stdio.h>

int rc_test_checks = 0;
int rc_test_failures = 0;

/* 각 테스트 파일이 노출하는 진입점. Phase가 늘면 여기에 등록한다. */
void test_map_context(void);
void test_graph(void);
void test_evaluate(void);
void test_trace(void);
void test_expect_and_cases(void);
void test_run(void);
void test_run_compose(void);

int main(void) {
    test_map_context();
    test_graph();
    test_evaluate();
    test_trace();
    test_expect_and_cases();
    test_run();
    test_run_compose();

    printf("%d checks, %d failures\n", rc_test_checks, rc_test_failures);
    return rc_test_failures == 0 ? 0 : 1;
}
