#define _POSIX_C_SOURCE 200809L   /* 전 .c 공통 규약 */

#include <rulecat/rulecat.h>

#include <string.h>

/* specs 배열에서 name과 일치하는 첫 spec을 선형 탐색으로 찾는다. 없으면 NULL. */
rc_spec *rc_find_spec(rc_spec *const *specs, size_t n_specs, const char *name) {
    for (size_t i = 0; i < n_specs; i++)
        if (strcmp(specs[i]->name(specs[i]), name) == 0)
            return specs[i];
    return NULL;
}
