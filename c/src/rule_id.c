#define _POSIX_C_SOURCE 200809L   /* strdup — 전 .c 공통 규약 */

#include "internal.h"

#include <stdlib.h>
#include <string.h>

/* 정체성 문자열 정렬 비교자 (Python sorted / Go sort.Strings 대응). */
static int rc__cmp_id(const void *a, const void *b) {
    const char *sa = *(const char *const *)a;
    const char *sb = *(const char *const *)b;
    return strcmp(sa, sb);
}

/* 각 spec의 정체성 문자열: id()가 있으면 id(), 없으면 name(). */
static const char *rc__spec_identity(const rc_spec *s) {
    const char *v = s->id ? s->id(s) : s->name(s);
    return v ? v : "";
}

char *rc__rule_id(rc_graph *g, rc_rule_func fn, void *user_data,
                  rc_spec *const *specs, size_t n_specs) {
    const char *fid = rc__func_id(g, fn, user_data);
    if (!fid) return NULL;

    if (n_specs == 0) return strdup(fid);

    const char **ids = malloc(n_specs * sizeof *ids);
    if (!ids) return NULL;
    for (size_t i = 0; i < n_specs; i++)
        ids[i] = rc__spec_identity(specs[i]);

    qsort(ids, n_specs, sizeof *ids, rc__cmp_id);

    /* fid + "#" + ids[0] + "+" + ids[1] + ... */
    size_t total = strlen(fid) + 1;             /* fid + '#' */
    for (size_t i = 0; i < n_specs; i++) {
        total += strlen(ids[i]);
        if (i + 1 < n_specs) total += 1;        /* '+' 구분자 */
    }
    total += 1;                                 /* NUL */

    char *out = malloc(total);
    if (!out) { free(ids); return NULL; }

    char *p = out;
    size_t flen = strlen(fid);
    memcpy(p, fid, flen); p += flen;
    *p++ = '#';
    for (size_t i = 0; i < n_specs; i++) {
        size_t len = strlen(ids[i]);
        memcpy(p, ids[i], len); p += len;
        if (i + 1 < n_specs) *p++ = '+';
    }
    *p = '\0';

    free(ids);
    return out;
}
