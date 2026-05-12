# pkg/

## Core

| Package | Description | Key API |
|---|---|---|
| `toulmin` | Core engine. Defeats graph + h-Categoriser verdict | `NewGraph`, `Evaluate`, `RunCases` |
| `tangl` | Markdown-based policy language (parser + validator) | `parser.Parse`, `validate.Validate` |

## Domain Frameworks

`toulmin` core 위에 domain-specific Rule 함수 + Spec 타입을 제공하는 프레임워크.

| Package | Domain | Key API | Spec Types |
|---|---|---|---|
| `policy` | 접근 제어 (인증, 인가, IP, Rate limit) | `Guard` (net/http middleware) | `RoleSpec`, `OwnerSpec`, `IPListSpec`, `HeaderSpec` |
| `state` | 상태 전이 (FSM) | `Machine.Can`, `Mermaid()` | `OwnerSpec`, `ExpirySpec` |
| `approve` | 다단계 결재 | `Flow.Evaluate` | `ApproverSpec`, `ThresholdSpec` |
| `price` | 할인 판정 (쿠폰, 멤버십) | `Pricer.Evaluate` | `DiscountSpec`, `MemberSpec`, `BulkOrderSpec` |
| `feature` | 피처 플래그 (롤아웃, 토글) | `Flags.IsEnabled` | `AttributeSpec`, `PercentageSpec`, `RegionSpec` |
| `moderate` | 콘텐츠 모더레이션 (혐오, 스팸) | `Moderator.Review` | `ClassifierSpec`, `TrustScoreSpec`, `MinPostsSpec` |
