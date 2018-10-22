/*
 * Copyright (c) 2018. Darwayne
 */

package services

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"strings"
)

type MemoryPackageStore struct {
	tree    *rbt.Tree
	depTree *rbt.Tree
}

type PackageData struct {
	Name      string
	Installed bool
	DependsOn []string
}

type DependencyData struct {
	Name       string
	RequiredBy *rbt.Tree
}

func NewMemoryPackageStore() *MemoryPackageStore {
	return &MemoryPackageStore{
		rbt.NewWithStringComparator(),
		rbt.NewWithStringComparator(),
	}
}

func (m *MemoryPackageStore) put(key string, value interface{}) {
	m.tree.Put(m.getKeyName(key), value)
}

func (m *MemoryPackageStore) remove(key string) {
	m.tree.Remove(m.getKeyName(key))
}

func (m *MemoryPackageStore) get(key string) *PackageData {
	result, found := m.tree.Get(m.getKeyName(key))

	if found {
		data, ok := result.(*PackageData)
		if ok {
			return data
		}
	}

	return nil

}

func (m *MemoryPackageStore) has(key string) bool {
	_, found := m.tree.Get(m.getKeyName(key))
	return found
}

func (m *MemoryPackageStore) Depend(packageNames ...string) bool {
	mainPackageName := packageNames[0]
	pkgKeyName := m.getKeyName(mainPackageName)
	data := &PackageData{Name: mainPackageName, DependsOn: packageNames[1:]}

	if m.has(mainPackageName) {
		return false
	}

	m.put(mainPackageName, data)
	for _, depName := range data.DependsOn {
		depKey := m.getKeyName(depName)
		rawDep, found := m.depTree.Get(depKey)

		// if dependency exits just append package name to its list
		if found {
			dep, ok := rawDep.(*DependencyData)
			if ok {
				dep.RequiredBy.Put(pkgKeyName, []string{})
			}
		} else {
			//create dependency
			tree := rbt.NewWithStringComparator()
			tree.Put(pkgKeyName, []string{})
			data := DependencyData{Name: depName, RequiredBy: tree}
			m.depTree.Put(depKey, &data)
		}
	}

	return true
}

func (m *MemoryPackageStore) Install(packageName string) (installed bool, installedDependencies []string) {
	pkg := m.get(packageName)

	if pkg != nil {
		if !pkg.Installed {
			installed = true
			pkg.Installed = true

			// if package has dependencies install dependencies
			for _, dep := range pkg.DependsOn {
				ins, _ := m.Install(dep)
				if ins {
					installedDependencies = append(installedDependencies, dep)
				}
			}
		}
	} else {
		m.put(packageName, &PackageData{Name: packageName, Installed: true})
		installed = true
	}

	return
}

func (m *MemoryPackageStore) Remove(packageName string) (removed bool, notInstalled bool, removedDependencies []string, inUseDependencies []string) {
	pkg := m.get(packageName)
	pkgKey := m.getKeyName(packageName)

	if pkg != nil {
		//check if package is depended on by other packages
		_, found := m.depTree.Get(pkgKey)
		if found {
			return
		}
		removed = true

		// check if dependencies are in use by other packages
		for _, depName := range pkg.DependsOn {
			depKey := m.getKeyName(depName)
			rawDep, found := m.depTree.Get(depKey)
			if found {
				dep, ok := rawDep.(*DependencyData)
				if ok {
					dep.RequiredBy.Remove(pkgKey)
					if dep.RequiredBy.Size() > 0 {
						inUseDependencies = append(inUseDependencies, dep.Name)
					} else {
						removedDependencies = append(removedDependencies, dep.Name)
						m.depTree.Remove(depKey)
					}
				}
			}
		}
	} else {
		notInstalled = true
	}
	return
}

func (m *MemoryPackageStore) List() (results []string) {
	for _, val := range m.tree.Values() {
		pkg, ok := val.(*PackageData)
		if ok && pkg != nil && pkg.Installed {
			results = append(results, pkg.Name)
		}
	}

	return
}

func (m *MemoryPackageStore) getKeyName(pkg string) string {
	return strings.ToLower(pkg)
}
