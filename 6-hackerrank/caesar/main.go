package main

import (
	"fmt"
)

var length, delta int
var input string

func init() {
	fmt.Scanf("%d\n", &length)
	fmt.Scanf("%s\n", &input)
	fmt.Scanf("%d\n", &delta)
}

func main() {
	var ret []rune
	for _, ch := range input[:length] {
		ret = append(ret, cipher(ch, delta))
	}
	for _, ch := range input[length:] {
		ret = append(ret, ch)
	}
	fmt.Println(string(ret))
}

func cipher(r rune, delta int) rune {
	if r >= 'A' && r <= 'Z' {
		return rotate(r, 'A', delta)
	}
	if r >= 'a' && r <= 'z' {
		return rotate(r, 'a', delta)
	}
	return r
}

func rotate(r rune, base, delta int) rune {
	tmp := int(r) - base
	tmp = (tmp + delta) % 26
	return rune(tmp + base)
}
