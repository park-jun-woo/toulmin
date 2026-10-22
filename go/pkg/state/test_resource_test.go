//ff:type feature=state type=model
//ff:what testResource — test helper resource type with owner and expiry
package state

import "time"

type testResource struct {
	OwnerID   string
	ExpiresAt time.Time
}
