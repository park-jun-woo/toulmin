//ff:func feature=feature type=engine control=iteration dimension=1
//ff:what TestActiveFeatures — tests retrieving active feature names from request context
package feature

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestActiveFeatures(t *testing.T) {
	cases := []struct {
		name string
		ctx  context.Context
		want []string
	}{
		{
			name: "present",
			ctx:  context.WithValue(context.Background(), featuresKey{}, []string{"a", "b"}),
			want: []string{"a", "b"},
		},
		{
			name: "absent",
			ctx:  context.Background(),
			want: nil,
		},
		{
			name: "wrong type",
			ctx:  context.WithValue(context.Background(), featuresKey{}, "not-a-slice"),
			want: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := ActiveFeatures(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("ActiveFeatures() = %v, want %v", got, c.want)
			}
		})
	}
}
