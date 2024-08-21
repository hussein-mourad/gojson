package main

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", 0)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [file.json]\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]
	// filePath := "test.json" // for debugger

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	lexer := NewLexer(string(data))
	parser := NewParser(lexer)
	json := parser.parse()

	// Assert the type to map[string]int
	m, ok := json.(map[string]interface{})
	if !ok {
		fmt.Println("Value is not of type map[string]interface{}")
		return
	}
	for k, v := range m {
		fmt.Printf("k: %v, v: %v\n", k, v)
	}
}
