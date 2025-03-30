package past

// Package represents the Go package which is the directory with source files, manages [File]s.
type Package struct {
	Path  string
	Name  string
	Files []*File

	owningModule *Module
}

// bind attaches the package to the owning [Module].
func (p *Package) bind(owningModule *Module) {
	p.owningModule = owningModule
}

// Module returns the [Module] which owns the package.
func (p *Package) Module() *Module {
	return p.owningModule
}

// AddFiles binds files to the package p and adds them to the list package files.
func (p *Package) AddFiles(files ...*File) {
	for _, file := range files {
		file.bind(p)
		p.Files = append(p.Files, file)
	}
}
