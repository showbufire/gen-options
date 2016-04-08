package handler

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"github.com/facebookgo/stackerr"
)

const optionsTagName = "options"

func GenFromStructType(tspec *ast.TypeSpec) ([]ast.Decl, error) {
	decls := []ast.Decl{}
	if _, ok := tspec.Type.(*ast.StructType); !ok {
		return decls, stackerr.Newf("not a struct type %v", tspec)
	}
	for _, field := range tspec.Type.(*ast.StructType).Fields.List {
		decl := genOptionFromField(tspec.Name.Name, field)
		decls = append(decls, decl)
	}
	return decls, nil
}

func genOptionFromField(structName string, field *ast.Field) *ast.FuncDecl {
	svarName := strings.ToLower(structName[0:1])
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

	tags := ""
	if field.Tag != nil {
		tags = field.Tag.Value
	}
	outerName := getOptionFuncName(tags, fieldName)

	return &ast.FuncDecl{
		Name: ast.NewIdent(outerName), // TODO: load field tag
		Type: outerType,
		Body: outerBody,
	}
}

func getOptionFuncName(tags string, fieldName string) string {
	if tags != "" && tags[0] == '`' && tags[len(tags)-1] == '`' {
		tag := reflect.StructTag(tags[1:len(tags)]).Get(optionsTagName)
		if tag != "" {
			return "Option" + tag
		}
	}
	return "Option" + strings.ToUpper(fieldName[0:1]) + fieldName[1:len(fieldName)]
}
