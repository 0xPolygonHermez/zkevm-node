package pgstatestorage

import "testing"

func Test_valuePlaceholdersForRow(t *testing.T) {
	tcs := []struct {
		description  string
		order, total int
		expected     string
	}{
		{
			description: "single element, first row",
			order:       0,
			total:       1,
			expected:    "($1)",
		},
		{
			description: "multiple elements, first row",
			order:       0,
			total:       3,
			expected:    "($1,$2,$3)",
		},
		{
			description: "single element, non-first row",
			order:       9,
			total:       1,
			expected:    "($10)",
		},
		{
			description: "multiple elements, non-first row",
			order:       9,
			total:       3,
			expected:    "($28,$29,$30)",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			actual := valuePlaceholdersForRow(tc.order, tc.total)

			if actual != tc.expected {
				t.Fatalf("Actual value placeholders %q don't match expected placeholders %q", actual, tc.expected)
			}
		})
	}
}
