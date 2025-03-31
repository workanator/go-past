package past

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
)

// ParseModule parses the go.mod file and extracts module information from it. fsys should have the root where go.mod
// file is located.
func ParseModule(fsys fs.FS, path string, oo ...ParsingOption) (*Module, error) {
	opts := defaultParsingOptions()
	for _, o := range oo {
		o(&opts)
	}

	file, err := fsys.Open("go.mod")
	if err != nil {
		return nil, fmt.Errorf("open go.mod: %w", err)
	}

	moduleName, err := readGoModule(file)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("read go.mod: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("close go.mod: %w", err)
	}

	return &Module{
		Path:        path,
		Name:        moduleName,
		Packages:    make(map[string]*Package),
		parsingOpts: opts,
		fs:          fsys,
	}, nil
}

func readGoModule(r io.Reader) (string, error) {
	buf := bufio.NewReader(r)
	for {
		s, readErr := buf.ReadString('\n')
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return "", nil
			}
			return "", fmt.Errorf("read: %w", readErr)
		}

		if strings.HasPrefix(s, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(s, "module ")), nil
		}
	}
}
