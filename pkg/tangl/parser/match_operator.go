//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what matchOperator — match operator tokens from word list and return operator string
package parser

import "fmt"

// matchOperator tries to match operator tokens from the beginning of words.
func matchOperator(words []string) (string, int, error) {
	operators := []struct {
		tokens []string
		op     string
	}{
		{[]string{"is", "not", "nil"}, "is not nil"},
		{[]string{"is", "nil"}, "is nil"},
		{[]string{"is", "greater", "than"}, "is greater than"},
		{[]string{"is", "less", "than"}, "is less than"},
		{[]string{"is", "at", "most"}, "is at most"},
		{[]string{"is", "at", "least"}, "is at least"},
		{[]string{"is", "not"}, "is not"},
		{[]string{"is", "in"}, "is in"},
		{[]string{"equals"}, "equals"},
	}

	for _, o := range operators {
		if len(words) < len(o.tokens) {
			continue
		}
		if tokensMatch(words, o.tokens) {
			return o.op, len(o.tokens), nil
		}
	}
	return "", 0, fmt.Errorf("no operator matched")
}
