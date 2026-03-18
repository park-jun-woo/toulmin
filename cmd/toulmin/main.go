//ff:func feature=cli type=command control=sequence
//ff:what main — toulmin CLI entrypoint
package main

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/toulmin/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
