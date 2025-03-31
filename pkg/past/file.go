package past

import (
	"github.com/workanator/go-past/pkg/past/decls"
)

// File is the Go source file with declarations in it.
type File struct {
	Path        string
	PackageName string
	Doc         string
	Imports     []*Import
	Functions   map[string]*Function
	Structs     map[string]*Struct
	Interfaces  map[string]*Interface
	Variables   map[string]*Variable
	Constants   map[string]*Constant

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
	if len(f.Imports) == 0 {
		f.Imports = make([]*Import, 0, len(decls))
	}

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
	if len(f.Functions) == 0 {
		f.Functions = make(map[string]*Function)
	}

	for _, d := range decls {
		item := &Function{
			Decl: d,
		}
		item.bind(f)

		f.Functions[item.Decl.Name] = item
	}
}

// AddStructs create [Struct]s from the declarations and adds to the list of file structs.
func (f *File) AddStructs(decls ...decls.Struct) {
	if len(f.Structs) == 0 {
		f.Structs = make(map[string]*Struct)
	}

	for _, d := range decls {
		item := &Struct{
			Decl: d,
		}
		item.bind(f)

		f.Structs[item.Decl.Name] = item
	}
}

// AddInterfaces create [Interface]s from the declarations and adds to the list of file interfaces.
func (f *File) AddInterfaces(decls ...decls.Interface) {
	if len(f.Interfaces) == 0 {
		f.Interfaces = make(map[string]*Interface)
	}

	for _, d := range decls {
		item := &Interface{
			Decl: d,
		}
		item.bind(f)

		f.Interfaces[item.Decl.Name] = item
	}
}

// AddVariables create [Variable]s from the declarations and adds to the list of file variables.
func (f *File) AddVariables(decls ...decls.Value) {
	if len(f.Variables) == 0 {
		f.Variables = make(map[string]*Variable)
	}

	for _, d := range decls {
		item := &Variable{
			Decl: d,
		}
		item.bind(f)

		f.Variables[item.Decl.Name] = item
	}
}

// AddConstants create [Constant]s from the declarations and adds to the list of file constants.
func (f *File) AddConstants(decls ...decls.Value) {
	if len(f.Constants) == 0 {
		f.Constants = make(map[string]*Constant)
	}

	for _, d := range decls {
		item := &Constant{
			Decl: d,
		}
		item.bind(f)

		f.Constants[item.Decl.Name] = item
	}
}
