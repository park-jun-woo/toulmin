//ff:func feature=feature type=engine control=iteration dimension=1
//ff:what List: 사용자에 대해 활성화된 전체 피처 목록
package feature

// List returns the names of all enabled features for the given user context.
func (f *Flags) List(ctx *UserContext) ([]string, error) {
	var enabled []string
	for _, name := range f.order {
		ok, err := f.IsEnabled(name, ctx)
		if err != nil {
			return nil, err
		}
		if ok {
			enabled = append(enabled, name)
		}
	}
	return enabled, nil
}
