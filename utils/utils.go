package utils

import "fmt"

func PosToLineColumn(input string) (int, int) {
	// Converts pos in string to line and column
	fmt.Printf("input: %q\n", input)
	pos := 14
	line, column := 1, 1
	for i := 0; i < pos; i++ {
		if pos >= len(input) {
			fmt.Println("Error: Out of bound")
			break
		}
		char := input[i]
		if char == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}
	return line, column
}
