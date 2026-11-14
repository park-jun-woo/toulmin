//ff:func feature=feature type=engine control=sequence
//ff:what TestBuildFeatureContext — tests building a toulmin.Context from a UserContext
package feature

import "testing"

func TestBuildFeatureContext(t *testing.T) {
	uc := &UserContext{
		User:       "someUser",
		ID:         "u-1",
		Region:     "KR",
		Attributes: map[string]any{"beta": true},
	}

	ctx := buildFeatureContext("myFeature", uc)

	if v, ok := ctx.Get("featureName"); !ok || v != "myFeature" {
		t.Fatalf("featureName = %v, %v; want %q, true", v, ok, "myFeature")
	}
	if v, ok := ctx.Get("user"); !ok || v != uc.User {
		t.Fatalf("user = %v, %v; want %v, true", v, ok, uc.User)
	}
	if v, ok := ctx.Get("id"); !ok || v != uc.ID {
		t.Fatalf("id = %v, %v; want %v, true", v, ok, uc.ID)
	}
	if v, ok := ctx.Get("region"); !ok || v != uc.Region {
		t.Fatalf("region = %v, %v; want %v, true", v, ok, uc.Region)
	}
	if v, ok := ctx.Get("attributes"); !ok {
		t.Fatalf("attributes = %v, %v; want present", v, ok)
	}
}
