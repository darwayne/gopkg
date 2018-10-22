/*
 * Copyright (c) 2018. Darwayne
 */

package interfaces

type PackageStore interface {
	Depend(packageNames ...string) bool
	Install(packageName string) (installed bool, installedDependencies []string)
	Remove(packageName string) (removed bool, notInstalled bool, removedDependencies []string, inUseDependencies []string)
	List() []string
}
