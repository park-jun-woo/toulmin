//ff:func feature=cli type=command control=sequence
//ff:what TestAstCmdTransferHasSubjectField — ast on the transfer .geul fixture prints JSON with a subject field
package tanglcli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"
)

// TestAstCmdTransferHasSubjectField saves the transfer fixture to a .geul file,
// runs "tangl ast" against it, and checks that the printed JSON decodes with
// subject == "transfer".
func TestAstCmdTransferHasSubjectField(t *testing.T) {
	path := filepath.Join(t.TempDir(), "transfer.geul")
	saveTangeul(t, path, buildTransferDoc())

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"ast", path})
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
