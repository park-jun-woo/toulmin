//ff:func feature=tangl type=parser control=sequence
//ff:what Parse — read a TANGL markdown file and parse into File AST
package parser

import "os"

// Parse reads a TANGL markdown file and parses it into a File AST.
func Parse(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseString(string(data))
}
