/*
 * Copyright (c) 2018. Darwayne
 */

package services

import (
	database "github.com/darwayne/gopkg/pkg/db"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"testing"
)

func cleanDB(db *gorm.DB) {
	db.Exec("DROP SCHEMA IF EXISTS public CASCADE")
	db.Exec("CREATE SCHEMA public")
	database.RunMigrations(db)
}

func getDB() *gorm.DB {
	connStr := "host=localhost port=5432 dbname=gopkg_test user=gopkg password=pass sslmode=disable"
	return database.GetDB(connStr)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		return contains(b, v)
	}
	return true
}

func TestDBPackageStore_Depend(t *testing.T) {
	// set up test cases
	var tests = []struct {
		params           []string
		expectedResponse bool
	}{
		{[]string{"hey"}, true},
		{[]string{"hey"}, false},
		{[]string{"yo", "hey", "how"}, true},
	}

	db := getDB()
	cleanDB(db)

	store := NewDBPackageStore(db)

	for _, test := range tests {
		result := store.Depend(test.params...)
		if result != test.expectedResponse {
			t.Errorf("Expected Depend(%v) = %v but got %v",
				test.params, test.expectedResponse, result)
		}
	}
}

func TestDBPackageStore_Remove(t *testing.T) {
	db := getDB()
	cleanDB(db)

	store := NewDBPackageStore(db)

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

func TestDBPackageStore_List(t *testing.T) {
	db := getDB()
	cleanDB(db)

	// set up test cases
	var tests = []struct {
		pkgsToInstall []string
	}{
		{[]string{"Just", "a", "test"}},
	}

	for _, test := range tests {
		cleanDB(db)
		store := NewDBPackageStore(db)

		for _, pkg := range test.pkgsToInstall {
			store.Install(pkg)
		}

		result := store.List()

		if !Equal(result, test.pkgsToInstall) {
			t.Errorf("Expected List() = %v but got %v", test.pkgsToInstall, result)
		}
	}
}
