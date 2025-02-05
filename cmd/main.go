package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang-repo-inspector/cmd/analyzer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a directory path to analyze")
		os.Exit(1)
	}

	analyzer := analyzer.New()
	dirPath := os.Args[1]

	// Verify it's a directory
	info, err := os.Stat(dirPath)
	if err != nil {
		log.Fatalf("Error accessing directory: %v", err)
	}
	if !info.IsDir() {
		log.Fatalf("Path %s is not a directory", dirPath)
	}

	// Walk through all .go files in the directory
	var allImports, allStructs []string
	allFunctionCalls := make(map[string][]string)
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			file, err := analyzer.AnalyzeFile(path)
			if err != nil {
				return fmt.Errorf("error analyzing %s: %v", path, err)
			}

			imports := analyzer.FindImports(file)
			structs := analyzer.FindStructs(file)
			functionCalls := analyzer.FindFunctionCalls(file)

			allImports = append(allImports, imports...)
			allStructs = append(allStructs, structs...)
			// Merge function calls
			for name, args := range functionCalls {
				allFunctionCalls[name] = args
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}

	fmt.Println("Imports found:")
	for _, imp := range allImports {
		fmt.Printf("  - %s\n", imp)
	}

	fmt.Println("\nStructs found:")
	for _, str := range allStructs {
		fmt.Printf("  - %s\n", str)
	}

	fmt.Println("\nFunction calls found:")
	for name, args := range allFunctionCalls {
		if len(args) > 0 {
			fmt.Printf("  - %s(%s)\n", name, strings.Join(args, ", "))
		} else {
			fmt.Printf("  - %s()\n", name)
		}
	}
}
