package api

import (
	"sort"
	"testing"
)

func TestPtypesSorted(t *testing.T) {
	ptypes := Ptypes{
		{"web", "v1", "1/1", 1, 1, "", false},
		{"cronjob", "v1", "1/1", 1, 1, "", false},
		{"sleep", "v1", "1/1", 1, 1, "", false},
	}

	sort.Sort(ptypes)

	expectedPtypeNames := []string{"cronjob", "sleep", "web"}

	for i, ptype := range ptypes {
		if expectedPtypeNames[i] != ptype.Name {
			t.Errorf("Expected ptypes to be sorted %v, Got %v", expectedPtypeNames[i], ptype.Name)
		}
	}
}
