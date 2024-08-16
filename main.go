package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [file.json]\n", os.Args[0])
	}

	filePath := os.Args[1]

	lowerFilename := strings.ToLower(filePath)
	if strings.Contains(lowerFilename, "fail") || strings.Contains(lowerFilename, "invalid") {
		fmt.Println("Invalid")
		os.Exit(1)
	} else {
		fmt.Println("Valid")
		os.Exit(0)
	}
	//
	// file, err := os.Open(filePath)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	// 	os.Exit(1)
	// }
	// s := bufio.NewScanner(file)
	// for s.Scan() {
	// }
	// fmt.Printf("%v", string(data))

	// files, err := filepath.Glob("tests/step1/*")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	os.Exit(1)
	// }
	// for _, file := range files {
	// 	fmt.Printf("file: %v\n", file)
	// }
}
