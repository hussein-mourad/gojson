package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hussein-mourad/go-json-parser/lexer"
)

var logger = log.New(os.Stdout, "", 0)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [file.json]\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]
	// filePath := "./testdata/final/pass1.json" // FIXME: for debugging only

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	lex := lexer.NewLexer(string(data))
	for {
		token := lex.NextToken()
		if token == nil {
			break
		}
		fmt.Printf("Type: %v\tValue: %q\n", token.TypeString(), token.Value)
		if token.Type == lexer.EOF {
			break
		}
	}
	// parser := NewParser(lexer)
	// json := parser.parse()
	//
	// fmt.Printf("reflect.TypeOf(json): %v\n", reflect.TypeOf(json))
	// Assert the type to map[string]int
	// m, ok := json.(map[string]interface{})
	// if !ok {
	// 	fmt.Println("Value is not of type map[string]interface{}")
	// 	return
	// }
	// for k, v := range m {
	// 	fmt.Printf("k: %v, v: %v\n", k, v)
	// }
}
