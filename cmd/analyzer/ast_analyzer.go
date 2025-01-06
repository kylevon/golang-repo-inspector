package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Analyzer struct {
	fset *token.FileSet
}

func New() *Analyzer {
	return &Analyzer{
		fset: token.NewFileSet(),
	}
}

// AnalyzeFile parses and analyzes a single Go file
func (a *Analyzer) AnalyzeFile(filepath string) (*ast.File, error) {
	node, err := parser.ParseFile(a.fset, filepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// FindFunctionCalls finds all function calls in the AST
func (a *Analyzer) FindFunctionCalls(node ast.Node) []string {
	var calls []string

	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				calls = append(calls, ident.Name)
			}
		}
		return true
	})

	return calls
}

// FindStructs finds all struct declarations
func (a *Analyzer) FindStructs(node ast.Node) []string {
	var structs []string

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				structs = append(structs, typeSpec.Name.Name)
			}
		}
		return true
	})

	return structs
}

// FindImports returns all imports in the file
func (a *Analyzer) FindImports(file *ast.File) []string {
	var imports []string

	for _, imp := range file.Imports {
		if imp.Path != nil {
			// Remove quotes from import path
			importPath := imp.Path.Value[1 : len(imp.Path.Value)-1]
			imports = append(imports, importPath)
		}
	}

	return imports
}
