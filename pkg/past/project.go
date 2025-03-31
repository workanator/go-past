package past

// Project is the top-level container for managing the set of [Module].
type Project struct {
	Modules map[string]*Module
}

// NewProject creates a new instance of [Project] which contains no parsed data.
func NewProject() *Project {
	return &Project{
		Modules: make(map[string]*Module),
	}
}

// AddModule binds the module to the project p and adds it to the list of project modules.
func (p *Project) AddModule(mod *Module, version string) {
	if p.Modules == nil {
		p.Modules = make(map[string]*Module)
	}

	mod.bind(p)
	p.Modules[mod.Name+"@"+version] = mod
}
