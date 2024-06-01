package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chains-project/capslock-analysis/gosurface/analysis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <module_path>")
		return
	}

	modulePath := os.Args[1]

	asciiArt := `
                                                                                                               
  ,ad8888ba,                 ad88888ba                               ad88                                      
 d8"'    '"8b               d8"     "8b                             d8"                                        
d8'                         Y8,                                     88                                         
88              ,adPPYba,   'Y8aaaaa,    88       88  8b,dPPYba,  MM88MMM  ,adPPYYba,   ,adPPYba,   ,adPPYba,  
88      88888  a8"     "8a    '"""""8b,  88       88  88P'   "Y8    88     ""     '8  a8"     ""  a8P_____88  
Y8,        88  8b       d8          '8b  88       88  88            88     ,adPPPPP88  8b          8PP"""""""  
 Y8a.    .a88  "8a,   ,a8"  Y8a     a8P  "8a,   ,a88  88            88     88,    ,88  "8a,   ,aa  "8b,   ,aa  
  '"Y88888P"    '"YbbdP"'    "Y88888P"    '"YbbdP'Y8  88            88     '"8bbdP"Y8   '"Ybbd8"'   '"Ybbd8"'  
                                                                                                               
                                                                                                          "
`
	fmt.Println(asciiArt)

	fmt.Println("GoSurface is a tool that aims to analyze the potential attack surface of open-source Go packages and modules.")
	fmt.Println("It looks for occurrences of various features and constructs that could potentially introduce security risks.")
	fmt.Println()

	// Get paths of packages imported by module (it includes the main package)
	// TODO: currently only fetches subdirectories in module, not external dependencies
	dependencies, err := analysis.GetDependencies(modulePath) // TODO rename get module/packages
	if err != nil {
		fmt.Printf("Error getting files in module: %v\n", err)
		return
	}

	// Print the dependencies
	//analysis.PrintDependencies(dependencies)

	// Analyze module and its dependencies
	for _, dep := range dependencies {
		analysis.AnalyzeModule(dep.Path, &analysis.InitOccurrences, analysis.InitFuncParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.AnonymOccurrences, analysis.AnonymFuncParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.ExecOccurrences, analysis.ExecParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.PluginOccurrences, analysis.PluginParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.GoGenerateOccurrences, analysis.GoGenerateParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.UnsafeOccurrences, analysis.UnsafeParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.CgoOccurrences, analysis.CgoParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.IndirectOccurrences, analysis.IndirectParser{})
		analysis.AnalyzeModule(dep.Path, &analysis.ReflectOccurrences, analysis.ReflectParser{})
	}

	// Convert occurrences to JSON
	occurrences := append(append(append(append(append(append(append(append(
		analysis.InitOccurrences,
		analysis.AnonymOccurrences...),
		analysis.ExecOccurrences...),
		analysis.PluginOccurrences...),
		analysis.GoGenerateOccurrences...),
		analysis.UnsafeOccurrences...),
		analysis.CgoOccurrences...),
		analysis.IndirectOccurrences...),
		analysis.ReflectOccurrences...)
	// Print occurrences
	analysis.PrintOccurrences(analysis.IndirectOccurrences)
	// analysis.PrintOccurrences(occurrences)

	// Count unique occurrences
	initCount, anonymCount, osExecCount, pluginCount, goGenerateCount, unsafeCount, cgoCount, indirectCount, reflectCount := analysis.CountUniqueOccurrences(occurrences)
	fmt.Println()
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ Attack Surface Analysis: %s			       ║\n", filepath.Base(strings.TrimSuffix(modulePath, "/"+filepath.Base(modulePath))))
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Unique occurrences of init() function:            %10d ║\n", initCount)
	fmt.Printf("║ Initialization with anonymous function:           %10d ║\n", anonymCount)
	fmt.Printf("║ Invocation from the os/exec package:              %10d ║\n", osExecCount)
	fmt.Printf("║ Plugin dynamically loaded:                        %10d ║\n", pluginCount)
	fmt.Printf("║ go:generate directive:                            %10d ║\n", goGenerateCount)
	fmt.Printf("║ Unsafe pointers:                                  %10d ║\n", unsafeCount)
	fmt.Printf("║ CGO pointers:                                     %10d ║\n", cgoCount)
	fmt.Printf("║ Indirect method calls via interfaces:             %10d ║\n", indirectCount)
	fmt.Printf("║ Imports of reflection:			    %10d ║\n", reflectCount)
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}
