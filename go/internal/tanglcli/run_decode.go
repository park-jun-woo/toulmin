//ff:func feature=cli type=command control=sequence
//ff:what runDecode — render a .geul back to text: verbatim mode 1, or codec-driven mode 2 with --locale
package tanglcli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/spf13/cobra"
)

// runDecode loads the .geul at args[0] and renders it. With no --locale it is
// mode 1: doc.RenderVerbatim dumps the authored Source (byte identity, README §3).
// With --locale it is mode 2: the target locale's codec regenerates statements via
// doc.Render; a locale with no registered codec (EN/KO land in tangl-en/tangl-ko)
// fails with "no codec registered for locale". Output goes to -o/--out or stdout.
func runDecode(cmd *cobra.Command, args []string) error {
	d, err := tangeul.Load(args[0])
	if err != nil {
		return err
	}
	locale, err := cmd.Flags().GetString("locale")
	if err != nil {
		return err
	}
	var rendered []byte
	if locale == "" {
		rendered = doc.RenderVerbatim(d)
	} else {
		codec, ok := doc.LookupCodec(locale)
		if !ok {
			return fmt.Errorf("no codec registered for locale %q", locale)
		}
		rendered, err = doc.Render(d, codec)
		if err != nil {
			return err
		}
	}
	out, err := cmd.Flags().GetString("out")
	if err != nil {
		return err
	}
	if out == "" {
		_, err = cmd.OutOrStdout().Write(rendered)
		return err
	}
	return os.WriteFile(out, rendered, 0644)
}
