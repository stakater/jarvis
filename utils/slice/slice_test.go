package slice

import "testing"

func TestContains(t *testing.T) {
	inputSlice := []string{"aa", "bb"}

	if !Contains(inputSlice, "bb") {
		t.Errorf("Contains didn't find the element as expected")
	}
}
