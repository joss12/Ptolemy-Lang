# Ptolemy Lang

A statically-scoped, dynamically-typed programming language built from scratch in Go. Features lexical closures, first-class functions, recursion, and a tree-walking interpreter.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

## Overview

Ptolemy Lang is a personal passion project born from deep diving into "Crafting Interpreters" by Robert Nystrom and the CodeCrafters.io interpreter track. What started as curiosity about how programming languages actually work turned into weeks of intense learning—experiencing the same joy a child feels discovering ice cream for the first time.

Building an interpreter from scratch revealed the beautiful complexity hidden beneath every line of code we write. Every challenge—from implementing Prat parsing to making closures capture their environments correctly—was a puzzle that demanded understanding, not just implementation.

This project represents the transition from using languages to understanding them at a fundamental level. It's a testament to the idea that the best way to truly learn something is to build it yourself, one token at a time.

**Key Features:**
- Lexical closures with captured environments
- First-class functions and higher-order programming
- Tail-recursive optimization-friendly design
- Pratt parser for correct operator precedence
- Dynamic typing with runtime type checking
- Garbage collection via Go's memory management

## Quick Start

### Prerequisites
- Go 1.21 or higher

### Installation

**Option 1: Build from source**
```bash
git clone https://github.com/joss12/ptolemy-lang.git
cd ptolemy-lang
go build -o ptolemy
```

### Running Programs

```bash
# Execute a program
./ptolemy examples/fibonacci.mini

# Check syntax without execution
./ptolemy check examples/closures.mini

# Show version
./ptolemy version
```

## Language Syntax

### Variables and Constants

```javascript
let mutableVar = 42
cst immutableConst = 3.14159

mutableVar = 100  // OK
// immutableConst = 3.14  // Error: cannot reassign constant
```

### Data Types

```javascript
// Numbers (float64 internally)
let integer = 42
let float = 3.14159
let negative = -17

// Strings (double-quotes only)
let text = "Hello, World"
let concatenated = "Hello" + " " + "World"

// Booleans
let isTrue = true
let isFalse = false

// Null
let empty = null

// Arrays (heterogeneous, dynamically-sized)
let numbers = [1, 2, 3, 4, 5]
let mixed = [42, "text", true, null]
```

### Functions

```javascript
// Function declaration
fn add(a, b) {
    return a + b
}

// Functions are first-class values
let operation = add
let result = operation(5, 3)  // 8

// Higher-order functions
fn applyTwice(func, value) {
    return func(func(value))
}

fn double(x) { return x * 2 }
print(applyTwice(double, 5))  // 20
```

### Closures

```javascript
// Functions capture their lexical environment
fn makeCounter() {
    let count = 0
    
    fn increment() {
        count = count + 1
        return count
    }
    
    return increment
}

let counter = makeCounter()
print(counter())  // 1
print(counter())  // 2
print(counter())  // 3
```

### Control Flow

```javascript
// Conditionals
if (x > 10) {
    print("Large")
} else if (x > 5) {
    print("Medium")
} else {
    print("Small")
}

// While loops
let i = 0
while (i < 5) {
    print(i)
    i = i + 1
}

// For loops
for (let j = 0; j < 10; j = j + 1) {
    print(j)
}
```

### Recursion

```javascript
fn fibonacci(n) {
    if (n <= 1) {
        return n
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

print(fibonacci(10))  // 55
```

### Built-in Functions

```javascript
print("Hello", "World")    // Output: Hello World
len([1, 2, 3])             // Returns: 3
type(42)                   // Returns: "NUMBER"
push(arr, value)           // Append to array
pop(arr)                   // Remove and return last element
```

## Architecture

### Execution Pipeline

```
Source Code (.po)
        ↓
    [Lexer]          → Tokenization
        ↓
    [Parser]         → AST Construction (Pratt Parsing)
        ↓
    [Evaluator]      → Tree-Walking Interpretation
        ↓
    Output
```

### Components

**Lexer (`src/lexer.go`)**
- Converts source text into token stream
- Handles single-line comments (`//`)
- Recognizes keywords, operators, and literals
- Tracks line/column for error reporting

**Parser (`src/parser.go`)**
- Implements Pratt parser (Top-Down Operator Precedence)
- Builds Abstract Syntax Tree (AST)
- Handles operator precedence correctly
- Generates detailed parse errors

**AST (`src/ast.go`)**
- Defines node types for all language constructs
- Separates statements (side effects) from expressions (values)
- Implements visitor-friendly interface

**Evaluator (`src/evaluator.go`)**
- Tree-walking interpreter
- Environment-based variable binding
- Lexical scoping with closure support
- Runtime type checking and error handling

**Object System (`src/object.go`)**
- Runtime value representation
- Type system implementation
- Built-in function interface

