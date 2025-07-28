package query

import (
	"regexp/syntax"
	"strings"
)

// makeInsensitive takes a regular expression and returns lowercased version of the expression
func makeInsensitive(re string) (string, error) {
	r, err := syntax.Parse(re, syntax.Perl)
	if err != nil {
		return "", err
	}
	r = r.Simplify()
	expr := []*syntax.Regexp{r}
	for len(expr) != 0 {
		e := expr[len(expr)-1]
		expr = expr[:len(expr)-1]
		expr = append(expr, e.Sub...)
		switch e.Op {
		case syntax.OpLiteral:
			e.Rune = []rune(strings.ToLower(string(e.Rune)))
		}
	}
	return r.String(), nil
}
