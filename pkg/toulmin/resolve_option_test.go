//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestResolveOption — tests resolveOption for empty opts, Recursive error, Duration-forces-Trace, and pass-through branches
package toulmin

import "testing"

func TestResolveOption(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Empty", func(t *testing.T) {
			opt, err := resolveOption(nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if opt != (EvalOption{}) {
				t.Fatalf("expected zero-value EvalOption, got %+v", opt)
			}
		}},
		{"RecursiveError", func(t *testing.T) {
			_, err := resolveOption([]EvalOption{{Method: Recursive}})
			if err == nil {
				t.Fatal("expected error for Recursive method")
			}
		}},
		{"DurationForcesTrace", func(t *testing.T) {
			opt, err := resolveOption([]EvalOption{{Duration: true}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !opt.Trace {
				t.Fatalf("expected Trace forced true when Duration is true")
			}
			if !opt.Duration {
				t.Fatalf("expected Duration to remain true")
			}
		}},
		{"PassThrough", func(t *testing.T) {
			opt, err := resolveOption([]EvalOption{{Trace: true}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !opt.Trace {
				t.Fatalf("expected Trace true")
			}
			if opt.Duration {
				t.Fatalf("expected Duration false")
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
