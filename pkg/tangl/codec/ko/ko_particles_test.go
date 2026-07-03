//ff:func feature=tangl type=parser control=sequence
//ff:what TestKoParticles — both allomorphs of each ignore pair are present and non-particles are absent
package ko

import "testing"

func TestKoParticles(t *testing.T) {
	t.Run("AllomorphPairsBothPresent", func(t *testing.T) {
		present := []string{
			"은", "는", "이", "가", "을", "를",
			"와", "과", "로", "으로", "면", "이면",
			"의", "에", "에서", "도",
		}
		for _, p := range present {
			if !koParticles[p] {
				t.Errorf("expected %q to be an ignored particle", p)
			}
		}
	})

	t.Run("NonParticlesAbsent", func(t *testing.T) {
		absent := []string{"까지", "보다", "중", "마다", "이다", "필수다", "무효다", "규칙이다"}
		for _, k := range absent {
			if koParticles[k] {
				t.Errorf("expected %q NOT to be an ignored particle (§주의)", k)
			}
		}
	})
}
