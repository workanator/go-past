package usages

import (
	"go/ast"
	"regexp"
	"strings"
)

var reSimplify = regexp.MustCompile(`^(?:\.\.\.)?(?:(?:<-)?chan\s+(?:<-)?|\[\d*]|\*|\s)*`)

// Type usage.
type Type string

func (u Type) String() string {
	return string(u)
}

// NumPointers returns the number of prefixed *.
func (u Type) NumPointers() int {
	n := 0
	for n < len(u) && u[n] == '*' {
		n++
	}
	return n
}

// TrimPointers trims all pointers and returns the exact type.
func (u Type) TrimPointers() Type {
	return Type(strings.TrimPrefix(u.String(), "*"))
}

// AddPointers adds pointer prefix with n pointers.
func (u Type) AddPointers(n int) Type {
	return Type(strings.Repeat("*", n) + u.String())
}

// PackageAndType splits the spec to the package name and the type name.
func (u Type) PackageAndType() (string, string) {
	parts := strings.Split(u.TrimPointers().String(), ".")
	if len(parts) == 1 {
		return "", parts[0]
	}
	return parts[0], parts[1]
}

// IsImported returns true if the spec contains dot.
func (u Type) IsImported() bool {
	return strings.Contains(u.String(), ".")
}

// IsExported returns true if the spec starts with the upper case.
func (u Type) IsExported() bool {
	return ast.IsExported(u.String())
}

// IsVariadic tests whether the type usage starts with `...`.
func (u Type) IsVariadic() bool {
	return strings.HasPrefix(u.String(), "...")
}

// TrimVariadic trims the leading ... if it presents.
func (u Type) TrimVariadic() Type {
	return Type(strings.TrimPrefix(u.String(), "..."))
}

// Simplify trims all pointers, slices, channesl, etc. and returns the last most type.
func (u Type) Simplify() Type {
	return Type(reSimplify.ReplaceAllString(u.String(), ""))
}
