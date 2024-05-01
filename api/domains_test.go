package api

import (
	"sort"
	"testing"
)

func TestDomainsSorted(t *testing.T) {
	domains := Domains{
		{"Alpha", "", "gamma.example.com", "web", "", ""},
		{"Alpha", "", "alpha1.example.com", "web", "", ""},
		{"Alpha", "", "zulu.example.com", "web", "", ""},
		{"Alpha", "", "delta.example.com", "web", "", ""},
	}

	sort.Sort(domains)
	expectedDomains := []string{"alpha1.example.com", "delta.example.com", "gamma.example.com", "zulu.example.com"}

	for i, domain := range domains {
		if expectedDomains[i] != domain.Domain {
			t.Errorf("Expected domains to be sorted %v, Got %v at index %v", expectedDomains[i], domain.Domain, i)
		}
	}
}
