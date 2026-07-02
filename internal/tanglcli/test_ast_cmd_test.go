//ff:func feature=cli type=command control=sequence
//ff:what TestAstCmdTransferHasSubjectField — ast on the transfer fixture prints JSON with a subject field
package tanglcli

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TestAstCmdTransferHasSubjectField runs "tangl ast" against transfer.md
// and checks that the printed JSON decodes with subject == "transfer".
func TestAstCmdTransferHasSubjectField(t *testing.T) {
	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"ast", "../../pkg/tangl/parser/testdata/transfer.md"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var decoded struct {
		Subject string `json:"subject"`
	}
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v (output: %s)", err, out.String())
	}
	if decoded.Subject != "transfer" {
		t.Errorf("subject = %q, want %q", decoded.Subject, "transfer")
	}
}
