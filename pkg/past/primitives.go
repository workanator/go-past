package past

import "github.com/workanator/go-past/pkg/past/decls"

// Import declaration of the file-level import.
type Import struct {
	Decl decls.Import

	owningFile *File
}

// bind attaches the import declaration to the containing file.
func (d *Import) bind(owningFile *File) {
	d.owningFile = owningFile
}

// File returns the [File] which contains the import declaration.
func (d *Import) File() *File {
	return d.owningFile
}

// Variable declaration of the file-level variable.
type Variable struct {
	Decl decls.Value

	owningFile *File
}

// bind attaches the variable declaration to the containing file.
func (d *Variable) bind(owningFile *File) {
	d.owningFile = owningFile
}

// File returns the [File] which contains the variable declaration.
func (d *Variable) File() *File {
	return d.owningFile
}

// Constant declaration of the file-level constant.
type Constant struct {
	Decl decls.Value

	owningFile *File
}

// bind attaches the constant declaration to the containing file.
func (d *Constant) bind(owningFile *File) {
	d.owningFile = owningFile
}

// File returns the [File] which contains the constant declaration.
func (d *Constant) File() *File {
	return d.owningFile
}

// Function declaration of the file-level function.
type Function struct {
	Decl decls.Func

	owningFile *File
}

// bind attaches the function declaration to the containing file.
func (d *Function) bind(owningFile *File) {
	d.owningFile = owningFile
}

// File returns the [File] which contains the function declaration.
func (d *Function) File() *File {
	return d.owningFile
}

// Struct declaration of the file-level struct.
type Struct struct {
	Decl decls.Struct

	owningFile *File
}

// bind attaches the struct declaration to the containing file.
func (d *Struct) bind(owningFile *File) {
	d.owningFile = owningFile
}

// File returns the [File] which contains the struct declaration.
func (d *Struct) File() *File {
	return d.owningFile
}

// Interface declaration of the file-level interface.
type Interface struct {
	Decl decls.Interface

	owningFile *File
}

// bind attaches the interface declaration to the containing file.
func (d *Interface) bind(owningFile *File) {
	d.owningFile = owningFile
}

// File returns the [File] which contains the interface declaration.
func (d *Interface) File() *File {
	return d.owningFile
}
