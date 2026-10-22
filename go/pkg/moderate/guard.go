//ff:func feature=moderate type=adapter control=sequence
//ff:what Guard: 콘텐츠 제출 엔드포인트에 모더레이션 판정 적용
package moderate

import "net/http"

// Guard returns a net/http middleware that moderates content.
// Block → 403, Flag → 202, Allow → next.
func Guard(m *Moderator, extractFn ExtractFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			content, ctx := extractFn(r)
			result, err := m.Review(content, ctx)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"moderation failed"}`))
				return
			}
			switch result.Action {
			case ActionBlock:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"content blocked"}`))
			case ActionFlag:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte(`{"status":"flagged for review"}`))
			default:
				next.ServeHTTP(w, r)
			}
		})
	}
}
