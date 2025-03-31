package past

import "io/fs"

// Module represents the Go module which is the directory with the go.mod file, manages the set of [Package].
type Module struct {
	Path     string
	Name     string
	Packages map[string]*Package

	parsingOpts   parsingOptions
	owningProject *Project
	fs            fs.FS
}

// bind attaches the module to the owning [Project].
func (m *Module) bind(owningProject *Project) {
	m.owningProject = owningProject
}

// Project returns the [Project] which owns the module.
func (m *Module) Project() *Project {
	return m.owningProject
}

// AddPackage binds the package to the module m and adds it to the list of module packages at the path relative to
// the module root, for example, if the module is github.com/author/go-project and the package import path is
// github.com/author/go-project/pkg/awesome then the relative path is pkg/awesome.
func (m *Module) AddPackage(pkg *Package, relativePath string) {
	if m.Packages == nil {
		m.Packages = make(map[string]*Package)
	}

	pkg.bind(m)
	m.Packages[relativePath] = pkg
}
