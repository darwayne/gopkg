package services

import (
	"testing"
)

func TestMemoryPackageStore_Depend(t *testing.T) {
	// set up test cases
	var tests = []struct {
		params           []string
		expectedResponse bool
	}{
		{[]string{"hey"}, true},
		{[]string{"hey"}, false},
		{[]string{"yo", "hey", "how"}, true},
	}

	store := NewMemoryPackageStore()

	for _, test := range tests {
		result := store.Depend(test.params...)
		if result != test.expectedResponse {
			t.Errorf("Expected Depend(%v) = %v but got %v",
				test.params, test.expectedResponse, result)
		}
	}
}

func TestMemoryPackageStore_Remove(t *testing.T) {
	store := NewMemoryPackageStore()

	store.Depend([]string{"yo", "son", "k"}...)
	store.Install("yo")
	store.Depend("Hmm", "son")
	store.Install("Hmm")

	// set up test cases
	var tests = []struct {
		pkgToRemove         string
		removed             bool
		notInstalled        bool
		removedDependencies []string
		inUseDependencies   []string
	}{
		{
			"yo",
			true,
			false,
			[]string{"k"},
			[]string{"son"},
		},
		{
			"Hmm",
			true,
			false,
			[]string{"son"},
			[]string{},
		},
		{
			"doesn't exit",
			false,
			true,
			[]string{},
			[]string{},
		},
	}

	for _, test := range tests {
		removed, notInstalled, removedDependencies, inUseDependencies := store.Remove(test.pkgToRemove)

		if removed != test.removed || notInstalled != test.notInstalled ||
			!Equal(removedDependencies, test.removedDependencies) ||
			!Equal(inUseDependencies, test.inUseDependencies) {
			t.Errorf("Expected Remove(%v) = %v, %v, %v, %v but got %v, %v, %v, %v",
				test.pkgToRemove, test.removed, test.notInstalled, test.removedDependencies, test.inUseDependencies,
				removed, notInstalled, removedDependencies, inUseDependencies)
		}
	}
}

func TestMemoryPackageStore_List(t *testing.T) {
	// set up test cases
	var tests = []struct {
		pkgsToInstall []string
	}{
		{[]string{"Just", "a", "test"}},
	}

	for _, test := range tests {
		store := NewMemoryPackageStore()

		for _, pkg := range test.pkgsToInstall {
			store.Install(pkg)
		}

		result := store.List()

		if !Equal(result, test.pkgsToInstall) {
			t.Errorf("Expected List() = %v but got %v", test.pkgsToInstall, result)
		}
	}
}
