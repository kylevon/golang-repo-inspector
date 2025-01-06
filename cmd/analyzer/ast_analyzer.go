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

// FindFunctionCalls finds all function calls in the AST with their arguments
func (a *Analyzer) FindFunctionCalls(node ast.Node) map[string][]string {
	calls := make(map[string][]string)

	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			var funcName string
			switch fun := call.Fun.(type) {
			case *ast.Ident:
				funcName = fun.Name
			case *ast.SelectorExpr:
				if pkg, ok := fun.X.(*ast.Ident); ok {
					funcName = pkg.Name + "." + fun.Sel.Name
				}
			}

			// Extract arguments
			var args []string
			for _, arg := range call.Args {
				switch v := arg.(type) {
				case *ast.BasicLit:
					// For literal values (strings, numbers, etc.)
					args = append(args, v.Value)
				case *ast.Ident:
					// For variable names
					args = append(args, "$"+v.Name)
				case *ast.SelectorExpr:
					// For package-qualified values (e.g., os.Stdout)
					if pkg, ok := v.X.(*ast.Ident); ok {
						args = append(args, pkg.Name+"."+v.Sel.Name)
					}
				case *ast.CallExpr:
					// For nested function calls
					args = append(args, "nested_call")
				default:
					args = append(args, "complex_expr")
				}
			}

			calls[funcName] = args
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
