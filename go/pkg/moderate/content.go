//ff:type feature=moderate type=model
//ff:what Content: 모더레이션 대상 콘텐츠
package moderate

// Content represents the content to be moderated.
type Content struct {
	Body        string
	MediaURLs   []string
	ContentType string
	Metadata    map[string]any
}
