# pkg/

## Core

| Package | Description | Key API |
|---|---|---|
| `toulmin` | 엔진 코어. defeats graph + h-Categoriser verdict. YAML 파서, 검증, 코드 생성 포함 | `NewGraph`, `Evaluate`, `ParseYAML`, `ValidateGraphDef`, `LoadGraph`, `GenerateGraph` |
| `analyzer` | Go AST 정적 분석. 소스 코드에서 defeats graph 추출 | `ExtractDefeats` |

## Domain Frameworks

`toulmin` 코어 위에 도메인별 Rule 함수 + Backing 타입을 제공하는 프레임워크.

| Package | Domain | Key API | Backing Types |
|---|---|---|---|
| `policy` | 접근 제어 (인증, 인가, IP, Rate limit) | `Guard` (net/http 미들웨어) | `RoleBacking`, `OwnerBacking`, `IPListBacking` |
| `state` | 상태 전이 (FSM) | `Machine.Can`, `Mermaid()` | — |
| `approve` | 다단계 결재 | `Flow.Evaluate` | — |
| `price` | 할인 판정 (쿠폰, 멤버십) | `Pricer.Evaluate` | `DiscountBacking` |
| `feature` | 피처 플래그 (롤아웃, 토글) | `Flags.IsEnabled` | — |
| `moderate` | 콘텐츠 모더레이션 (혐오, 스팸) | `Moderator.Review` | — |
