package past

// Location descriptive package location. Here import and path are relative to the module root.
type Location struct {
	PackageName    string `yaml:"package_name"`
	RelativeImport string `yaml:"relative_import"`
	RelativePath   string `yaml:"relative_path"`
}
