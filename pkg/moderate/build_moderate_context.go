//ff:func feature=moderate type=adapter control=sequence
//ff:what buildModerateContext: Content + ContentContext → toulmin.Context 변환
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// buildModerateContext converts a Content and ContentContext into a toulmin.Context.
func buildModerateContext(content *Content, cc *ContentContext) toulmin.Context {
	ctx := toulmin.NewContext()
	ctx.Set("body", content.Body)
	ctx.Set("mediaURLs", content.MediaURLs)
	ctx.Set("contentType", content.ContentType)
	ctx.Set("contentMetadata", content.Metadata)
	ctx.Set("author", cc.Author)
	ctx.Set("channel", cc.Channel)
	ctx.Set("metadata", cc.Metadata)
	return ctx
}
