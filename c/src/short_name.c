#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include "internal.h"

#include <string.h>

/* short_name은 full 경로에서 마지막 '.'( '#' 이전 구간 기준) 다음부터 끝까지의
   연속 부분열을 돌려준다 (Python/Go 포트와 동작 일치). 결과는 항상 full의 접미사이므로
   새 할당 없이 full 내부를 가리키는 포인터를 반환한다.
   예: "github.com/example/pkg.IsAdult#Spec{...}" → "IsAdult#Spec{...}" */
const char *rc__short_name(const char *full) {
    const char *hash = strchr(full, '#');
    const char *base_end = hash ? hash : full + strlen(full);
    const char *dot = NULL;
    for (const char *p = full; p < base_end; p++)
        if (*p == '.') dot = p;
    return dot ? dot + 1 : full;
}