**Environment (`src/environment.go`)**
- Variable storage with scope chaining
- Closure environment capture
- Scope management for blocks and functions

## Examples

See the `examples/` directory for comprehensive demonstrations:

| File | Description | Concepts |
|------|-------------|----------|
| `01_variables.mini` | Variable scoping | let, cst, scope |
| `10_closures.mini` | Closure mechanics | Environment capture |
| `11_recursion.mini` | Recursive algorithms | Base cases, tail calls |
| `15_sorting.mini` | Sorting algorithms | Bubble, selection, insertion |
| `20_memoization.mini` | Dynamic programming | Optimization via caching |

Run all examples:
```bash
for file in examples/*.mini; do
    echo "Running $file..."
    ./ptolemy "$file"
    echo "---"
done
```

## Implementation Notes

### Design Decisions

**Why Go?**
- Manual memory management abstracted by GC
- Strong typing catches implementation errors
- Fast compilation and execution
- Excellent standard library for file I/O and string manipulation

**Why Tree-Walking Interpreter?**
- Simplicity: Direct AST execution without bytecode compilation
- Clarity: Easy to understand and modify
- Educational: Clear mapping between source and execution
- Trade-off: Performance sacrificed for implementation simplicity

**Why No Objects/Hash Maps?**
- Focus: Demonstrates core interpreter concepts without overwhelming complexity
- Scope: Arrays and closures provide sufficient data structures for educational purposes
- Future: Can be added as enhancement without architectural changes

### Limitations

- No module system or imports
- No string interpolation (`` `template ${var}` ``)
- No break/continue statements in loops
- No exception handling (errors stop execution)
- No standard library beyond built-ins
- Single-quote strings not supported (use double-quotes)

### Performance Characteristics

- **Lexing**: O(n) in source length
- **Parsing**: O(n) in token count
- **Evaluation**: Dependent on program (tree-walking overhead)
- **Recursion**: Limited by Go's call stack (no TCO)
- **Closures**: O(1) environment lookup via hash map

## Testing

Run the test suite:
```bash
go test ./tests -v
```

Test individual components:
```bash
go test ./tests -run TestLexer
go test ./tests -run TestParser
go test ./tests -run TestEvaluator
```

## Contributing

This is an a personal project built to understand how compiler are designed from scratch. Contributions that enhance clarity or add well-documented features are welcome.

**Areas for Enhancement:**
- Additional built-in functions
- Performance optimizations
- Better error messages with suggestions
- REPL (Read-Eval-Print-Loop)
- Syntax highlighting support
- Debugger interface

## Technical Details

### Operator Precedence (Lowest to Highest)

```
1. Assignment        =
2. Logical OR        ||
3. Logical AND       &&
4. Equality          == !=
5. Comparison        < > <= >=
6. Addition          + -
7. Multiplication    * / %
8. Unary             - !
9. Call              func()
10. Index            arr[i]
```

### Scope Rules

- **Lexical scoping**: Variables resolved at definition time, not call time
- **Closure capture**: Inner functions capture outer environment by reference
- **Block scope**: Variables in `{ }` blocks are locally scoped
- **Shadow prevention**: Inner scopes can shadow outer variables

### Memory Model

- **Garbage Collection**: Handled by Go runtime
- **Reference Semantics**: Arrays and functions passed by reference
- **Value Semantics**: Numbers, strings, booleans copied on assignment
- **Environment Chaining**: Parent scopes kept alive by closure references

## Project Structure

```
ptolemy-lang/
├── src/
│   ├── token.go         # Token definitions
│   ├── lexer.go         # Lexical analysis
│   ├── ast.go           # AST node types
│   ├── parser.go        # Syntax analysis
│   ├── object.go        # Runtime value system
│   ├── environment.go   # Variable storage
│   └── evaluator.go     # Execution engine
├── examples/            # Language examples
├── tests/               # Unit tests
├── main.go              # CLI entry point
├── go.mod               # Go module file
└── README.md            # This file
```

## Author

Built by [Eddy Mouity](https://github.com/joss12) as a demonstration of interpreter implementation fundamentals.

## Acknowledgments

Inspired by:
- "Writing An Interpreter In Go" by Thorsten Ball
- "Crafting Interpreters" by Robert Nystrom
- The simplicity of Lua's design philosophy
- Classic LISP interpreter implementations

## Resources

**Learn More About Interpreters:**
- [Pratt Parsing Explained](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/)
- [Tree-Walking Interpreters](https://craftinginterpreters.com/a-tree-walk-interpreter.html)
- [Closure Implementation](https://en.wikipedia.org/wiki/Closure_(computer_programming))

---

**Note:** This is a personal project. For curiosity use cases, consider mature languages like Python, JavaScript, or Lua.

