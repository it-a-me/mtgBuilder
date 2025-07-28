package query

import (
	"regexp"
	"regexp/syntax"
	"strings"
	"testing"
)

func reEqual(a, b string) (bool, error) {
	expr1, err := syntax.Parse(a, syntax.Perl)
	if err != nil {
		return false, err
	}
	expr2, err := syntax.Parse(b, syntax.Perl)
	if err != nil {
		return false, err
	}
	return expr1.Equal(expr2), nil
}

func TestMakeInsensitive(t *testing.T) {
	cases := map[string]string{
		"AB":       "ab",
		"aB":       "ab",
		"r[A-c]":   "r[A-c]",
		`cowid:\D`: `cowid:\D`,
	}
	for expr, expected := range cases {
		got, err := makeInsensitive(expr)
		if err != nil {
			t.Fatal(err)
		}
		equal, err := reEqual(got, expected)
		if err != nil {
			t.Fatal(err)
		}
		if !equal {
			t.Fatalf("expected `%s`, got `%s`", expected, got)
		}
	}
}

func FuzzMakeInsensitive(f *testing.F) {
	f.Add(`abc[A-c]\D`, "abcb@")
	f.Fuzz(func(t *testing.T, re string, query string) {
		simple, err := regexp.Compile("(?i)" + re)
		if err != nil {
			t.Skip()
		}
		manual, err := makeInsensitive(re)
		if err != nil {
			t.Fatalf(`failed to make '%s' insensitive`, re)
		}
		compiled, err := regexp.Compile(manual)
		if err != nil {
			t.Fatalf(`failed to compile '%s'`, re)
		}
		expected := simple.MatchString(query)
		got := compiled.MatchString(strings.ToLower(query))
		if expected != got {
			t.Fatalf("expected matching `%s` on '%s', to return %t but got %t", re, query, expected, got)
		}
	})
}
