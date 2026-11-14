//ff:func feature=cli type=command control=sequence
//ff:what TestRunEvaluate — runEvaluate builds the demo graph, evaluates it, and prints JSON verdicts
package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

// TestRunEvaluate covers runEvaluate's success path: NewContext never returns
// nil and the demo graph is acyclic, so Evaluate always succeeds and the
// result is JSON-encoded to stdout.
func TestRunEvaluate(t *testing.T) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = origStdout }()

	cmd := NewEvaluateCmd()
	runErr := cmd.RunE(cmd, nil)

	if closeErr := w.Close(); closeErr != nil {
		t.Fatalf("w.Close() error: %v", closeErr)
	}
	os.Stdout = origStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error: %v", err)
	}

	if runErr != nil {
		t.Fatalf("runEvaluate returned error: %v", runErr)
	}

	var results []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if len(results) == 0 {
		t.Error("expected at least one evaluation result")
	}
}
