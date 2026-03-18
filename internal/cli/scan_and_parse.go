//ff:func feature=scanner type=scanner control=iteration dimension=1
//ff:what scanAndParse — scans Go files and parses all rule declarations
package cli

import (
	"github.com/park-jun-woo/toulmin/internal/scanner"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// scanAndParse extracts rule declarations from all Go files and returns
// the package name and parsed RuleMetas.
func scanAndParse(paths []string) (string, []toulmin.RuleMeta, error) {
	var pkgName string
	var metas []toulmin.RuleMeta
	for _, path := range paths {
		pkg, decls, err := scanner.ExtractRules(path)
		if err != nil {
			return "", nil, err
		}
		if len(decls) == 0 {
			continue
		}
		pkgName = pkg
		metas = append(metas, parseDecls(decls)...)
	}
	return pkgName, metas, nil
}
