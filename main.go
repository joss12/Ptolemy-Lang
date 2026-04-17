package main

import (
	"fmt"
	"os"

	"github.com/joss12/ptolemy-lang/src"
)

func main() {

	args := os.Args[1:]

	//No argument provided
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "run":
		if len(args) < 2 {
			fmt.Println("Error: missing file path")
			fmt.Println("Usage: ptolemy run <file.po>")
			os.Exit(1)
		}
		runFile(args[1])

	case "check":
		if len(args) < 2 {
			fmt.Println("Error: missing file path")
			fmt.Println("Usage: ptolemy check <file.po>")
			os.Exit(1)
		}
		checkFile(args[1])

	case "version":
		fmt.Println("Ptolemy Lang v1.0.0")

	case "help":
		printUsage()

	default:
		//Allow running directly: ptolemy file.mini
		runFile(args[0])

	}
}

// runFile lexes, parses and evaluates a .mini '.oz' file
func runFile(path string) {
	//read file
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: cannot open file '%s'\n", path)
		os.Exit(1)
	}

	//Lex
	lexer := src.NewLexer(string(source))

	//Parse
	parser := src.NewParser(lexer)
	program := parser.ParseProgram()

	//Check for parse errors
	if len(parser.Errors()) > 0 {
		printParseErrors(parser.Errors())
		os.Exit(1)
	}

	//Evaluate
	env := src.NewEnvironment()
	result := src.Eval(program, env)

	//Print runtime errors
	if result != nil {
		if errObj, ok := result.(*src.ErrorObject); ok {
			fmt.Println(errObj.Inspect())
			os.Exit(1)
		}
	}
}

// check syntax without executing
func checkFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: cannot open file '%s'\n", path)
		os.Exit(1)
	}

	lexer := src.NewLexer(string(source))
	parser := src.NewParser(lexer)
	parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		printParseErrors(parser.Errors())
		os.Exit(1)
	}
	fmt.Printf(" %s looks good\n", path)
}

// printParseErrors prints all parser errors
func printParseErrors(errors []string) {
	fmt.Println("Parse errors:")
	for _, msg := range errors {
		fmt.Printf(" ->%s\n", msg)
	}
}

// printUsage
func printUsage() {
	fmt.Println(`
	Ptolemy Lang - Educational Programming Language

Usage:
  ptolemy <file.po>          Run a program
  ptolemy run <file.po>      Run a program
  ptolemy check <file.po>    Check syntax without running
  ptolemy version              Show version
  ptolemy help                 Show this help message

Examples:
  ptolemy examples/hello.po
  ptolemy run examples/fibonacci.po
  ptolemy check examples/loops.po
	`)
}
