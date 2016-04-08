package handler

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/facebookgo/stackerr"
	"github.com/facebookgo/structtag"
)

const (
	optionsTagName = "options"
)

type GenResult struct {
	Comment *ast.CommentGroup
	Decl    ast.Decl
	Name    string
}

func GenFromStructType(prefix string, tspec *ast.TypeSpec) ([]*GenResult, error) {
	gened := []*GenResult{}
	if _, ok := tspec.Type.(*ast.StructType); !ok {
		return gened, stackerr.Newf("not a struct type %v", tspec)
	}
	for _, field := range tspec.Type.(*ast.StructType).Fields.List {
		res := genOptionFromField(tspec.Name.Name, field, prefix)
		if res != nil {
			gened = append(gened, res)
		}
	}
	return gened, nil
}

func genOptionFromField(structName string, field *ast.Field, prefix string) *GenResult {
	found, tag, err := getTag(field)
	if err != nil || found == false {
		// skip a field by default if no field tag
		return nil
	}
	if len(field.Names) == 0 {
		// skip anonymous field
		return nil
	}

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

	nameSuffix := ""
	if tag != "" {
		nameSuffix = tag
	} else {
		nameSuffix = strings.ToUpper(fieldName[0:1]) + fieldName[1:len(fieldName)]
	}
	name := prefix + nameSuffix
	return &GenResult{
		Decl: &ast.FuncDecl{
			Name: ast.NewIdent(name),
			Type: outerType,
			Body: outerBody,
		},
		Comment: getDoc(field, name),
		Name:    name,
	}
}

func getTag(field *ast.Field) (bool, string, error) {
	if field.Tag == nil {
		return false, "", nil
	}
	tags := field.Tag.Value
	if tags != "" && tags[0] == '`' && tags[len(tags)-1] == '`' {
		return structtag.Extract(optionsTagName, tags[1:len(tags)])
	}
	return false, "", nil
}

func getDoc(field *ast.Field, funcName string) *ast.CommentGroup {
	if field.Doc == nil {
		return nil
	}
	c := field.Doc.List[0]
	if strings.HasPrefix(c.Text, "//") {
		c.Text = "// " + funcName + c.Text[2:len(c.Text)]
	}

	return field.Doc
}
