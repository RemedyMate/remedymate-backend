package validation_test

import (
	"testing"

	"remedymate-backend/util/validation"
)

func TestISO3166Alpha2(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"ET", true}, {"US", true}, {"GB", true},
		{"Et", false}, {"E", false}, {"ETH", false}, {"1T", false},
	}
	for _, c := range cases {
		if got := validation.IsISO3166Alpha2(c.in); got != c.want {
			t.Errorf("IsISO3166Alpha2(%q)=%v want %v", c.in, got, c.want)
		}
	}
}

func TestISO6391(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"en", true}, {"am", true}, {"es", true},
		{"EN", false}, {"e", false}, {"eng", false},
	}
	for _, c := range cases {
		if got := validation.IsISO6391(c.in); got != c.want {
			t.Errorf("IsISO6391(%q)=%v want %v", c.in, got, c.want)
		}
	}
}

func TestE164(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"+251911123456", true},
		{"+12025550123", true},
		{"251911123456", false},
		{"+001234", false},
		{"+1234567", false},
	}
	for _, c := range cases {
		if got := validation.IsE164Phone(c.in); got != c.want {
			t.Errorf("IsE164Phone(%q)=%v want %v", c.in, got, c.want)
		}
	}
}
