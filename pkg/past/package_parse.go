package past

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func (o *parsingOptions) nameMatches(name string) bool {
	if o.excludeFileMask != nil {
		if o.excludeFileMask.MatchString(name) {
			return false
		}
	}

	if o.includeFileMask != nil {
		return o.includeFileMask.MatchString(name)
	}

	return strings.HasSuffix(name, ".go")
}

// ParsePackages reads entries in the root directory of the fsys and parses all Go files matching options opts.
// In the result here can be more than one package because Go allows to keep production and test code in one directory.
func ParsePackages(fsys fs.FS, path string, oo ...ParsingOption) ([]*Package, error) {
	opts := defaultParsingOptions()
	for _, o := range oo {
		o(&opts)
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	pm := make(map[string]*Package)
	for i, entry := range entries {
		if !opts.nameMatches(entry.Name()) {
			continue
		}

		filename := entry.Name()
		fsysFile, openErr := fsys.Open(filename)
		if openErr != nil {
			return nil, fmt.Errorf("open file %s: %w", filename, openErr)
		}

		file, parseErr := ParseFile(fsysFile, filepath.Join(path, filename), withParsingOptions(opts))
		if parseErr != nil {
			_ = fsysFile.Close()
			return nil, fmt.Errorf("parse file %s: %w", filename, parseErr)
		}

		closeErr := fsysFile.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("close file %s: %w", filename, closeErr)
		}

		p, ok := pm[file.PackageName]
		if !ok {
			p = &Package{
				Path:        path,
				Name:        file.PackageName,
				Files:       make([]*File, 0, len(entries)-i),
				parsingOpts: opts,
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
