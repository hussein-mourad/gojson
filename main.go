package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [file.json]\n", os.Args[0])
	}

	filePath := os.Args[1]

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	lexer := NewLexer(string(data))
	for {
		t := lexer.NextToken()
		fmt.Printf("%v\n", t)
		if t.Type == EOF {
			break
		}
	}
}
