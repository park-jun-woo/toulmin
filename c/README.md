# rulecat (C 포트)

Toulmin 논증 모델 기반 규칙 엔진의 C11 포트. **규칙은 함수, 예외는 그래프 간선이다.**
Go 코어(`pkg/toulmin`), TypeScript(`ts/`), Python(`py/src/rulecat/`) 포트와 동작을 일치시킨다.

verdict는 h-Categoriser로 계산하는 연속값 `[-1, +1]`이다:

```
raw(a) = w(a) / (1 + Σ raw(attackers))
verdict = 2 × raw − 1
```

`+1.0` 확정 · `0.0` 미결 · `-1.0` 완전 반박. 이진값이 아니며 프레임워크가 해석한다
(예: 보안 `verdict ≤ 0` → deny).

## 빌드

```bash
make build        # librulecat.a 생성
make test         # 자체 테스트 러너 빌드 후 실행
make example      # examples/quickstart.c 빌드 후 실행
make clean
```

- 표준 C11, 런타임 외부 의존성 0개 (libc + `strdup`/`clock_gettime`용 POSIX만 가정).
- 모든 `.c` 파일 최상단에 `#define _POSIX_C_SOURCE 200809L` 규약.
- 공개 심볼 접두사 `rc_`, 비공개 내부 심볼 `rc__`. 공개 헤더는 `include/rulecat/rulecat.h` 하나.

## 사용 / 링크

```c
#include <rulecat/rulecat.h>
```

정적 라이브러리를 직접 링크:

```bash
cc app.c -Ic/include c/librulecat.a -lm -o app
```

`h-Categoriser`가 `<math.h>`를 쓰므로 `-lm`이 필요하다.

## 배포 형태

C 생태계에는 npm/PyPI 같은 단일 레지스트리가 없으므로 세 가지를 지원한다:

1. **소스 vendoring (권장)** — `c/` 디렉토리를 그대로 복사한다. 의존성이 0이라 별도 설정 없이
   `make build` 한 번이면 빌드된다. `src/*.c`를 자신의 빌드에 직접 넣어도 된다.
2. **`make install`** — 정적 라이브러리 + 헤더 + pkg-config 파일 설치.

   ```bash
   make install PREFIX=/usr/local                     # 시스템 설치 (권한 필요)
   make install PREFIX=/usr/local DESTDIR=/tmp/stage  # 스테이징 디렉토리로 검증

   # 설치 후 pkg-config로 플래그 조회
   cc app.c $(pkg-config --cflags --libs rulecat) -lm -o app
   ```

   설치물: `$(PREFIX)/lib/librulecat.a`, `$(PREFIX)/include/rulecat/rulecat.h`,
   `$(PREFIX)/lib/pkgconfig/rulecat.pc`.
3. (후속) 단일 파일 합본(amalgamation) — sqlite 스타일. 아직 미제공.

## Quick start — access control

킬러 예제(접근 제어)를 C로. 규칙 함수 3개 + 그래프 구성 + evaluate + 결과 출력.
전체 실행본은 [`examples/quickstart.c`](examples/quickstart.c)이며 `make example`로 실행된다.

```c
#include <rulecat/rulecat.h>
#include <stdio.h>
#include <string.h>

/* 규칙 = 함수. 클로저가 없으므로 사실은 ctx에서 읽는다.
   각 함수는 별도 심볼이어야 한다 — 규칙 정체성은 (fn, user_data) 쌍. */
static rc_status is_authenticated(rc_context *ctx, rc_spec *const *s, size_t n,
                                  void *ud, bool *act, void **ev) {
    (void)s; (void)n; (void)ud; *ev = NULL;
    *act = ctx->get(ctx, "user") != NULL;
    return RC_OK;
}
static rc_status is_ip_blocked(rc_context *ctx, rc_spec *const *s, size_t n,
                               void *ud, bool *act, void **ev) {
    (void)s; (void)n; (void)ud; *ev = NULL;
    const char *ip = ctx->get(ctx, "ip");
    *act = ip && (strcmp(ip, "6.6.6.6") == 0 || strcmp(ip, "10.6.6.6") == 0);
    return RC_OK;
}
static rc_status is_internal_ip(rc_context *ctx, rc_spec *const *s, size_t n,
                                void *ud, bool *act, void **ev) {
    (void)s; (void)n; (void)ud; *ev = NULL;
    const char *ip = ctx->get(ctx, "ip");
    *act = ip && strncmp(ip, "10.", 3) == 0;
    return RC_OK;
}

int main(void) {
    rc_graph *g = rc_graph_create("api:access");
    rc_rule *auth    = rc_graph_rule(g, is_authenticated, NULL);
    rc_rule *blocked = rc_graph_counter(g, is_ip_blocked, NULL);
    rc_rule *exempt  = rc_graph_except(g, is_internal_ip, NULL);
    rc_rule_attacks(blocked, auth);    /* 차단 IP가 인증을 공격 */
    rc_rule_attacks(exempt, blocked);  /* 내부 IP가 차단을 무력화 (예외의 예외) */

    rc_context *ctx = rc_map_context_create();
    ctx->set(ctx, "user", (void *)"alice");
    ctx->set(ctx, "ip",   (void *)"10.6.6.6");   /* 차단 + 내부 → except가 구제 */

    rc_results *res = NULL;
    if (rc_graph_evaluate(g, ctx, NULL, &res) == RC_OK) {
        size_t n; const rc_eval_result *items = rc_results_items(res, &n);
        if (n > 0) printf("verdict=%+.3f\n", items[0].verdict);  /* +0.333 → allow */
        rc_results_destroy(res);       /* results는 graph보다 먼저 해제 */
    }
    rc_map_context_destroy(ctx);
    rc_graph_destroy(g);
    return 0;
}
```

