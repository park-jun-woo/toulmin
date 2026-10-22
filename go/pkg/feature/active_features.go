//ff:func feature=feature type=adapter control=sequence
//ff:what ActiveFeatures — retrieves enabled feature names from request context
package feature

import "net/http"

// ActiveFeatures retrieves the list of enabled feature names from the request context.
func ActiveFeatures(r *http.Request) []string {
	v, _ := r.Context().Value(featuresKey{}).([]string)
	return v
}
