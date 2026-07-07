#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 — clock_gettime 선언 노출 */

#include "internal.h"

#include <time.h>

/* CLOCK_MONOTONIC ms 단위 double. C11 timespec_get은 단조성 미보장이라 미채택. */
double rc__now_ms(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (double)ts.tv_sec * 1000.0 + (double)ts.tv_nsec / 1e6;
}