요청이 진화해도 매번 2줄 추가일 뿐 구조는 그대로다.
`blocked.attacks(auth)`, `exempt.attacks(blocked)` — 예외의 예외가 간선 하나다.

## C 규약

다른 포트에는 없는, C 포트를 안전하게 쓰기 위한 규약이다.

### 에러 처리 — 예외 대신 상태 코드 + 지연 보고

- 대부분의 함수는 `rc_status`를 반환한다 (`RC_OK`, `RC_ERR_CYCLE`, `RC_ERR_RUN_CYCLE`,
  `RC_ERR_RUN_DEPTH`, `RC_ERR_INVALID_ARG` 등).
- 등록 함수(`rc_graph_rule/counter/except`)는 실패 시 `NULL`을 반환한다.
- 빌더 체인(`rc_rule_with`, `rc_rule_qualifier`, `rc_rule_run_on`, `rc_rule_run`)은
  체이닝을 위해 `self`를 반환하고, **오류를 래치(지연 기록)** 한다 — 실제 보고는
  `rc_graph_evaluate` / `rc_graph_run` 진입 시 `RC_ERR_BUILDER`로 이뤄진다.
- 마지막 오류 메시지는 `rc_graph_last_error(g)`로 조회한다 (graph 소유 버퍼, 다음 오류가 덮어씀).

### 메모리 소유권 (반드시 이 표를 따른다)

| 객체 | 생성 | 소유자 | 해제 |
|---|---|---|---|
| `rc_graph` | `rc_graph_create` | 호출자 | `rc_graph_destroy` — 내부 rule 메타·핸들·ID 문자열·간선 일괄 해제 |
| `rc_rule` 핸들 | `rc_graph_rule` 등 | **graph** | `rc_graph_destroy`가 함께 해제. 호출자는 절대 free 금지 |
| `rc_spec` | 호출자 | **호출자** | graph는 포인터만 차용. **graph보다 오래 살아야 한다** |
| `rc_context` | 호출자 | 호출자 | `rc_map_context_destroy` (또는 사용자 구현의 소멸 함수) |
| evidence (`void *`) | 규칙 함수 | 기본: 규칙 작성자 (엔진은 free 안 함) / dtor 등록 시: 결과 객체 | dtor 등록 시 결과 해제가 evidence마다 정확히 1회 호출 |
| `rc_results` | `rc_graph_evaluate` | 호출자 | `rc_results_destroy` |
| `rc_run_result` | `rc_graph_run` | 호출자 | `rc_run_result_destroy` — results + trace + evidence 풀 일괄 해제 |
| 결과/trace 내 `const char *name` | graph | **graph** | 결과는 문자열을 차용 → **results/trace는 graph보다 먼저 해제** |

> **핵심 규칙: `rc_results` / `rc_run_result`는 그것을 만든 `rc_graph`보다 먼저 해제한다.**
> 결과 객체가 graph 소유 문자열을 차용하기 때문이다.

### 클로저 대체 — `fn + void *user_data`

C에는 클로저가 없다. 모든 콜백은 **함수 포인터 + `void *user_data`** 쌍이다.
규칙의 정체성 키도 `(fn, user_data)` 쌍이다:

- 같은 `fn`을 서로 다른 `user_data`로 등록 → **서로 다른 규칙** (다른 ID).
- 같은 `(fn, user_data)` 쌍을 두 번 등록 → 중복 오류.

같은 함수를 다른 기준으로 재사용하려면 `user_data` 또는 `rc_rule_with(spec)`로 구분한다
(클로저 팩토리 불필요).

### evidence 소유권

