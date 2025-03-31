package past

import (
	"fmt"
	"github.com/workanator/go-past/pkg/past/decls"
	"github.com/workanator/go-past/pkg/past/usages"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"strings"
)

// ParseFile reads top-level declarations from the Go file.
func ParseFile(r io.Reader, path string, oo ...ParsingOption) (*File, error) {
	opts := defaultParsingOptions()
	for _, o := range oo {
		o(&opts)
	}

	if opts.fileSet == nil {
		opts.fileSet = token.NewFileSet()
	}

	f, err := parser.ParseFile(opts.fileSet, path, r, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	file := &File{
		Path:        path,
		PackageName: f.Name.Name,
		Doc:         strings.Trim(f.Doc.Text(), "\n"),
		Imports:     make([]*Import, 0, len(f.Imports)),
		Functions:   make(map[string]*Function, len(f.Decls)),
		Structs:     make(map[string]*Struct, len(f.Decls)),
		Interfaces:  make(map[string]*Interface, len(f.Decls)),
		Variables:   make(map[string]*Variable, len(f.Decls)),
		Constants:   make(map[string]*Constant, len(f.Decls)),
	}

	for _, imp := range f.Imports {
		file.AddImports(parseAstImportDecl(imp))
	}

	for _, decl := range f.Decls {
		switch typedDecl := decl.(type) {
		case *ast.GenDecl:
			parseAstGenDecl(file, typedDecl)
		case *ast.FuncDecl:
			file.AddFunctions(parseAstFuncDecl(typedDecl))
		}
	}

	return file, nil
}

func parseAstImportDecl(spec *ast.ImportSpec) decls.Import {
	name := ""
	if spec.Name != nil {
		name = spec.Name.Name
	}

	return decls.Import{
		Path:    strings.Trim(spec.Path.Value, `"`),
		Name:    name,
		Doc:     strings.Trim(spec.Doc.Text(), "\n"),
		Comment: strings.Trim(spec.Comment.Text(), "\n"),
	}
}

func parseAstGenDecl(dst *File, gen *ast.GenDecl) {
	for i := range gen.Specs {
		valSpec, ok := gen.Specs[i].(*ast.ValueSpec)
		if ok {
			vals := parseAstValueSpec(valSpec)
			for _, vs := range vals {
				if gen.Tok == token.CONST {
					dst.AddConstants(vs)
				} else {
					dst.AddVariables(vs)
				}
			}

			continue
		}

		typeSpec, ok := gen.Specs[i].(*ast.TypeSpec)
		if !ok {
			continue
		}

		switch kind := typeSpec.Type.(type) {
		case *ast.StructType:
			dst.AddStructs(parseAstStructDecl(gen, typeSpec, kind))
		case *ast.InterfaceType:
			dst.AddInterfaces(parseAstInterfaceDecl(gen, typeSpec, kind))
		}
	}
}

func parseAstFuncDecl(decl *ast.FuncDecl) decls.Func {
	funcDecl := decls.Func{
		Name:     decl.Name.String(),
		Doc:      strings.Trim(decl.Doc.Text(), "\n"),
		Exported: decl.Name.IsExported(),
	}

	if decl.Recv != nil {
		fieldDecls := parseAstFieldDecl(decl.Recv.List[0])
		funcDecl.Receiver = &fieldDecls[0]
	}
	if decl.Type.TypeParams != nil && decl.Type.TypeParams.NumFields() > 0 {
		funcDecl.TypeParams = parseAstFieldList(decl.Type.TypeParams)
	}
	if decl.Type.Params != nil && decl.Type.Params.NumFields() > 0 {
		funcDecl.Params = parseAstFieldList(decl.Type.Params)
	}
	if decl.Type.Results != nil && decl.Type.Results.NumFields() > 0 {
		funcDecl.Results = parseAstFieldList(decl.Type.Results)
	}

	return funcDecl
}

func parseAstStructDecl(gen *ast.GenDecl, spec *ast.TypeSpec, decl *ast.StructType) decls.Struct {
	structDecl := decls.Struct{
		Name:     spec.Name.String(),
		Doc:      strings.Trim(gen.Doc.Text(), "\n"),
		Exported: spec.Name.IsExported(),
	}

	if decl.Fields != nil && decl.Fields.NumFields() > 0 {
		structDecl.Fields = make([]decls.Field, 0, decl.Fields.NumFields())
		for _, fld := range decl.Fields.List {
			fieldDecls := parseAstFieldDecl(fld)
			for _, fd := range fieldDecls {
				structDecl.Fields = append(structDecl.Fields, fd)
			}
		}
	}

	return structDecl
}

func parseAstInterfaceDecl(gen *ast.GenDecl, spec *ast.TypeSpec, decl *ast.InterfaceType) decls.Interface {
	ifaceDecl := decls.Interface{
		Name:     spec.Name.String(),
		Doc:      strings.Trim(gen.Doc.Text(), "\n"),
		Exported: spec.Name.IsExported(),
	}

	if decl.Methods != nil && decl.Methods.NumFields() > 0 {
		ifaceDecl.Embeds = make([]usages.Type, 0, decl.Methods.NumFields())
		ifaceDecl.Methods = make([]decls.Func, 0, decl.Methods.NumFields())
		for _, fld := range decl.Methods.List {
			if len(fld.Names) == 0 {
				ifaceDecl.Embeds = append(ifaceDecl.Embeds, usages.Type(types.ExprString(fld.Type)))
			} else {
				fieldDecl := parseAstFieldDecl(fld)
				for _, fd := range fieldDecl {
					ifaceDecl.Methods = append(ifaceDecl.Methods, convertFieldToFuncDecl(fd))
				}
			}
		}
	}

	return ifaceDecl
}

func parseAstFieldDecl(fld *ast.Field) []decls.Field {
	tag := ""
	if fld.Tag != nil {
		tag = strings.Trim(fld.Tag.Value, "`")
	}

	typ := usages.Type(types.ExprString(fld.Type))
	doc := strings.Trim(fld.Doc.Text(), "\n")
	comment := strings.Trim(fld.Comment.Text(), "\n")

	if len(fld.Names) == 0 {
		return []decls.Field{{
			Type:     typ,
			Doc:      doc,
			Comment:  comment,
			Tag:      tag,
			Exported: ast.IsExported(typ.String()),
		}}
	}

	result := make([]decls.Field, 0, len(fld.Names))
	for i := range fld.Names {
		result = append(result, decls.Field{
			Name:     fld.Names[i].Name,
			Type:     typ,
			Doc:      doc,
			Comment:  comment,
			Tag:      tag,
			Exported: fld.Names[i].IsExported(),
		})
	}

	return result
}

func parseAstValueSpec(spec *ast.ValueSpec) []decls.Value {
	typ := types.ExprString(spec.Type)
	if typ == "(ast: <nil>)" {
		typ = ""
	}

	result := make([]decls.Value, 0, len(spec.Names))
	for i := range spec.Names {
		result = append(result, decls.Value{
			Name:     spec.Names[i].Name,
			Type:     usages.Type(typ),
			Doc:      strings.Trim(spec.Doc.Text(), "\n"),
			Comment:  strings.Trim(spec.Comment.Text(), "\n"),
			Exported: spec.Names[i].IsExported(),
			Value:    types.ExprString(spec.Values[i]),
		})
	}

	return result
}

func convertFieldToFuncDecl(fieldDecl decls.Field) decls.Func {
	funcDecl := decls.Func{
		Name:     fieldDecl.Name,
		Doc:      strings.Trim(fieldDecl.Doc, "\n"),
		Exported: fieldDecl.Exported,
	}

	expr, err := parser.ParseExpr(fieldDecl.Type.String())
	if err != nil {
		return funcDecl
	}

	funcType, ok := expr.(*ast.FuncType)
	if !ok {
		return funcDecl
	}

	if funcType.TypeParams != nil && funcType.TypeParams.NumFields() > 0 {
		funcDecl.TypeParams = parseAstFieldList(funcType.TypeParams)
	}
	if funcType.Params != nil && funcType.Params.NumFields() > 0 {
		funcDecl.Params = parseAstFieldList(funcType.Params)
	}
	if funcType.Results != nil && funcType.Results.NumFields() > 0 {
		funcDecl.Results = parseAstFieldList(funcType.Results)
	}

	return funcDecl
}

func parseAstFieldList(list *ast.FieldList) []decls.Field {
	result := make([]decls.Field, 0, list.NumFields())
	for _, fld := range list.List {
		fieldDecls := parseAstFieldDecl(fld)
		result = append(result, fieldDecls...)
	}

	return result
}
