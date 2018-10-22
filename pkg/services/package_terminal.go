/*
 * Copyright (c) 2018. Darwayne
 */

package services

import (
	"bufio"
	"fmt"
	"github.com/darwayne/gopkg/pkg/interfaces"
	"os"
	"regexp"
	"strings"
)

type PackageTerminal struct {
	store         interfaces.PackageStore
	commandsRegex *regexp.Regexp
	endRegex      *regexp.Regexp
}

func NewPackageTerminal(store interfaces.PackageStore) *PackageTerminal {
	regexStr := "(?i)\\s*(depend|install|list|end|remove)"
	return &PackageTerminal{
		store:         store,
		commandsRegex: regexp.MustCompile(regexStr),
		endRegex:      regexp.MustCompile("(?i)^\\s*end"),
	}
}

func (p *PackageTerminal) hasSupportedCommand(line string) bool {
	return p.commandsRegex.MatchString(line)
}

func (p *PackageTerminal) isEnd(line string) bool {
	return p.endRegex.MatchString(line)
}

func (p *PackageTerminal) sanitizeLine(line string) string {
	result := strings.Replace(line, "  ", " ", -1)
	result = strings.TrimSpace(result)

	return result
}

func (p *PackageTerminal) getCommandOutput(line string) (result string) {
	line = p.sanitizeLine(line)
	args := strings.Split(line, " ")
	cmd := strings.ToLower(args[0])

	switch cmd {
	case "depend":
		p.store.Depend(args[1:]...)
	case "install":
		pkgName := strings.TrimSpace(args[1])
		installed, installedDependencies := p.store.Install(pkgName)

		if installed {
			items := append(installedDependencies, pkgName)
			for _, pkg := range items {
				result += fmt.Sprintf("%s successfully installed\n", pkg)
			}
		} else {
			result += fmt.Sprintf("%s is already installed\n", pkgName)
		}
	case "remove":
		pkgName := strings.TrimSpace(args[1])
		removed, notInstalled, removedDependencies, inUseDependencies := p.store.Remove(pkgName)
		if removed {
			result += fmt.Sprintf("%s successfully removed\n", pkgName)
			for _, dep := range inUseDependencies {
				result += fmt.Sprintf("%s is still needed\n", dep)
			}

			for _, dep := range removedDependencies {
				result += fmt.Sprintf("%s is no longer needed\n", dep)
			}
		} else if notInstalled {
			result += fmt.Sprintf("%s is not installed\n", pkgName)
		} else {
			result += fmt.Sprintf("%s is still needed\n", pkgName)
		}

	case "list":
		for _, pkg := range p.store.List() {
			result += fmt.Sprintf("%s\n", pkg)
		}

	}

	return result
}

func (p *PackageTerminal) Read() {
	reader := bufio.NewReader(os.Stdin)
	output := ""

	for {
		text, _ := reader.ReadString('\n')
		if p.hasSupportedCommand(text) {
			output += text
			str := p.getCommandOutput(text)
			output += str
			if p.isEnd(text) {
				break
			}
		}
	}

	fmt.Println(output)
}
