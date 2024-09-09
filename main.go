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
	// filePath := "./testdata/step2/valid.json" // FIXME: for debugging only

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	lex := lexer.NewLexer(string(data))
	p := parser.NewParser(lex)
	stmt := p.Parse().Body
	litter.Config.StripPackageNames = true
	litter.Dump(stmt)
}
