//ff:func feature=engine type=model control=sequence
//ff:what validateSpecFailingSpec.Validate — always returns a domain validation error
package toulmin

import "errors"

func (s *validateSpecFailingSpec) Validate() error { return errors.New("domain validation failed") }
