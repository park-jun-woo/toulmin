//ff:func feature=moderate type=adapter control=sequence
//ff:what Guard: 콘텐츠 제출 엔드포인트에 모더레이션 판정 적용
package moderate

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Guard returns a gin.HandlerFunc that moderates content.
// Block → 403, Flag → 202, Allow → c.Next().
func Guard(m *Moderator, extractor ContentExtractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		content, ctx := extractor(c)
		result, err := m.Review(content, ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "moderation failed"})
			return
		}
		switch result.Action {
		case ActionBlock:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "content blocked"})
		case ActionFlag:
			c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"status": "flagged for review"})
		default:
			c.Next()
		}
	}
}
