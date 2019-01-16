package main

import (
	"fmt"
)

var input string

func init() {
	fmt.Scanf("%s\n", &input)
}

func main() {
	answer := camelcase(input)
	fmt.Println(answer)
}

func camelcase(s string) int32 {
	if s == "" {
		return 0
	}

	var count int32 = 1

	for _, v := range s {
		if isUpper(v) {
			count++
		}
	}

	return count
}

func isUpper(v rune) bool {
	return v >= 'A' && v <= 'Z'
}
