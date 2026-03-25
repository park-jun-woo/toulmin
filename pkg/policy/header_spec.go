//ff:type feature=policy type=model
//ff:what HeaderSpec — spec for HasHeader rule specifying which header to check
package policy

// HeaderSpec specifies the header name to check.
type HeaderSpec struct {
	Header string `yaml:"header"`
}