규칙 함수의 `out_evidence`는 기본적으로 **차용(`void *`)** 이다 — 엔진은 절대 free하지 않는다.
소유를 결과 객체로 넘기려면 `rc_graph_set_evidence_destructor(g, dtor)`를 등록한다. 그러면
결과 해제 시 평가 패스에서 생산된 **서로 다른 evidence 포인터마다 정확히 1회** dtor가 호출된다
(evidence는 result와 trace 양쪽에 같은 포인터로 나타날 수 있으므로 "등장 횟수"가 아니라
노드당 1개 기준).

### 스레드

단일 스레드 가정. graph 구성·평가·해제를 동시 스레드에서 섞지 않는다 (Go 코어도 동일).

## Run / RunOn — 판단 후 행동

`rc_graph_evaluate`는 판단만 하고 반환한다(순수·멱등). `rc_graph_run`은 먼저 full pass로
판단한 뒤, **Active** 노드마다 핸들러를 발화한다. 노드는 `activated && verdict > 0`일 때만
Active다 — 그것이 유일한 이벤트다.

```c
static rc_status audit(const rc_trace_entry *self, const rc_trace *t, void *ud) {
    (void)t; (void)ud;
    printf("fired: %s (verdict=%+.3f)\n", self->name, self->verdict);
    return RC_OK;   /* RC_OK 외 반환 시 Run이 RC_ERR_HANDLER로 중단 */
}

rc_rule *r = rc_graph_rule(g, is_authenticated, NULL);
rc_rule_run_on(r, audit, NULL);                 /* Active 발화 핸들러 */
rc_rule_run(r, notify_subgraph);                /* Active → 하위 그래프 Run (실행 합성) */

rc_run_result *rr = NULL;
if (rc_graph_run(g, ctx, NULL, &rr) == RC_OK) {
    const rc_trace *tr = rc_run_result_trace(rr);
    size_t n; const rc_trace_entry *all = rc_trace_all(tr, &n);   /* 전 노드, 등록 순서 */
    /* rc_trace_get(tr, name), rc_trace_ctx(tr)로도 조회 */
}
if (rr) rc_run_result_destroy(rr);              /* NULL이 아니면 항상 해제 */
```

- 핸들러는 등록 순서로 발화하며 첫 오류에서 Run이 중단된다.
- 각 `rc_trace_entry`는 **Claim**(`name`) · **Ground**(`ground` = ctx 그대로) ·
  **Backing**(`specs`) · **Verdict**를 노출한다 — 별도 뷰 없이 감사/설명이 가능하다.
- 실행 합성: `rc_rule_run(r, sub)`은 Active 노드에서 하위 그래프 `sub`를 같은 ctx로 Run한다.
  판단은 `attacks`로 위로, 실행은 `run`으로 아래로 합성된다. 하위 그래프 verdict는 격리(즉시
  폐기)되고 **오류만** 전파된다. 실행 합성은 DAG여야 하며 순환은 `RC_ERR_RUN_CYCLE`,
  깊이는 64로 제한(`RC_ERR_RUN_DEPTH`)된다. `sub`는 차용 — 부모 graph destroy가 해제하지 않는다.

## 테스트 헬퍼 — 테이블 주도 케이스

```c
rc_test_case cases[] = {
    { .name = "authed normal",   .ctx = NULL, .opt = NULL, .expect = rc_verdict_above(0.0)       },
    { .name = "blocked ip",      .ctx = NULL, .opt = NULL, .expect = rc_verdict_at_most(0.0)     },
    { .name = "internal rescue", .ctx = NULL, .opt = NULL, .expect = rc_verdict_between(0.0, 0.5)},
    { .name = "anonymous",       .ctx = NULL, .opt = NULL, .expect = rc_no_result()              },
};
char why[256];
rc_status st = rc_run_cases(g, cases, 4, why, sizeof why);   /* 첫 실패에서 중단, why에 "[케이스명] 사유" */
```

- `.ctx = NULL`이면 케이스마다 임시 `rc_map_context`를 만들고 해제한다.
- 단건 대조는 `rc_expect_check(results, expectation, errbuf, errlen)`.

| 기대값 | 조건 |
|---|---|
| `rc_verdict_above(v)` | verdict > v |
| `rc_verdict_at_most(v)` | verdict ≤ v |
| `rc_verdict_between(lo, hi)` | lo < verdict ≤ hi |
| `rc_no_result()` | 결과 0건 |

## 4개 포트 API 대응표

