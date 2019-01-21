package main

import "testing"

type TestCase struct {
	actual   string
	expected string
}

func TestNormalize(t *testing.T) {
	testCases := []TestCase{
		{"234567890", "234567890"},
		{"23 456 7891", "234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"23-456-7894", "234567894"},
		{"(123)456-7892", "1234567892"},
	}

	for _, tc := range testCases {
		t.Run(tc.actual, func(t *testing.T) {
			actual := normalize(tc.actual)
			if actual != tc.expected {
				t.Errorf("got %s; expected %s", actual, tc.expected)
			}
		})
	}

}
