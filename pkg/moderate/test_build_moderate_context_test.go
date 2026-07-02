//ff:func feature=moderate type=engine control=sequence
//ff:what TestBuildModerateContext — tests building a toulmin.Context from Content and ContentContext
package moderate

import "testing"

func TestBuildModerateContext(t *testing.T) {
	content := &Content{
		Body:        "hello world",
		MediaURLs:   []string{"http://example.com/a.png"},
		ContentType: "text",
		Metadata:    map[string]any{"lang": "en"},
	}
	cc := &ContentContext{
		Author:   &Author{ID: "u1", Verified: true, PostCount: 5, TrustScore: 0.8},
		Channel:  &Channel{ID: "c1", Type: "public", AgeGated: false},
		Metadata: map[string]any{"source": "web"},
	}

	ctx := buildModerateContext(content, cc)

	if v, ok := ctx.Get("body"); !ok || v != content.Body {
		t.Fatalf("body = %v, %v; want %q, true", v, ok, content.Body)
	}
	if v, ok := ctx.Get("mediaURLs"); !ok {
		t.Fatalf("mediaURLs = %v, %v; want present", v, ok)
	}
	if v, ok := ctx.Get("contentType"); !ok || v != content.ContentType {
		t.Fatalf("contentType = %v, %v; want %q, true", v, ok, content.ContentType)
	}
	if v, ok := ctx.Get("contentMetadata"); !ok {
		t.Fatalf("contentMetadata = %v, %v; want present", v, ok)
	}
	if v, ok := ctx.Get("author"); !ok || v != cc.Author {
		t.Fatalf("author = %v, %v; want %v, true", v, ok, cc.Author)
	}
	if v, ok := ctx.Get("channel"); !ok || v != cc.Channel {
		t.Fatalf("channel = %v, %v; want %v, true", v, ok, cc.Channel)
	}
	if v, ok := ctx.Get("metadata"); !ok {
		t.Fatalf("metadata = %v, %v; want present", v, ok)
	}
}
