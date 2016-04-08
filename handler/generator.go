package handler

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"github.com/facebookgo/stackerr"
)

const (
	optionsTagName = "options"
	omitTag        = "_omit"
)

func GenFromStructType(tspec *ast.TypeSpec) ([]ast.Decl, error) {
	decls := []ast.Decl{}
	if _, ok := tspec.Type.(*ast.StructType); !ok {
		return decls, stackerr.Newf("not a struct type %v", tspec)
	}
	for _, field := range tspec.Type.(*ast.StructType).Fields.List {
		decl := genOptionFromField(tspec.Name.Name, field)
		if decl != nil {
			decls = append(decls, decl)
		}
	}
	return decls, nil
}

func genOptionFromField(structName string, field *ast.Field) *ast.FuncDecl {
	tag := getTag(field)
	if tag == omitTag {
		// skip this field if omitted
		return nil
	}

	svarName := strings.ToLower(structName[0:1])

	if len(field.Names) == 0 {
		// skip anonymous field
		return nil
	}

	fieldName := field.Names[0].Name

	outerType := &ast.FuncType{
		Params: &ast.FieldList{List: []*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent(fieldName)},
				Type:  field.Type,
			},
		}},
		Results: &ast.FieldList{List: []*ast.Field{
			{
				Type: &ast.FuncType{
					Params: &ast.FieldList{List: []*ast.Field{
						{
							Type: &ast.StarExpr{X: ast.NewIdent(structName)},
						}},
					}},
			},
		}},
	}

	innerParams := &ast.FieldList{List: []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(svarName)},
			Type:  &ast.StarExpr{X: ast.NewIdent(structName)},
		}},
	}

	assignStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(svarName),
				Sel: ast.NewIdent(fieldName),
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{ast.NewIdent(fieldName)},
	}

	retFunc := &ast.FuncLit{
		Type: &ast.FuncType{Params: innerParams},
		Body: &ast.BlockStmt{List: []ast.Stmt{assignStmt}},
	}
	outerBody := &ast.BlockStmt{List: []ast.Stmt{
		&ast.ReturnStmt{
			Results: []ast.Expr{retFunc},
		}},
	}

	nameSuffix := ""
	if tag != "" {
		nameSuffix = tag
	} else {
		nameSuffix = strings.ToUpper(fieldName[0:1]) + fieldName[1:len(fieldName)]

	}
	return &ast.FuncDecl{
		Name: ast.NewIdent("Option" + nameSuffix),
		Type: outerType,
		Body: outerBody,
	}
}

func getTag(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	tags := field.Tag.Value
	if tags != "" && tags[0] == '`' && tags[len(tags)-1] == '`' {
		return reflect.StructTag(tags[1:len(tags)]).Get(optionsTagName)
	}
	return ""
}
