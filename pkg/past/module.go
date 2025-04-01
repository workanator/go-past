package past

import (
	"cmp"
	"fmt"
	"github.com/denormal/go-gitignore"
	"golang.org/x/mod/modfile"
	"io/fs"
	"path/filepath"
	"strings"
)

// ModuleData serializable [Module] data.
type ModuleData struct {
	Name     string              `yaml:"name,omitempty"`
	Packages map[string]*Package `yaml:"packages,omitempty"`
}

// Module represents the Go module which is the directory with the go.mod file, manages the set of [Package].
type Module struct {
	ModuleData `yaml:",inline"`
	Path       string

	parsingOpts   parsingOptions
	owningProject *Project
	fs            fs.FS
	goMod         *modfile.File
	gitIgnore     gitignore.GitIgnore
}

// bind attaches the module to the owning [Project].
func (m *Module) bind(owningProject *Project) {
	m.owningProject = owningProject
}

// Project returns the [Project] which owns the module.
func (m *Module) Project() *Project {
	return m.owningProject
}

// AddPackage binds the package to the module m and adds it to the list of module packages at the import path relative to
// the module root, for example, if the module is github.com/author/go-project and the package import path is
// github.com/author/go-project/pkg/awesome then the relative path is pkg/awesome.
func (m *Module) AddPackage(pkg *Package, relativeImport string) {
	if m.Packages == nil {
		m.Packages = make(map[string]*Package)
	}

	pkg.bind(m)
	m.Packages[relativeImport] = pkg
}

func (m *Module) FindPackage(relativePath string) (*Package, error) {
	if pkg, ok := m.Packages[relativePath]; ok {
		return pkg, nil
	}

	parsingOpts := make([]ParsingOption, 0, 2)
	parsingOpts = append(parsingOpts, withParsingOptions(m.parsingOpts))
	if m.gitIgnore != nil {
		parsingOpts = append(parsingOpts, WithGitIgnore(m.gitIgnore))
	}

	if relativePath == "" {
		pkgs, parseErr := ParsePackages(m.fs, ".", parsingOpts...)
		if parseErr != nil {
			return nil, fmt.Errorf("parse packages: %w", parseErr)
		}

		if len(pkgs) == 0 {
			words := strings.Split(m.Name, "/")
			pkgs = []*Package{
				NewPackage("", words[len(words)-1], parsingOpts...),
			}
		}

		for _, pkg := range pkgs {
			m.AddPackage(pkg, "")
		}
	}

	first := true
	last := ""
	for {
		closestPkg, closestPkgPath, _ := m.findClosesPackage(relativePath)
		subFsys, err := fs.Sub(m.fs, cmp.Or(closestPkgPath, "."))
		if err != nil {
			return nil, fmt.Errorf("make sub-fs for %s: %w", closestPkgPath, err)
		}

		var pkgRelativePath string
		if closestPkg != nil {
			pkgRelativePath = closestPkg.Path
		}

		err = m.parseAndAddPackagesAtPath(subFsys, pkgRelativePath, closestPkgPath)
		if err != nil {
			return nil, fmt.Errorf("parse packages: %w", err)
		}

		if first {
			first = false
		} else if last == pkgRelativePath {
			return nil, fmt.Errorf("not found")
		}

		last = pkgRelativePath
	}

	return nil, nil
}

func (m *Module) findClosesPackage(relativePath string) (*Package, string, []string) {
	var lastPkg *Package
	var accum string
	crumbs := strings.Split(relativePath, "/")
	for i := range crumbs {
		if accum != "" {
			accum += "/"
		}
		accum += crumbs[i]

		if pkg, ok := m.Packages[accum]; ok {
			lastPkg = pkg
			continue
		}

		return lastPkg, strings.Join(crumbs[:i], "/"), crumbs[i:]
	}

	return nil, "", nil
}

func (m *Module) parseAndAddPackagesAtPath(fsys fs.FS, basePath, baseRelativeImport string) error {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subFsys, subErr := fs.Sub(fsys, entry.Name())
		if subErr != nil {
			return fmt.Errorf("make sub-fs for %s: %w", entry.Name(), subErr)
		}

		pkgPath := filepath.Join(basePath, entry.Name())
		pkgs, parseErr := ParsePackages(
			subFsys,
			pkgPath,
			withParsingOptions(m.parsingOpts),
		)
		if parseErr != nil {
			return fmt.Errorf("parse packages in %s: %w", entry.Name(), parseErr)
		}

		for _, pkg := range pkgs {
			var pkgImport string
			if baseRelativeImport == "" {
				pkgImport = pkg.Name
			} else {
				pkgImport = baseRelativeImport + "/" + pkg.Name
			}

			m.AddPackage(pkg, pkgImport)
		}

		if len(pkgs) == 0 {
			var pkgImport string
			if baseRelativeImport == "" {
				pkgImport = entry.Name()
			} else {
				pkgImport = baseRelativeImport + "/" + entry.Name()
			}

			m.AddPackage(
				NewPackage(pkgPath, entry.Name(), withParsingOptions(m.parsingOpts)),
				pkgImport,
			)
		}
	}

	return nil
}
