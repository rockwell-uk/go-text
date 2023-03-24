package text

import (
	"reflect"
	"testing"
)

func TestSplitStringInTwo(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			"St Marychurch",
			[]string{"St Marychurch"},
		},
		{
			"TORQUAY",
			[]string{"TORQUAY"},
		},
		{
			"Babbacombe Bay",
			[]string{"Babbacombe Bay"},
		},
		{
			"Great Hill",
			[]string{"Great Hill"},
		},
		{
			"Combe Palford",
			[]string{"Combe", "Palford"},
		},
		{
			"Pilsworth Road",
			[]string{"Pilsworth", "Road"},
		},
		{
			"Royle Barn Road",
			[]string{"Royle", "Barn Road"},
		},
		{
			"6 Barn Road",
			[]string{"6 Barn Road"},
		},
	}

	for _, tt := range tests {
		actual := SplitStringInTwo(tt.input, ShouldSplit)

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("%v: Expected [%#v]\nGot [%#v]", tt.input, tt.expected, actual)
		}
	}
}
