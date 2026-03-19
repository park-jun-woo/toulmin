package moderate

import "testing"

type mockClassifier struct {
	scores map[string]float64
}

func (m *mockClassifier) Predict(text string, category string) float64 {
	return m.scores[category]
}

func TestIsVerifiedUser(t *testing.T) {
	tests := []struct {
		name     string
		verified bool
		want     bool
	}{
		{"verified", true, true},
		{"not verified", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &ContentContext{Author: &Author{Verified: tt.verified}}
			got, _ := IsVerifiedUser(nil, ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTrustedUser(t *testing.T) {
	tests := []struct {
		name  string
		score float64
		min   float64
		want  bool
	}{
		{"trusted", 0.95, 0.9, true},
		{"not trusted", 0.5, 0.9, false},
		{"equal", 0.9, 0.9, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &ContentContext{Author: &Author{TrustScore: tt.score}}
			got, _ := IsTrustedUser(nil, ctx, tt.min)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsHateSpeech(t *testing.T) {
	tests := []struct {
		name  string
		score float64
		want  bool
	}{
		{"hate", 0.95, true},
		{"not hate", 0.3, false},
		{"borderline", 0.8, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &mockClassifier{scores: map[string]float64{"hate_speech": tt.score}}
			content := &Content{Body: "test"}
			got, _ := ContainsHateSpeech(content, nil, c)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsSpam(t *testing.T) {
	c := &mockClassifier{scores: map[string]float64{"spam": 0.9}}
	got, _ := ContainsSpam(&Content{Body: "test"}, nil, c)
	if !got {
		t.Error("expected spam detected")
	}
}

func TestContainsNSFW(t *testing.T) {
	c := &mockClassifier{scores: map[string]float64{"nsfw": 0.9}}
	got, _ := ContainsNSFW(&Content{Body: "test"}, nil, c)
	if !got {
		t.Error("expected nsfw detected")
	}
}

func TestIsNewsContext(t *testing.T) {
	tests := []struct {
		name    string
		chType  string
		want    bool
	}{
		{"news", "news", true},
		{"general", "general", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &ContentContext{Channel: &Channel{Type: tt.chType}}
			got, _ := IsNewsContext(nil, ctx, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEducational(t *testing.T) {
	ctx := &ContentContext{Channel: &Channel{Type: "education"}}
	got, _ := IsEducational(nil, ctx, nil)
	if !got {
		t.Error("expected true for education channel")
	}
}

func TestIsAdultChannel(t *testing.T) {
	ctx := &ContentContext{Channel: &Channel{AgeGated: true}}
	got, _ := IsAdultChannel(nil, ctx, nil)
	if !got {
		t.Error("expected true for age-gated channel")
	}
}

func TestHasMinPosts(t *testing.T) {
	tests := []struct {
		name  string
		count int
		min   int
		want  bool
	}{
		{"enough", 100, 10, true},
		{"not enough", 5, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &ContentContext{Author: &Author{PostCount: tt.count}}
			got, _ := HasMinPosts(nil, ctx, tt.min)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
