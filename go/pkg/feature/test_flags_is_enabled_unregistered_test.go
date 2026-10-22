//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_IsEnabled_Unregistered — tests error for unregistered feature
package feature

import "testing"

func TestFlags_IsEnabled_Unregistered(t *testing.T) {
	flags := NewFlags()
	_, err := flags.IsEnabled("nonexistent", &UserContext{})
	if err == nil {
		t.Error("expected error for unregistered feature")
	}
}
