package api

import (
	"sort"
	"testing"

	"github.com/drycc/controller-sdk-go/pkg/time"
)

func TestPtypesSorted(t *testing.T) {
	ptypes := Ptypes{
		{"web", "v1", "1/1", 1, 1, time.Time{}},
		{"cronjob", "v1", "1/1", 1, 1, time.Time{}},
		{"sleep", "v1", "1/1", 1, 1, time.Time{}},
	}

	sort.Sort(ptypes)

	expectedPtypeNames := []string{"cronjob", "sleep", "web"}

	for i, ptype := range ptypes {
		if expectedPtypeNames[i] != ptype.Name {
			t.Errorf("Expected ptypes to be sorted %v, Got %v", expectedPtypeNames[i], ptype.Name)
		}
	}
}
