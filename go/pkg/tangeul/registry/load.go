//ff:func feature=tangl type=util control=sequence
//ff:what Load — recursively loads every document reachable via checking-alias See references
package registry

// Load loads rootPath and recursively resolves every checking-alias
// reference reachable from it: for each node with n.Checking != nil &&
// n.Checking.Alias != "", it resolves the alias's See.Package via
// resolveCheckingAliasPath, loads that .geul, and recurses into it.
// Returns an error if a referenced .geul cannot be loaded, if an alias
// used in a checking clause has no matching See declaration, or if a
// document-level load cycle is detected. Diamond references (A sees B and C,
// both seeing D) parse D only once — already-loaded documents are reused
// from the Registry cache. Load wires the production tangeul-backed loader;
// the recursion itself lives in loadRegistry so tests can inject a fake.
func Load(rootPath string) (*Registry, error) {
	return loadRegistry(rootPath, loadDocAST)
}
