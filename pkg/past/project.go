package past

import (
	"errors"
	"fmt"
	"strings"
)

// ProjectData serializable [Project] data.
type ProjectData struct {
	Modules map[string]map[string]*Module `yaml:"modules,omitempty"`
}

// Project is the top-level container for managing the set of [Module].
type Project struct {
	ProjectData `yaml:",inline"`

	parsingOpts parsingOptions
}

// NewProject creates a new instance of [Project] which contains no parsed data.
func NewProject(oo ...ParsingOption) *Project {
	opts := defaultParsingOptions()
	for _, o := range oo {
		o(&opts)
	}

	return &Project{
		ProjectData: ProjectData{
			Modules: make(map[string]map[string]*Module),
		},
		parsingOpts: opts,
	}
}

// AddModule binds the module to the project p and adds it to the list of project modules.
func (p *Project) AddModule(mod *Module, version string) {
	if p.Modules == nil {
		p.Modules = make(map[string]map[string]*Module)
	}

	versioned := p.Modules[mod.Name]
	if versioned == nil {
		versioned = make(map[string]*Module)
		p.Modules[mod.Name] = versioned
	}

	mod.bind(p)
	versioned[version] = mod
}

func (p *Project) FindPackage(pkgPath, moduleVersion string) (*Package, error) {
	mod, modFound := p.getPackageModule(pkgPath, moduleVersion)
	if !modFound {
		return nil, errors.New("module not found") // TODO better error handling
	}

	fmt.Println(mod)

	return nil, nil
}

func (p *Project) getPackageModule(pkgPath, moduleVersion string) (*Module, bool) {
	for name, versioned := range p.Modules {
		if name == pkgPath || strings.HasPrefix(pkgPath, name+"/") {
			mod, ok := versioned[moduleVersion]
			return mod, ok
		}
	}
	return nil, false
}
