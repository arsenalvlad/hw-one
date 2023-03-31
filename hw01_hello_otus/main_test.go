package main_test

import (
	"testing"

	"golang.org/x/example/stringutil"
)

func TestReverse(t *testing.T) {
	t.Parallel()

	for _, value := range []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	} {
		got := stringutil.Reverse(value.in)
		if got != value.want {
			t.Errorf("Reverse(%q) == %q, want %q", value.in, got, value.want)
		}
	}
}
