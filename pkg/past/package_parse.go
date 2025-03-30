package past

import (
	"fmt"
	"go/token"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

// ParsePackageOptions options for package parsing.
type ParsePackageOptions struct {
	FileSet  *token.FileSet
	BasePath string
	Path     string
	FileMask *regexp.Regexp
}

func (o *ParsePackageOptions) normalize() *ParsePackageOptions {
	if o == nil {
		return &ParsePackageOptions{
			Path: ".",
		}
	}

	if o.Path == "" {
		o.Path = "."
	}

	return o
}

func (o *ParsePackageOptions) nameMatches(name string) bool {
	if o == nil || o.FileMask == nil {
		return strings.HasSuffix(name, ".go")
	}
	return o.FileMask.MatchString(name)
}

// ParsePackages reads entries in the root directory of the fsys and parses all Go files matching options opts.
// In the result here can be more than one package because Go allows to keep production and test code in one directory.
func ParsePackages(fsys fs.FS, opts *ParsePackageOptions) ([]*Package, error) {
	opts = opts.normalize()

	entries, err := fs.ReadDir(fsys, opts.Path)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	pm := make(map[string]*Package)
	for i, entry := range entries {
		if !opts.nameMatches(entry.Name()) {
			continue
		}

		path := filepath.Join(opts.Path, entry.Name())
		fsysFile, openErr := fsys.Open(path)
		if openErr != nil {
			return nil, fmt.Errorf("open file %s: %w", path, openErr)
		}

		file, parseErr := ParseFile(fsysFile, &ParseFileOptions{
			FileSet:  opts.FileSet,
			Filename: path,
		})
		if parseErr != nil {
			_ = fsysFile.Close()
			return nil, fmt.Errorf("parse file %s: %w", path, parseErr)
		}

		p, ok := pm[file.PackageName]
		if !ok {
			p = &Package{
				Path:  opts.BasePath,
				Name:  file.PackageName,
				Files: make([]*File, 0, len(entries)-i),
			}
			pm[file.PackageName] = p
		}

		p.AddFiles(file)
	}

	pp := make([]*Package, 0, len(pm))
	for _, p := range pm {
		pp = append(pp, p)
	}

	return pp, nil
}
