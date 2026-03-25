//ff:type feature=policy type=model
//ff:what HeaderBacking — backing for HasHeader rule specifying which header to check
package policy

// HeaderBacking specifies the header name to check.
type HeaderBacking struct {
	Header string `yaml:"header"`
}
