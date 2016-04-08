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

func GenFromStructType(prefix string, tspec *ast.TypeSpec) ([]ast.Decl, error) {
	decls := []ast.Decl{}
	if _, ok := tspec.Type.(*ast.StructType); !ok {
		return decls, stackerr.Newf("not a struct type %v", tspec)
	}
	for _, field := range tspec.Type.(*ast.StructType).Fields.List {
		decl := genOptionFromField(tspec.Name.Name, field, prefix)
		if decl != nil {
			decls = append(decls, decl)
		}
	}
	return decls, nil
}

func genOptionFromField(structName string, field *ast.Field, prefix string) *ast.FuncDecl {
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
	outerName := prefix + nameSuffix
	return &ast.FuncDecl{
		Name: ast.NewIdent(outerName),
		Type: outerType,
		Body: outerBody,
		Doc:  getDoc(field, outerName),
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

func getDoc(field *ast.Field, funcName string) *ast.CommentGroup {
	if field.Doc == nil {
		return nil
	}

	// to make go lint happy, only works with "//" style for now
	fstComment := field.Doc.List[0].Text
	c := &ast.Comment{}
	if strings.HasPrefix(fstComment, "//") {
		c.Text = "// " + funcName + fstComment[2:len(fstComment)]
	} else {
		c.Text = fstComment
	}
	doc := &ast.CommentGroup{}
	doc.List = make([]*ast.Comment, len(field.Doc.List), len(field.Doc.List))
	copy(doc.List, field.Doc.List)
	doc.List[0] = c
	return doc
}
