//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what cloneBinderEnv — copies a binder env and adds one new binder
package gen

// cloneBinderEnv copies the codegen binder environment and adds the new
// binder, so a nested exists sees its enclosing binders without mutating the
// caller's scope.
func cloneBinderEnv(env map[string]binderInfo, binder string, info binderInfo) map[string]binderInfo {
	inner := make(map[string]binderInfo, len(env)+1)
	for k, v := range env {
		inner[k] = v
	}
	inner[binder] = info
	return inner
}
