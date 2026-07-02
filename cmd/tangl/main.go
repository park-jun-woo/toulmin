//ff:func feature=cli type=command control=sequence
//ff:what main — tangl CLI entrypoint
package main

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/toulmin/internal/tanglcli"
)

func main() {
	if err := tanglcli.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
