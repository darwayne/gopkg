/*
 * Copyright (c) 2018. Darwayne
 */

package services

import (
	"fmt"
	"github.com/darwayne/gopkg/pkg/models"
	"github.com/jinzhu/gorm"
	"sort"
	"strings"
)

type DBPackageStore struct {
	db *gorm.DB
}

func NewDBPackageStore(db *gorm.DB) *DBPackageStore {
	return &DBPackageStore{db: db}
}

func (d *DBPackageStore) Depend(packageNames ...string) bool {
	mainPackageName := packageNames[0]
	mainPackage := models.Package{Name: mainPackageName}
	res := d.db.Create(&mainPackage)

	if res.Error != nil {
		return false
	}

	if len(packageNames) > 1 {
		queryFmt := "INSERT INTO packages (name) values %s ON CONFLICT DO NOTHING"
		values := make([]string, 0)
		args := make([]interface{}, 0)
		dependentPackageNames := packageNames[1:]

		for _, name := range dependentPackageNames {
			args = append(args, name)
			values = append(values, "?")
		}
		valuesStr := fmt.Sprintf("(%s)", strings.Join(values, "), ("))
		query := fmt.Sprintf(queryFmt, valuesStr)

		res = d.db.Exec(query, args...)

		if res.Error != nil {
			return false
		}

		res = d.db.Exec(`
	INSERT INTO package_dependencies (package_id, needed_package_id)
	SELECT ?, id from packages WHERE name in (?)
`, mainPackage.Id, dependentPackageNames)

		if res.Error != nil {
			return false
		}
	}

	return true
}

func (d *DBPackageStore) Install(packageName string) (installed bool, installedDependencies []string) {
	mainPackage := models.Package{}
	res := d.db.FirstOrCreate(&mainPackage, models.Package{Name: packageName})

	if res.Error != nil {
		return
	}

	// run query to see what dependencies need to be installed
	res = d.db.Raw(`
      SELECT p.name from package_dependencies pd
      INNER JOIN packages p ON p.id = pd.needed_package_id
      LEFT JOIN installed_packages ip ON ip.package_id = pd.needed_package_id
      WHERE  pd.package_id = ? AND ip.package_id IS NULL
`, mainPackage.Id).Pluck("name", &installedDependencies)

	if res.Error != nil {
		return
	}

	res = d.db.Exec(`
	INSERT INTO installed_packages (package_id)
	SELECT ?
	UNION ALL
	SELECT pd.needed_package_id from package_dependencies pd
    WHERE pd.package_id = ?
ON CONFLICT DO NOTHING
`, mainPackage.Id, mainPackage.Id)

	if res.Error != nil || res.RowsAffected == 0 {
		return
	} else if res.RowsAffected > 0 {
		installed = true
	}

	return
}

func inArray(arr []int, value int) bool {
	length := len(arr)
	i := sort.Search(length, func(i int) bool { return arr[i] >= value })
	return i < length && arr[i] == value
}

func (d *DBPackageStore) Remove(packageName string) (removed bool, notInstalled bool, removedDependencies []string, inUseDependencies []string) {
	totalUses := 0
	installedCount := 0

	mainPackage := models.Package{}

	res := d.db.Find(&mainPackage, models.Package{Name: packageName})

	if res.Error != nil || mainPackage.Id == 0 {
		notInstalled = true
		return
	}

	inUseDependencyIds := make([]int, 0)
	idsToRemove := []int{mainPackage.Id}

	// check if this package is installed
	d.db.Raw(`
	SELECT COUNT(*) from installed_packages
	WHERE package_id = ?
`, mainPackage.Id).Count(&installedCount)

	if installedCount == 0 {
		notInstalled = true
		return
	}

	// check if package is a dependency of another
	d.db.Raw(`
	SELECT COUNT(*) from package_dependencies
	WHERE needed_package_id = ?
`, mainPackage.Id).Count(&totalUses)

	// if another package depends on this package we cannot remove this package
	if totalUses > 0 {
		return
	}

	// get required packages that are in use by another package
	d.db.Raw(`SELECT pd.needed_package_id FROM package_dependencies pd
		INNER JOIN package_dependencies pd2 ON pd.package_id != pd2.package_id 
          AND pd2.needed_package_id = pd.needed_package_id
        WHERE pd.package_id = ?
        ORDER BY pd.needed_package_id ASC
`, mainPackage.Id).Pluck("needed_package_id", &inUseDependencyIds)

	dependencies := make([]models.Package, 0)
	// get all dependencies of this package
	rows, err := d.db.Raw(`
	SELECT p.id, p.name FROM package_dependencies pd
    INNER JOIN packages p ON p.id = pd.needed_package_id
    WHERE pd.package_id = ?
`, mainPackage.Id).Rows()

	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var pkg models.Package
			d.db.ScanRows(rows, &pkg)
			dependencies = append(dependencies, pkg)

			// check if dependency is in use
			if inArray(inUseDependencyIds, pkg.Id) {
				inUseDependencies = append(inUseDependencies, pkg.Name)
			} else {
				// if dependency not used we can also remove that dependency
				idsToRemove = append(idsToRemove, pkg.Id)
				removedDependencies = append(removedDependencies, pkg.Name)
			}
		}
	}

	//run delete command
	d.db.Exec("DELETE FROM packages WHERE id IN (?)", idsToRemove)

	removed = len(idsToRemove) > 0

	return
}

func (d *DBPackageStore) List() (results []string) {
	rows, err := d.db.Raw(`
SELECT p.* from packages p
INNER JOIN installed_packages ip ON ip.package_id = p.id
`).Rows()

	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var pkg models.Package
			d.db.ScanRows(rows, &pkg)
			results = append(results, pkg.Name)
		}
	}

	return

}
