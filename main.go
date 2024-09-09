package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hussein-mourad/go-json-parser/lexer"
	"github.com/hussein-mourad/go-json-parser/parser"
	"github.com/sanity-io/litter"
)

var logger = log.New(os.Stdout, "", 0)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [file.json]\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]
	// filePath := "./testdata/final/fail22.json" // FIXME: for debugging only

	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Fatalln(err)
	}
	fmt.Println("Input Data: ")
	fmt.Println(string(data))
	fmt.Println()

	lex := lexer.NewLexer(string(data))
	p := parser.NewParser(lex)
	json := p.Parse()
	ast := p.GetAST()

	litter.Config.StripPackageNames = true
	fmt.Println("Json: ")
	litter.Dump(json)
	fmt.Println("\nAST: ")
	litter.Dump(ast)

	// f, _ := os.Open(filePath)
	// var jsonData map[string]interface{}
	// decoder := json.NewDecoder(f)
	// _ = decoder.Decode(&jsonData)
	// litter.Dump(jsonData)
}
