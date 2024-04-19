package external

import "testing"

func TestConvertDurationToMinutes(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"PT1H27M", 87},
		{"PT4H6M", 246},
		{"PT2H", 120},
		{"PT10M", 10},
	}

	for _, test := range tests {
		if result, err := ConvertDurationToMinutes(test.input); result != test.expected {
			t.Errorf("%s: expected %d, got %d (error: %v)", test.input, test.expected, result, err)
		}
	}
}
