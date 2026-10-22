//ff:func feature=cli type=command control=sequence
//ff:what runEncode — detect locale, encode a .md surface via its codec, verify tiling, and save the .geul
package tanglcli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/spf13/cobra"
)

// runEncode reads the .md surface at args[0], detects its authoring locale from
// the first section header, looks up that locale's codec, and encodes the source
// into a Doc. The result's tiling is verified (the encoder-output gate, README §4)
// and the authoring locale is stamped on the Doc before it is saved as .geul.
// With no codec registered for the detected locale (EN/KO land in tangl-en /
// tangl-ko) it fails with "no codec registered for locale". The output path is
// args[0] with a .geul extension unless -o/--out overrides it.
func runEncode(cmd *cobra.Command, args []string) error {
	in := args[0]
	if filepath.Ext(in) != ".md" {
		return fmt.Errorf("encode expects a .md file, got %q", in)
	}
	src, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	locale, err := detectLocale(src)
	if err != nil {
		return err
	}
	codec, ok := doc.LookupCodec(locale)
	if !ok {
		return fmt.Errorf("no codec registered for locale %q", locale)
	}
	d, err := codec.Encode(src, in)
	if err != nil {
		return err
	}
	d.Locale = locale
	if err := doc.VerifyTiling(d); err != nil {
		return err
	}
	out, err := cmd.Flags().GetString("out")
	if err != nil {
		return err
	}
	if out == "" {
		out = strings.TrimSuffix(in, filepath.Ext(in)) + ".geul"
	}
	return tangeul.Save(d, out)
}
