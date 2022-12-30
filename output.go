package aster

import (
	"go/ast"
	"reflect"
)

func resultFromDecl(decl ast.Decl) (Result, error) {
	r := Result{}

	switch reflect.TypeOf(decl).Elem().Name() {
	case "GenDecl":
		genDecl := decl.(*ast.GenDecl)
		r.Comments = genDecl.Doc.Text()
		switch reflect.TypeOf(genDecl.Specs[0]).Elem().Name() {
		case "TypeSpec":
			typeSpec := genDecl.Specs[0].(*ast.TypeSpec)
			r.Name = typeSpec.Name.Name
			switch reflect.TypeOf(typeSpec.Type).Elem().Name() {
			case "StructType":
				structType := typeSpec.Type.(*ast.StructType)
				r.Type = "struct"
				for _, f := range structType.Fields.List {
					fld := Field{
						Name:     f.Names[0].Name,
						Comments: f.Doc.Text(),
						Type:     f.Type.(*ast.Ident).Name,
					}
					r.Fields = append(r.Fields, fld)
				}
			}
		}
	case "FuncDecl":
		funcDecl := decl.(*ast.FuncDecl)
		r.Comments = funcDecl.Doc.Text()
		r.Name = funcDecl.Name.Name
		r.Type = "func"
	}

	return r, nil
}

// Result is an abstracted / simplified view of the matching Go objects,
// their attached comments, and, if applicable, fields. The idea here is to
// hide the complexity of the underlying 'ast' package structures.
type Result struct {
	// The name of the object
	Name string
	// The type of the object, accessible via the 'Name()' method
	Type string
	// Comments associated with the object
	Comments string
	// Fields for a struct or func
	Fields []Field
}

// A Field for a struct or funtion
type Field struct {
	// The name of the field
	Name string
	// The type of the field
	Type string
	// Comments associated with the field
	Comments string
}
