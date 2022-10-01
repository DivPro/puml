package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type PField struct {
	Name    string
	Type    string
	Tag     string
	IsPtr   bool
	IsSlice bool
	IsEmbed bool
}

type PStruct struct {
	File    string
	Package string
	Name    string
	Fields  []PField
}

func ParseFile(path string) ([]PStruct, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	var (
		res       []PStruct
		curStruct *PStruct
	)
	ast.Inspect(node, func(n ast.Node) bool {
		switch curNode := n.(type) {
		case *ast.TypeSpec:
			_, ok := curNode.Type.(*ast.StructType)
			if !ok {
				break
			}
			curStruct = &PStruct{
				File:    path,
				Package: node.Name.String(),
				Name:    curNode.Name.String(),
			}
			return true
		case *ast.FieldList:
			if curStruct == nil {
				break
			}
			for _, f := range curNode.List {
				curStruct.Fields = append(curStruct.Fields, visitField(f))
			}
			res = append(res, *curStruct)
			curStruct = nil
			return false
		}

		return true
	})

	return res, nil
}

func visitField(f *ast.Field) PField {
	var res PField
	if f.Tag != nil {
		res.Tag = f.Tag.Value
	}
	el := &res.Name
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch curNode := n.(type) {
		case *ast.ArrayType:
			res.IsSlice = true
		case *ast.SelectorExpr:
			switch sub := curNode.X.(type) {
			case *ast.Ident:
				*el = sub.String() + "."
			}
			*el = *el + curNode.Sel.String()
			el = &res.Type
			return false
		case *ast.Ident:
			*el = curNode.String()
			el = &res.Type
			return false
		case *ast.StarExpr:
			res.IsPtr = true
		}
		return true
	})
	if res.Type == "" {
		res.IsEmbed = true
	}

	return res
}
