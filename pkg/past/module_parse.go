package past

import (
	"fmt"
	"github.com/denormal/go-gitignore"
	"golang.org/x/mod/modfile"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ParseModule parses the go.mod file and extracts module information from it. fsys should have the root where go.mod
// file is located.
func ParseModule(fsys fs.FS, path string, oo ...ParsingOption) (*Module, error) {
	opts := defaultParsingOptions()
	for _, o := range oo {
		o(&opts)
	}

	// read and parse go.mod
	mf, err := readAndParseGoModule(fsys)
	if err != nil {
		return nil, fmt.Errorf("read and parse go.mod: %w", err)
	}

	// read and parse .gitignore if any
	var gf gitignore.GitIgnore
	if opts.followGitIgnoreRules {
		gf, err = readAndParseGitIgnore(fsys, ".gitignore")
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("read and parse .gitignore: %w", err)
		}
	}

	return &Module{
		ModuleData: ModuleData{
			Name:     mf.Module.Mod.Path,
			Packages: make(map[string]*Package),
		},
		Path:        path,
		parsingOpts: opts,
		fs:          fsys,
		goMod:       mf,
		gitIgnore:   gf,
	}, nil
}

func readAndParseGoModule(fsys fs.FS) (*modfile.File, error) {
	file, err := fsys.Open("go.mod")
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	mf, err := parseGoModule(filepath.Base("go.mod"), file)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("read: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("close: %w", err)
	}

	return mf, nil
}

func parseGoModule(name string, r io.Reader) (*modfile.File, error) {
	p, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	f, err := modfile.Parse(name, p, nil)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	return f, nil
}

func readAndParseGitIgnore(fsys fs.FS, path string) (gitignore.GitIgnore, error) {
	file, err := fsys.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	gf := gitignore.New(file, filepath.Dir(path), nil)

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("close: %w", err)
	}

	return gf, nil
}