| Go | TypeScript | Python | C |
|---|---|---|---|
| `toulmin.NewGraph("name")` | `new Graph("name")` | `Graph("name")` | `rc_graph_create("name")` |
| `g.Rule(fn)` | `g.rule(fn)` | `g.rule(fn)` | `rc_graph_rule(g, fn, ud)` |
| `g.Counter(fn)` | `g.counter(fn)` | `g.counter(fn)` | `rc_graph_counter(g, fn, ud)` |
| `g.Except(fn)` | `g.except(fn)` | `g.except_(fn)` | `rc_graph_except(g, fn, ud)` |
| `r.Attacks(t)` | `r.attacks(t)` | `r.attacks(t)` | `rc_rule_attacks(r, t)` |
| `r.With(spec)` | `r.with(spec)` | `r.with_spec(spec)` | `rc_rule_with(r, spec)` |
| `r.Qualifier(q)` | `r.qualifier(q)` | `r.qualifier(q)` | `rc_rule_qualifier(r, q)` |
| `r.RunOn(h)` | `r.runOn(h)` | `r.run_on(h)` | `rc_rule_run_on(r, h, ud)` |
| `r.Run(sub)` | `r.run(sub)` | `r.run(sub)` | `rc_rule_run(r, sub)` |
| `g.Evaluate(ctx, opt)` | `g.evaluate(ctx, opt)` | `g.evaluate(ctx, opt)` | `rc_graph_evaluate(g, ctx, opt, &res)` |
| `g.Run(ctx, opt)` | `g.run(ctx, opt)` | `g.run(ctx, opt)` | `rc_graph_run(g, ctx, opt, &rr)` |
| `toulmin.NewContext()` | `newContext()` | `MapContext()` | `rc_map_context_create()` |
| `toulmin.VerdictAbove(v)` | `verdictAbove(v)` | `verdict_above(v)` | `rc_verdict_above(v)` |
| `toulmin.NoResult` | `noResult` | `no_result` | `rc_no_result()` |
| `toulmin.RunCases(t, g, cs)` | `TestCase` 타입만 (러너 없음) | `TestCase` 타입만 (러너 없음) | `rc_run_cases(g, cs, n, buf, len)` |
| (GC) | (GC) | (GC) | `rc_*_destroy` 쌍 (호출자 책임) |

> **포트 차이**: 테이블 주도 러너는 Go의 `toulmin.RunCases`와 C의 `rc_run_cases`에만 있다.
> TS/Py는 `TestCase` 타입과 `verdictAbove`/`verdict_above` 등 기대값 헬퍼만 노출하고, 케이스를
> 도는 루프는 사용자가 직접 작성한다.

## 4개 포트 동작 일치 검증표

| 시나리오 | Go | TS | Python | C |
|---|---|---|---|---|
| 단일 warrant | 1.0 | 1.0 | 1.0 | 1.0 |
| warrant + counter | 0.0 | 0.0 | 0.0 | 0.0 |
| 보상 체인 (W←R←D) | 1/3 | 1/3 | 1/3 | 1/3 |
| qualifier 0.0 | -1.0 | -1.0 | -1.0 | -1.0 |
| 비활성 규칙 | 결과 없음 | 결과 없음 | 결과 없음 | 결과 없음 |
| defeats 순환 | error | throw | raise | `RC_ERR_CYCLE` |
| run 순환 / 깊이 초과 | error | throw | raise | `RC_ERR_RUN_CYCLE` / `RC_ERR_RUN_DEPTH` |
| Run: Active만 발화, 핸들러→하위그래프, verdict 격리 | 동일 | 동일 | 동일 | 동일 |

## 공개 API 표면

전 심볼은 `include/rulecat/rulecat.h` 하나에 선언된다.

| 분류 | 심볼 |
|---|---|
| 상태/에러 | `rc_status`, `rc_graph_last_error` |
| Context | `rc_context`(vtable), `rc_map_context_create/destroy` |
| Spec | `rc_spec`(vtable), `rc_find_spec` |
| 타입 | `rc_strength`, `rc_eval_method`, `rc_eval_option`, `rc_trace_entry`, `rc_eval_result` |
| Graph | `rc_graph_create/destroy`, `rc_graph_rule/counter/except`, `rc_graph_set_evidence_destructor` |
| Rule 빌더 | `rc_rule_attacks`, `rc_rule_with`, `rc_rule_qualifier`, `rc_rule_run_on`, `rc_rule_run`, `rc_rule_id` |
| 평가 | `rc_graph_evaluate`, `rc_results_items`, `rc_results_destroy` |
| Run | `rc_graph_run`, `rc_run_result_results/trace/destroy`, `rc_trace_all/get/ctx`, `rc_node_handler` |
| 테스트 | `rc_expectation`, `rc_verdict_above/at_most/between`, `rc_no_result`, `rc_expect_check`, `rc_test_case`, `rc_run_cases` |

## License

MIT
