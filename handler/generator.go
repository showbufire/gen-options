package handler

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/facebookgo/stackerr"
)

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

	optionType := &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{Type: &ast.StarExpr{X: ast.NewIdent(structName)}},
			},
		},
	}
	outerType := &ast.FuncType{
		Params: &ast.FieldList{List: []*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent(fieldName)},
				Type:  field.Type,
			},
		}},
		Results: &ast.FieldList{List: []*ast.Field{
			{
				Type: optionType,
			},
		}},
	}

	innerParams := &ast.FieldList{List: []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(svarName)},
			Type:  optionType,
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

	return &ast.FuncDecl{
		Name: ast.NewIdent("Option" + fieldName), // TODO: load field tag
		Type: outerType,
		Body: outerBody,
	}
}
