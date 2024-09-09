package main

import (
	"encoding/json"
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
	jsonData := p.Parse()
	astData := p.GetAST()

	litter.Config.StripPackageNames = true
	fmt.Println("Json: ")
	litter.Dump(jsonData)
	fmt.Println("\nAST: ")
	litter.Dump(astData)

	outputJSON, err := json.Marshal(astData)
	file, err := os.Create("ast.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = file.Write(outputJSON)
	if err != nil {
		fmt.Println(err)
		return
	}

	outputJSON, err = json.Marshal(jsonData)
	file, err = os.Create("data.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = file.Write(outputJSON)
	if err != nil {
		fmt.Println(err)
		return
	}

	// f, _ := os.Open(filePath)
	// var jsonData map[string]interface{}
	// decoder := json.NewDecoder(f)
	// _ = decoder.Decode(&jsonData)
	// litter.Dump(jsonData)
}
