package past

import (
	"github.com/workanator/go-past/pkg/past/decls"
)

// File is the Go source file with declarations in it.
type File struct {
	Filename    string
	PackageName string
	Doc         string
	Imports     []*Import
	Function    map[string]*Function
	Struct      map[string]*Struct
	Interface   map[string]*Interface
	Variable    map[string]*Variable
	Constant    map[string]*Constant

	owningPackage *Package
}

// bind attaches the file to the owning [Package].
func (f *File) bind(owningPackage *Package) {
	f.owningPackage = owningPackage
}

// Package returns the [Package] which owns the file.
func (f *File) Package() *Package {
	return f.owningPackage
}

// AddImports create [Import]s from the declarations and adds to the list of file imports.
func (f *File) AddImports(decls ...decls.Import) {
	for _, d := range decls {
		item := &Import{
			Decl: d,
		}
		item.bind(f)

		f.Imports = append(f.Imports, item)
	}
}

// AddFunctions create [Function]s from the declarations and adds to the list of file functions.
func (f *File) AddFunctions(decls ...decls.Func) {
	for _, d := range decls {
		item := &Function{
			Decl: d,
		}
		item.bind(f)

		f.Function[item.Decl.Name] = item
	}
}

// AddStructs create [Struct]s from the declarations and adds to the list of file structs.
func (f *File) AddStructs(decls ...decls.Struct) {
	for _, d := range decls {
		item := &Struct{
			Decl: d,
		}
		item.bind(f)

		f.Struct[item.Decl.Name] = item
	}
}

// AddInterfaces create [Interface]s from the declarations and adds to the list of file interfaces.
func (f *File) AddInterfaces(decls ...decls.Interface) {
	for _, d := range decls {
		item := &Interface{
			Decl: d,
		}
		item.bind(f)

		f.Interface[item.Decl.Name] = item
	}
}

// AddVariables create [Variable]s from the declarations and adds to the list of file variables.
func (f *File) AddVariables(decls ...decls.Value) {
	for _, d := range decls {
		item := &Variable{
			Decl: d,
		}
		item.bind(f)

		f.Variable[item.Decl.Name] = item
	}
}

// AddConstants create [Constant]s from the declarations and adds to the list of file constants.
func (f *File) AddConstants(decls ...decls.Value) {
	for _, d := range decls {
		item := &Constant{
			Decl: d,
		}
		item.bind(f)

		f.Constant[item.Decl.Name] = item
	}
}
